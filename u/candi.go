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

// Dummy1 rawQuery e.g.
// "func(rune) bool"
// "[T, R any, K comparable] func (collection []T, iteratee func(item T, index int) R) []R"
func Dummy1(rawQuery string) (*types.Signature, error) {
	augmentedQuery := `package dummy
import (
	//"golang.org/x/exp/constraints"
	"sort"
	"time"
	"cmp"
	"sync"
	"sync/atomic"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"unicode"
	"unsafe"
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

func Dummy2(rawQuery string) (*types.Signature, error) {
	augmentedQuery := `package dummy
type dummy ` + rawQuery

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", augmentedQuery, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	spec := file.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec)

	var typeParams []*types.TypeParam
	if spec.TypeParams != nil {
		typeParamFields := spec.TypeParams.List
		for i, typeParam := range typeParamFields { // [「a,b any」,「c comparable」]
			i = i
			constraintName := ""
			switch typeParam.Type.(type) {
			case *ast.SelectorExpr:
				selector := typeParam.Type.(*ast.SelectorExpr)
				x := selector.X.(*ast.Ident).Name // TODO 允许没有？
				sel := selector.Sel.Name
				constraintName = x + "." + sel // e.g. cmp.Ordered
			case *ast.Ident:
				ident := typeParam.Type.(*ast.Ident)
				constraintName = ident.String()
				// TODO 其他？~[]T
			}

			constraint := types.NewNamed(types.NewTypeName(token.NoPos, nil, constraintName, types.NewInterfaceType(nil, nil)),
				types.NewInterfaceType(nil, nil),
				nil)

			for j, tpName := range typeParam.Names { // [「a」,「b」 any]
				j = j
				//tpName.Name // 用来 anonymize 时排序 TODO 暂且塞进去
				typesTypeParam := types.NewTypeParam(
					types.NewTypeName(token.NoPos, nil, tpName.Name, nil /*这里？TODO*/),
					constraint)
				typeParams = append(typeParams, typesTypeParam)
			}
		}
	}

	funcType := spec.Type.(*ast.FuncType) // TODO 防
	var params []*types.Var
	if funcType.Params != nil {
		params = astFieldListToTypesVars(funcType.Params.List, augmentedQuery)
	}
	var results []*types.Var
	if funcType.Results != nil {
		results = astFieldListToTypesVars(funcType.Results.List, augmentedQuery)
	}

	//funcType.Results.List

	sig := types.NewSignatureType(
		nil,        // TODO
		nil,        // TODO
		typeParams, // 相同 type 的自动扎堆？！……
		types.NewTuple(params...),
		types.NewTuple(results...),
		false, // TODO 看 Ellipsis type?
	)
	//anonSig := Anonymize(sig) // 弄好关联再匿（可能也不需要了）
	return sig, nil

	////// 先只管 .TypeParams()；TODO .RecvTypeParams()
	//typeParams := make([]*types.TypeParam, sig.TypeParams().Len())
	//nameToTypeParam := make(map[string]*types.TypeParam)
	//for i := 0; i < sig.TypeParams().Len(); i++ {
	//	tp := sig.TypeParams().At(i)
	//	newTP := types.NewTypeParam(types.NewTypeName(
	//		tp.Obj().Pos(), tp.Obj().Pkg(), fmt.Sprintf("_T%d", i), tp.Obj().Type()),
	//		tp.Constraint())
	//	typeParams[i] = newTP
	//	nameToTypeParam[tp.Obj().Name()] = newTP
	//}
	//
	//params, err := RebindVars(sig.Params(), nameToTypeParam, "params of "+sig.String())
	//if err != nil {
	//	// TODO
	//	panic(fmt.Errorf("rebinding params of %s: %w", sig, err))
	//}
	//results, err := RebindVars(sig.Results(), nameToTypeParam, "results of "+sig.String())
	//if err != nil {
	//	// TODO
	//	panic(fmt.Errorf("rebinding results of %s: %w", sig, err))
	//}
	//recv, err := rebindVar(sig.Recv(), nameToTypeParam /*TODO 用 sig.RecvTypeParams() 对应物*/, "recv of "+sig.String())
	//
	//anonSig := types.NewSignatureType(
	//	recv,
	//	CopyTypeParams(TypeParamListToSliceOfTypeParam(sig.RecvTypeParams())),
	//	typeParams, // with new name!
	//	types.NewTuple(params...),
	//	types.NewTuple(results...),
	//	sig.Variadic(),
	//)
	//return anonSig

}

// astFieldListToTypesVars not only ast.Field -> types.Var,
// but also flatten: (a,b T, c F) -> (a T, b T, c F)
func astFieldListToTypesVars(list []*ast.Field, src string) []*types.Var {
	var vars []*types.Var
	for _, astField := range list {
		s := astField.Type.Pos()
		e := astField.Type.End()
		tn := src[s-1 : e-1]

		////1. 类型名，能不能从 pos 读源码子串，不想分类讨论。这玩意 trim 之后没啥要规范化的吧，除了 named 跟着类型变量换名。
		////	啊好吧，……还真得手工逐类型构造 underlying 啊……？不对，我只需要规范字符串…… 但是 类型变量换名 这个理由就够了。TODO 知道咋做，就是麻烦。
		//// 可以分批做！！dummy -> dummy2 兜底，功能不会劣化。
		////2. 恐怕要现场看 type param 然后绑上
		//var typeName string
		//typeName = typeName
		//// TODO 已知标准库也有直接 struct 而不是 named 的……； 以及 github.com/samber/lo.Async0 func(f func()) <-chan struct{}
		//switch astField.Type.(type) {
		//case *ast.SelectorExpr:
		//	selector := astField.Type.(*ast.SelectorExpr)
		//	x := selector.X.(*ast.Ident).Name
		//	sel := selector.Sel.Name
		//	typeName = x + "." + sel // e.g. time.Time
		//case *ast.Ident:
		//	typeName = astField.Type.(*ast.Ident).Name
		//case *ast.StarExpr:
		//	typeName = "*" + astField.Type.(*ast.StarExpr).X.(*ast.Ident).Name
		//case *ast.ArrayType:
		//	typeName = "[]" + astField.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name
		//	// TODO 含 slice。
		//case *ast.FuncType:
		//case *ast.InterfaceType:
		//case *ast.StructType:
		//
		//}
		//

		if astField.Names == nil { // 可以没有！没有就是 1 个。copy obj or copy code !
			typesTypeName := types.NewTypeName(token.NoPos, nil, tn, nil)
			vars = append(vars, types.NewVar(token.NoPos, nil, "", /*直接匿名*/
				types.NewNamed(typesTypeName, types.NewInterfaceType(nil, nil) /*?TODO*/, nil)))
		}
		for /*j, ident*/ _, _ = range astField.Names {
			typesTypeName := types.NewTypeName(token.NoPos, nil, tn, nil)
			vars = append(vars, types.NewVar(token.NoPos, nil, "", /*直接匿名*/
				types.NewNamed(typesTypeName, types.NewInterfaceType(nil, nil) /*?TODO*/, nil)))
		}
	}

	return vars
}

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
		//Error: // TODO try this hook? 能忽略部分错误（不认识的 type 和 constraint）继续么？

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
	pkg, err := config.Check("", fset, []*ast.File{file}, info)
	pkg = pkg

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

// TODO del this draft
// // 只 parse，还需要补 return 么？
// z 改成防止与用户输入撞名的怪东西
func z[T any]() (v T) { return v }

// func f()(int,error)
// func f()(int,error){}
func f[F comparable]() (int, error) { return z[int](), z[error]() }

// func [F comparable]()(int,error)
func f___[F comparable]() (int, error) { return z[int](), z[error]() }
