package ast

import (
	"errors"
	"flag"
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"go/types"
	"golang.org/x/tools/go/packages"
	"strings"
)

var (
	buildFlags = flag.String("build_flags", "", "(package mode) Additional flags for go build.")
)

type packageModeParser struct{}

func (p *packageModeParser) parsePackage(packageName string, structName string) (*entity.TypeMetadata, error) {
	pkg, err := p.loadPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("load package: %w", err)
	}

	typeMetadata, err := p.extractMetadataFromPackage(pkg, structName)
	if err != nil {
		return nil, fmt.Errorf("extract typeMetadata from package: %w", err)
	}

	return typeMetadata, nil
}

func (p *packageModeParser) loadPackage(packageName string) (*packages.Package, error) {
	var buildFlagsSet []string
	if *buildFlags != "" {
		buildFlagsSet = strings.Split(*buildFlags, " ")
	}

	cfg := &packages.Config{
		Mode:       packages.NeedDeps | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedEmbedFiles | packages.LoadSyntax,
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

func (p *packageModeParser) extractMetadataFromPackage(pkg *packages.Package, structName string) (*entity.TypeMetadata, error) {
	obj := pkg.Types.Scope().Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("struct %s does not exist", structName)
	}

	modelIface, err := p.parseStruct(obj)
	if err != nil {
		return nil, err
	}

	return modelIface, nil
}

func (p *packageModeParser) parseStruct(obj types.Object) (*entity.TypeMetadata, error) {
	named, ok := types.Unalias(obj.Type()).(*types.Named)
	if !ok {
		return nil, fmt.Errorf("%s is not an struct. it is a %T", obj.Name(), obj.Type().Underlying())
	}

	strct, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("%s is not an struct. it is a %T", obj.Name(), obj.Type().Underlying())
	}

	mainMetadata := &entity.TypeMetadata{
		Package:     obj.Pkg().Path(),
		TypeName:    obj.Name(),
		TypeKind:    "struct",
		Nodes:       make([]*entity.TypeMetadata, strct.NumFields()),
		Tags:        nil,
		IsReference: false,
	}
	initNodes(mainMetadata)
	parseStructInRecursion(mainMetadata, strct)
	return mainMetadata, nil
}

func initNodes(metadata *entity.TypeMetadata) {
	for i := 0; i < len(metadata.Nodes); i++ {
		metadata.Nodes[i] = &entity.TypeMetadata{}
	}
}

func parseStructInRecursion(typeMetadata *entity.TypeMetadata, strct *types.Struct) {
	for i := 0; i < strct.NumFields(); i++ {
		field := strct.Field(i)

		_, ok := field.Type().(*types.Pointer)
		typeMetadata.Nodes[i].IsReference = ok
		typeMetadata.Nodes[i].FieldName = field.Name()
		typeMetadata.Nodes[i].Package = field.Pkg().Path()

		typeMetadata.Nodes[i].Tags = make(map[string][]string)
		tg := strct.Tag(i)
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
			typeMetadata.Nodes[i].Tags[tmp[0]] = splited
		}

		tmpNamed, ok := field.Type().(*types.Named)
		if ok {
			typeMetadata.Nodes[i].TypeName = tmpNamed.Obj().Name()
			typeMetadata.Nodes[i].Package = tmpNamed.Obj().Pkg().Path()
		}
		tmpPtr, ok := field.Type().(*types.Pointer)
		if ok {
			if tmpNamed, ok = tmpPtr.Elem().(*types.Named); ok {
				typeMetadata.Nodes[i].TypeName = tmpNamed.Obj().Name()
				typeMetadata.Nodes[i].Package = tmpNamed.Obj().Pkg().Path()
			}
		}

		t := field.Type()
		switch v := t.(type) {
		case *types.Basic:
			typeMetadata.Nodes[i].TypeKind = v.String()
			typeMetadata.Nodes[i].TypeName = v.String()
			typeMetadata.Nodes[i].Package = ""
		case *types.Named, *types.Pointer:
			underlying := t.Underlying()
			if val, ok := t.(*types.Pointer); ok {
				underlying = val.Elem().Underlying()
			}
			tmpStruct, ok := underlying.(*types.Struct)
			if ok {
				typeMetadata.Nodes[i].Nodes = make([]*entity.TypeMetadata, tmpStruct.NumFields())
				initNodes(typeMetadata.Nodes[i])
				typeMetadata.Nodes[i].TypeKind = "struct"
				parseStructInRecursion(typeMetadata.Nodes[i], tmpStruct)
			} else {
				typeMetadata.Nodes[i].TypeKind = "undefined"
			}
		default:
			typeMetadata.Nodes[i].TypeKind = "undefined"
		}
	}
}
