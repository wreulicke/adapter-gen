package generator

import (
	"errors"
	"fmt"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

func Load(packageName string) (*packages.Package, error) {
	mode := packages.NeedFiles | packages.NeedSyntax |
		packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo
	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, packageName)
	if err != nil {
		return nil, fmt.Errorf("cannot load package: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, errors.New("cannot load package")
	}

	pkg := pkgs[0]
	if pkg.Errors != nil && len(pkg.Errors) != 0 {
		return nil, pkg.Errors[0]
	}

	return pkg, nil
}

func ShortTypeString(t types.Type) string {
	return types.TypeString(t, func(p *types.Package) string {
		return p.Name()
	})
}

func extractShortPackage(packageName string) string {
	i := strings.LastIndex(packageName, "/")
	return packageName[i+1:]
}

func Methods(t types.Type) []*types.Func {
	methods := []*types.Func{}
	ms := types.NewMethodSet(t)
	for i := 0; i < ms.Len(); i++ {
		m, _ := ms.At(i).Obj().(*types.Func)
		if m != nil && m.Exported() {
			methods = append(methods, m)
		}
	}
	sort.Slice(methods, func(i, j int) bool {
		return strings.Compare(methods[i].Name(), methods[j].Name()) <= 0
	})
	return methods
}

func PointerMethods(t types.Type) []*types.Func {
	if _, isPtr := t.(*types.Pointer); !isPtr {
		return Methods(types.NewPointer(t))
	}
	return []*types.Func{}
}
