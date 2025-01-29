package ast

import (
	"errors"
	"fmt"
	"github.com/paulrozhkin/jsonschema/pkg/entity"
	"golang.org/x/mod/modfile"
	"log"
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
			log.Fatalf("Get current directory failed: %v", err)
		}
		p.packageName, err = packageNameOfDir(dir)
		if err != nil {
			log.Fatalf("Parse package name failed: %v", err)
		}
	}
	parser := packageModeParser{}
	pkg, err := parser.parsePackage(p.packageName, p.typeName)
	return pkg, err
}

var errOutsideGoPath = errors.New("source directory is outside GOPATH")

// packageNameOfDir get package import path via dir
func packageNameOfDir(srcDir string) (string, error) {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
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
