package parser

import (
	"errors"
	"flag"
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"go/types"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
	"strings"
)

type AstParser struct {
	typeName    string
	packageName string
}

func NewAstParser(typeName, packageName string) *AstParser {
	return &AstParser{
		typeName:    typeName,
		packageName: packageName,
	}
}

func (p *AstParser) Parse() (*entity.JsonSchemaMetadata, error) {
	if p.packageName == "." {
		dir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("get current directory failed: %v", err)
		}
		p.packageName, err = packageNameOfDir(dir)
		if err != nil {
			return nil, fmt.Errorf("parse package name failed: %v", err)
		}
	}

	pkg, err := parsePackage(p.packageName, p.typeName)
	return pkg, err
}

var errOutsideGoPath = errors.New("source directory is outside GOPATH")

// packageNameOfDir get package import path via dir
func packageNameOfDir(srcDir string) (string, error) {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return "", err
	}

	var goFilePath string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			goFilePath = file.Name()
			break
		}
	}
	if goFilePath == "" {
		return "", fmt.Errorf("go source file not found %s", srcDir)
	}

	packageImport, err := parsePackageImport(srcDir)
	if err != nil {
		return "", err
	}
	return packageImport, nil
}

// parseImportPackage get package import path via source file
// an alternative implementation is to use:
// cfg := &packages.Config{Mode: packages.NeedName, Tests: true, Dir: srcDir}
// pkgs, err := packages.Load(cfg, "file="+source)
// However, it will call "go list" and slow down the performance
func parsePackageImport(srcDir string) (string, error) {
	moduleMode := os.Getenv("GO111MODULE")
	// trying to find the module
	if moduleMode != "off" {
		currentDir := srcDir
		for {
			dat, err := os.ReadFile(filepath.Join(currentDir, "go.mod"))
			if os.IsNotExist(err) {
				if currentDir == filepath.Dir(currentDir) {
					// at the root
					break
				}
				currentDir = filepath.Dir(currentDir)
				continue
			} else if err != nil {
				return "", err
			}
			modulePath := modfile.ModulePath(dat)
			return filepath.ToSlash(filepath.Join(modulePath, strings.TrimPrefix(srcDir, currentDir))), nil
		}
	}
	// fall back to GOPATH mode
	goPaths := os.Getenv("GOPATH")
	if goPaths == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}
	goPathList := strings.Split(goPaths, string(os.PathListSeparator))
	for _, goPath := range goPathList {
		sourceRoot := filepath.Join(goPath, "src") + string(os.PathSeparator)
		if strings.HasPrefix(srcDir, sourceRoot) {
			return filepath.ToSlash(strings.TrimPrefix(srcDir, sourceRoot)), nil
		}
	}
	return "", errOutsideGoPath
}

var (
	buildFlags = flag.String("build_flags", "", "(package mode) Additional flags for go build.")
)

func parsePackage(packageName string, structName string) (*entity.JsonSchemaMetadata, error) {
	pkg, err := loadPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("load package: %w", err)
	}

	typeMetadata, err := extractMetadataFromPackage(pkg, structName)
	if err != nil {
		return nil, fmt.Errorf("extract typeMetadata from package: %w", err)
	}

	return typeMetadata, nil
}

func loadPackage(packageName string) (*packages.Package, error) {
	var buildFlagsSet []string
	if *buildFlags != "" {
		buildFlagsSet = strings.Split(*buildFlags, " ")
	}

	cfg := &packages.Config{
		Mode:       packages.NeedDeps | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedEmbedFiles | packages.NeedSyntax,
		BuildFlags: buildFlagsSet,
	}
	pkgs, err := packages.Load(cfg, packageName)
	if err != nil {
		return nil, fmt.Errorf("load packages: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("packages length must be 1: %d", len(pkgs))
	}

	if len(pkgs[0].Errors) > 0 {
		errs := make([]error, len(pkgs[0].Errors))
		for i, err := range pkgs[0].Errors {
			errs[i] = err
		}

		return nil, errors.Join(errs...)
	}

	return pkgs[0], nil
}

func extractMetadataFromPackage(pkg *packages.Package, structName string) (*entity.JsonSchemaMetadata, error) {
	obj := pkg.Types.Scope().Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("struct %s does not exist", structName)
	}

	modelIface, err := parseStruct(obj)
	if err != nil {
		return nil, err
	}

	return modelIface, nil
}

func parseStruct(obj types.Object) (*entity.JsonSchemaMetadata, error) {
	named, ok := types.Unalias(obj.Type()).(*types.Named)
	if !ok {
		return nil, fmt.Errorf("%s is not an struct. it is a %T", obj.Name(), obj.Type().Underlying())
	}

	strct, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("%s is not an struct. it is a %T", obj.Name(), obj.Type().Underlying())
	}

	rootMetadata := entity.NewDataTypeMetadata(obj.Pkg().Path(), obj.Name(), "struct", false)

	mainMetadata := entity.NewJsonSchemaMetadata()

	root, _, err := parseStructInRecursion(mainMetadata, strct, rootMetadata)
	if err != nil {
		return nil, err
	}
	mainMetadata.Root = root
	return mainMetadata, nil
}

func parseStructInRecursion(schemaMetadata *entity.JsonSchemaMetadata, typ types.Type, currentMetadata *entity.DataTypeMetadata) (metadata *entity.DataTypeMetadata, isStruct bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unknown error in parseStructInRecursion: %v", r)
		}
	}()
	metadata = currentMetadata

	pointer, isPointer := typ.(*types.Pointer)
	if isPointer {
		typ = pointer.Elem()
	}
	named, isNamed := typ.(*types.Named)
	if isNamed {
		obj := named.Obj()
		//todo скорее всего неправильно выставлять "struct", могут быть другие типы для Named
		metadata = entity.NewDataTypeMetadataWithBaseMetadata(currentMetadata, obj.Pkg().Path(), obj.Name(), "struct", isPointer)
		typ = named.Underlying()
	}

	switch specificType := typ.(type) {
	case *types.Basic:
		metadata = entity.NewDataTypeMetadataWithBaseMetadata(currentMetadata, "", specificType.String(), specificType.String(), false)
		return metadata, false, nil
	case *types.Struct:
		if dataTypeMetadata, ok := schemaMetadata.Types[metadata.ID()]; ok {
			return dataTypeMetadata, true, nil
		}
		schemaMetadata.Types[metadata.ID()] = metadata

		for i := 0; i < specificType.NumFields(); i++ {
			field := specificType.Field(i)
			fieldMetadata := &entity.DataTypeMetadata{
				Package: field.Pkg().Path(),
			}
			tags := parseTags(specificType.Tag(i))
			nodeTypeMetadata, isStruct, err := parseStructInRecursion(schemaMetadata, field.Type(), fieldMetadata)
			if err != nil {
				return nil, false, err
			}

			var nodeMetadata *entity.DataTypeMetadata
			if isStruct {
				nodeMetadata = entity.NewDataTypeRefMetadata(nodeTypeMetadata)
			} else {
				nodeMetadata = entity.NewDataTypeMetadata(nodeTypeMetadata.Package, nodeTypeMetadata.TypeName, nodeTypeMetadata.TypeKind, nodeTypeMetadata.IsPointer)
				nodeMetadata.Nodes = nodeTypeMetadata.Nodes
			}
			nodeMetadata.Tags = tags
			nodeMetadata.FieldName = field.Name()
			_, isPointer = field.Type().(*types.Pointer)
			nodeMetadata.IsPointer = isPointer

			metadata.Nodes = append(metadata.Nodes, nodeMetadata)
		}
		return metadata, true, nil
	default:
		return nil, false, errors.New("incorrect and unexpected ype")
	}
}

func parseTags(tg string) map[string][]string {
	tagsResult := make(map[string][]string)
	tags := strings.Split(tg, " ")
	for _, tag := range tags {
		tmp := strings.Split(tag, ":")
		splited := strings.Split(tmp[1], ",")
		length := len(splited)
		if length > 1 {
			splited[0] = splited[0][1:]
			splited[length-1] = splited[length-1][:len(splited[length-1])-1]
		} else if len(splited) == 1 {
			splited[0] = splited[0][1 : len(splited[0])-1]
		}
		tagsResult[tmp[0]] = splited
	}
	return tagsResult
}
