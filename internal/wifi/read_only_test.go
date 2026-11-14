package wifi

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestPackageHasNoWriteDependencies(t *testing.T) {
	paths, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		t.Run(path, func(t *testing.T) {
			fset := token.NewFileSet()
			imports, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatal(err)
			}
			deviceNames := map[string]bool{}
			for _, spec := range imports.Imports {
				importPath, err := strconv.Unquote(spec.Path.Value)
				if err != nil {
					t.Fatal(err)
				}
				switch importPath {
				case "github.com/blackdark-org/travo-greenfield/internal/txn",
					"github.com/blackdark-org/travo-greenfield/internal/uciwrite":
					t.Errorf("read-only package imports writer %s", importPath)
				case "github.com/blackdark-org/travo-greenfield/internal/device":
					name := "device"
					if spec.Name != nil {
						name = spec.Name.Name
					}
					if name == "." {
						t.Error("dot-importing device can expose UCI writers")
					} else {
						deviceNames[name] = true
					}
				}
			}

			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				t.Fatal(err)
			}
			ast.Inspect(file, func(node ast.Node) bool {
				selector, ok := node.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				qualifier, ok := selector.X.(*ast.Ident)
				if ok && deviceNames[qualifier.Name] {
					switch selector.Sel.Name {
					case "Device", "UCI", "Reloader":
						t.Errorf("read-only package uses writer device.%s", selector.Sel.Name)
					}
				}
				return true
			})
		})
	}
}
