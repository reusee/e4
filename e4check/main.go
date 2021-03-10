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
		for _, check := range checkObjects {
			for checkIdent, obj := range pkg.TypesInfo.Uses {
				if obj != check {
					continue
				}
				for _, file := range pkg.Syntax {
					path, exact := astutil.PathEnclosingInterval(file, checkIdent.Pos(), checkIdent.Pos())
					if !exact {
						continue
					}
					for _, node := range path {

						var body *ast.BlockStmt
						var signature *ast.FuncType

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
								signature = fnLit.Type
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
								signature = fnDecl.Type
							}
						}

						if body != nil {
							checkOK := false
							// find handle call before check call
							for _, stmt := range body.List {
								if stmt.Pos() > checkIdent.End() {
									// stmt after check
									break
								}
								deferStmt, ok := stmt.(*ast.DeferStmt)
								if !ok {
									// not defer
									continue
								}
								callIdent, ok := deferStmt.Call.Fun.(*ast.Ident)
								if !ok {
									// not call by identifier
									continue
								}
								callObj := pkg.TypesInfo.Uses[callIdent]
								isHandleObject := false
								for _, obj := range handleObjects {
									if callObj == obj {
										isHandleObject = true
										break
									}
								}
								if !isHandleObject {
									// not handle object
									continue
								}
								// check target argument of handle call
								target := deferStmt.Call.Args[0]
								targetExpr, ok := target.(*ast.UnaryExpr)
								if !ok {
									pt("expecting unary expression: %s\n", pkg.Fset.Position(target.Pos()).String())
									return false
								}
								errIdent, ok := targetExpr.X.(*ast.Ident)
								if !ok {
									pt("expecting error identifier: %s\n", pkg.Fset.Position(target.Pos()).String())
									return false
								}
								errObj := pkg.TypesInfo.Uses[errIdent]
								// find def of error object
								for defIdent, defObj := range pkg.TypesInfo.Defs {
									if defObj != errObj {
										continue
									}
									// must define inside signature
									if !(defIdent.Pos() > signature.Results.Pos() && defIdent.End() < signature.Results.End()) {
										pt("should pass error defined at %s\n",
											pkg.Fset.Position(signature.Results.Pos()))
										return false
									}
									checkOK = true
								}

							}

							if !checkOK {
								pt(
									"check %s should be handle inside %s\n",
									pkg.Fset.Position(checkIdent.Pos()).String(),
									pkg.Fset.Position(body.Pos()).String(),
								)
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
