package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"

	"github.com/reusee/e4"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

var (
	pt = fmt.Printf
	ce = e4.Check
)

func main() {

	// load
	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedTypesInfo |
				packages.NeedDeps |
				packages.NeedImports |
				packages.NeedFiles |
				packages.NeedSyntax |
				packages.NeedTypes |
				packages.NeedName,
		},
		os.Args[1:]...,
	)
	ce(err)
	if packages.PrintErrors(pkgs) > 0 {
		return
	}

	// get object and types
	var checkFunc types.Object
	var handleFunc types.Object
	packages.Visit(pkgs, func(pkg *packages.Package) bool {
		if pkg.PkgPath == "github.com/reusee/e4" {
			checkFunc = pkg.Types.Scope().Lookup("Check")
			handleFunc = pkg.Types.Scope().Lookup("Handle")
		}
		return true
	}, nil)
	if checkFunc == nil {
		panic("no Check")
	}
	if handleFunc == nil {
		panic("no Handle")
	}

	// check
	packages.Visit(pkgs, func(pkg *packages.Package) bool {

		// get Check and Handle aliases
		findAlias := func(target types.Object) []types.Object {
			var objs []types.Object
			for ident, obj := range pkg.TypesInfo.Uses {
				if obj != target {
					continue
				}
				objs = append(objs, obj)
				for _, file := range pkg.Syntax {
					path, exact := astutil.PathEnclosingInterval(file, ident.Pos(), ident.End())
					if !exact {
						continue
					}
					for _, node := range path {
						valueSpec, ok := node.(*ast.ValueSpec)
						if !ok {
							continue
						}
						for i, value := range valueSpec.Values {
							selExpr, ok := value.(*ast.SelectorExpr)
							if !ok {
								continue
							}
							if selExpr.Sel != ident {
								continue
							}
							id := valueSpec.Names[i]
							obj := pkg.TypesInfo.Defs[id]
							objs = append(objs, obj)
						}
					}
				}
			}
			return objs
		}
		checkObjects := findAlias(checkFunc)
		handleObjects := findAlias(handleFunc)

		// error type
		tv, err := types.Eval(pkg.Fset, pkg.Types, 0, "error(nil)")
		ce(err)
		errorType := tv.Type

		// check usages
		found := make(map[string]bool)
		for _, check := range checkObjects {
			for ident, obj := range pkg.TypesInfo.Uses {
				if obj != check {
					continue
				}
				for _, file := range pkg.Syntax {
					path, exact := astutil.PathEnclosingInterval(file, ident.Pos(), ident.Pos())
					if !exact {
						continue
					}
					for _, node := range path {

						var body *ast.BlockStmt

						if fnLit, ok := node.(*ast.FuncLit); ok {
							// function literal
							fnSig := pkg.TypesInfo.Types[fnLit].Type.(*types.Signature)
							rets := fnSig.Results()
							if rets.Len() < 1 {
								continue
							}
							if types.Identical(
								rets.At(rets.Len()-1).Type(),
								errorType,
							) {
								body = fnLit.Body
							}

						} else if fnDecl, ok := node.(*ast.FuncDecl); ok {
							// function decl
							fnSig := pkg.TypesInfo.Defs[fnDecl.Name].Type().(*types.Signature)
							rets := fnSig.Results()
							if rets.Len() < 1 {
								continue
							}
							if types.Identical(
								rets.At(rets.Len()-1).Type(),
								errorType,
							) {
								body = fnDecl.Body
							}
						}

						if body != nil {
							checkOK := false
							ast.Inspect(body, func(node ast.Node) bool {
								id, ok := node.(*ast.Ident)
								if !ok {
									return true
								}
								obj := pkg.TypesInfo.Uses[id]
								for _, handleObj := range handleObjects {
									if obj == handleObj {
										checkOK = true
										return false
									}
								}
								return true
							})
							if !checkOK {
								pos := pkg.Fset.Position(body.Pos()).String()
								if _, ok := found[pos]; !ok {
									pt("%s\n", pos)
									found[pos] = true
								}
							}
							break
						}

					}
				}
			}
		}

		return true
	}, nil)

}
