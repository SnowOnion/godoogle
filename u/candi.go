package u

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
)

// Dummy rawQuery e.g.
// "func(rune) bool"
// "[T, R any, K comparable] func (collection []T, iteratee func(item T, index int) R) []R"
func Dummy(rawQuery string) (*types.Signature, error) {
	augmentedQuery := `package dummy
import (
	//"golang.org/x/exp/constraints"
	"sort"
	"time"
)
type dummy ` + rawQuery
	sigs, err := ParseGenDeclTypeSpecFuncSigs(augmentedQuery)
	if err != nil {
		return nil, err
	}
	if len(sigs) == 0 {
		return nil, errors.New("no type signature in augmentedQuery")
	}
	return sigs[0], nil
}

// Naive
func ParseGenDeclTypeSpecFuncSigs(src string) (sigs []*types.Signature, err error) {
	// 创建一个新的 token 文件集，用于语法解析。
	fset := token.NewFileSet()

	// 解析 Go 源代码字符串。
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// 创建类型信息配置并初始化一个新的类型检查器。
	config := &types.Config{
		Importer:                 importer.Default(),
		DisableUnusedImportCheck: true,
		IgnoreFuncBodies:         true,

		//types.DefaultImporter(), // gpt-4 幻觉 or 旧版？2024-03-06 01:49:43
	}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Instances:  make(map[*ast.Ident]types.Instance), // 感觉这个对 signature.tparams 比较重要……
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
		InitOrder:  make([]*types.Initializer, 0),
	}

	// 对 AST 进行类型检查，填充 info。
	_, err = config.Check("", fset, []*ast.File{file}, info)
	if err != nil {
		return nil, err
	}

	// 遍历所有的顶级声明。
	for _, decl := range file.Decls {
		//
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		// 遍历一般声明中的规格（Specs）。
		for _, spec := range genDecl.Specs {
			// 确保规格是类型规格。
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 查找类型名称对应的对象。
			obj := info.Defs[typeSpec.Name]
			if obj == nil {
				continue
			}

			//typeObj, ok := obj.(*types.Named)
			// 确保对象是类型对象。

			typeObj, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}

			//if typeParam, ok := typeObj.Type().Underlying().(*types.TypeParam); ok {
			//	fmt.Println("typeParam", typeParam)
			//
			//}

			// 获取类型对象的类型并断言为 *types.Signature。
			if signature, ok := typeObj.Type().Underlying().(*types.Signature); ok {
				//fmt.Printf("Function type %s: %s\n", typeObj.Name(), signature)

				// todo support method
				recvTypeParams0 := signature.RecvTypeParams()
				recvTypeParams := make([]*types.TypeParam, recvTypeParams0.Len())
				for i := 0; i < recvTypeParams0.Len(); i++ {
					recvTypeParams[i] = recvTypeParams0.At(i)
				}

				typeParams := make([]*types.TypeParam, 0)
				nameToTypeParam := make(map[string]*types.TypeParam) // oldName -> TypeParamWithNewName
				// 填补 signature.TypeParams() 的空白
				// 大概也可以提取到「for _, decl := range file.Decls」外面，复用 u.Anonymize

				//fmt.Println("typeSpec.TypeParams", typeSpec.TypeParams)
				if typeSpec.TypeParams != nil {
					typeParamID := 0 // 类型参数换名 1
					for _, tp := range typeSpec.TypeParams.List {
						//fmt.Println(tp)
						for _, tpn := range tp.Names {
							obj := info.Defs[tpn]
							if obj == nil {
								continue
							}

							if tn, ok := obj.Type().(*types.TypeParam); ok {
								//fmt.Println("tnnnnnnn", tn)
								//typeParams = append(typeParams, tn) //panic: type parameter bound more than once
								newTP := types.NewTypeParam(types.NewTypeName(
									tn.Obj().Pos(), tn.Obj().Pkg() /*tn.Obj().Name()*/, fmt.Sprintf("_T%d", typeParamID), tn.Obj().Type()),
									tn.Constraint())
								typeParams = append(typeParams, newTP)
								//nameToTypeParam[newTP.Obj().Name()] = newTP
								nameToTypeParam[tn.Obj().Name()] = newTP

								typeParamID++
							}

						}

					}
				}

				params, err := RebindVars(signature.Params(), nameToTypeParam, "params of "+signature.String())
				if err != nil {
					return sigs, fmt.Errorf("rebinding params of %s: %w", signature, err)
				}
				results, err := RebindVars(signature.Results(), nameToTypeParam, "results of "+signature.String())
				if err != nil {
					return sigs, fmt.Errorf("rebinding results of %s: %w", signature, err)
				}

				sigWithTypeParams := types.NewSignatureType(
					signature.Recv(),
					recvTypeParams,
					typeParams, //!
					types.NewTuple(params...),
					types.NewTuple(results...),
					signature.Variadic(),
				)
				//panic: type parameter bound more than once
				//考虑让用户输入带函数名了。至少，泛型的情况。 // 我还得费劲加函数体啊，别了

				sigs = append(sigs, sigWithTypeParams)
			}
		}

	}

	return sigs, nil
}
