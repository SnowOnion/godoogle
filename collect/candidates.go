package collect

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"

	"github.com/SnowOnion/godoogle/u"
)

var (
	FuncDatabase []u.T2
)

func InitFuncDatabase() {
	var err error
	pkgIDs := []string{
		//`golang.org/x/exp/slices`,
		`github.com/samber/lo`,
		`std`,
		`github.com/dominikbraun/graph`,

		//`sort`,
		// when siggraph has no depth limit: |V|=60509; |E|=351739
		// depthTTL=2: |V|=7950; |E|=13607
		// depthTTL=1: |V|=4465; |E|=4187

		//`strconv`,
		//`slices`,
		//`strings`,
		//`maps`,
	}
	FuncDatabase, err = ParseFuncSigsFromPackage(pkgIDs...)
	if err != nil {
		panic("InitFuncDatabase: " + err.Error())
	}

}

// ParseFuncSigs extracts all function signatures declared in [src].
func ParseFuncSigs(src string) (sigs []u.T2, err error) {
	// 创建一个新的 token 文件集，用于语法解析。
	fset := token.NewFileSet()

	// 解析 Go 源代码字符串。
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments) // 不同级别注意一下
	if err != nil {
		return nil, err
	}

	// 创建类型信息配置并初始化一个新的类型检查器。
	config := &types.Config{
		Importer: importer.Default(),
		//types.DefaultImporter(), // gpt-4 幻觉 or 旧版？2024-03-06 01:49:43
	}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}

	// 对 AST 进行类型检查，填充 info。
	_, err = config.Check("", fset, []*ast.File{file}, info)
	if err != nil {
		return nil, err
	}

	// 遍历所有的顶级声明。
	for _, decl := range file.Decls {
		// 确保声明是函数声明。
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// 获取函数定义的对象。
		obj := info.Defs[fn.Name]
		if obj == nil {
			continue
		}

		// 确保对象是函数对象。
		fnObj, ok := obj.(*types.Func)
		if !ok {
			continue
		}

		// 打印函数的签名。
		typ := fnObj.Type()
		sig := typ.(*types.Signature)
		//sig.TypeParams()
		//sig.Params()
		//sig.Results()
		//similarity 开整！

		//fmt.Printf("Function %s: %s\n", fnObj.Name(), sig)
		sigs = append(sigs, u.T2(lo.T2(sig, fnObj)))
		//sig.Params().Len()
	}

	return sigs, nil
}

// ParseFuncSigsFromPackage 解析 path 里 Go 程序所有函数的签名。
func ParseFuncSigsFromDirNaive(path string) (sigs []u.T2, err error) {
	// 创建一个新的 token 文件集，用于语法解析。
	fset := token.NewFileSet()

	// Parse the package's *.go files.
	pkgs, err := parser.ParseDir(fset, path, nil, parser.AllErrors) // 不同级别注意一下；似乎是 .go ∪ filter
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		// Type-check the package.
		// 创建类型信息配置并初始化一个新的类型检查器。
		conf := types.Config{Importer: importer.Default()}
		info := &types.Info{
			Defs: make(map[*ast.Ident]types.Object),
			Uses: make(map[*ast.Ident]types.Object), // 真有用还是幻觉？
		}

		// Create the package and type-check it.
		files := lo.Values(pkg.Files)

		_, err := conf.Check(path, fset, files, info) // path 正确吗 // 这个返回 pkgs 用起来？
		if err != nil {
			return nil, err
		}

		for ind, file := range files {
			fmt.Println("ind", ind, file)

			// 遍历所有的顶级声明。
			for _, decl := range file.Decls {
				// 确保声明是函数声明。
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				// 获取函数定义的对象。
				obj := info.Defs[fn.Name]
				if obj == nil {
					continue
				}

				// 确保对象是函数对象。
				fnObj, ok := obj.(*types.Func)
				if !ok {
					continue
				}

				// 打印函数的签名。
				typ := fnObj.Type()
				sig := typ.(*types.Signature)
				//sig.TypeParams()
				//sig.Params()
				//sig.Results()
				//similarity 开整！

				//fmt.Printf("Function %s: %s\n", fnObj.Name(), sig)
				sigs = append(sigs, u.T2(lo.T2(sig, fnObj)))
				//sig.Params().Len()
			}
		}

	}

	return sigs, nil
}

// ParseFuncSigsFromPackage 解析 ~~~ 里 Go 程序所有函数的签名，正确处理 import
// patterns[i] can be:
// `github.com/samber/lo` (already being go get)
func ParseFuncSigsFromPackage(patterns ...string) (sigs []u.T2, err error) {
	pkgs, _ := LoadDirDoc(patterns...)

	for _, pkg := range pkgs {
		// 在 packages.Package 类型中，
		//Types 字段是一个 *types.Package 对象，包含了类型信息，
		//TypesInfo 字字段是一个 *types.Info 对象，包含了关于包中每个语法节点的详细类型信息。

		// Types provides type information for the package.
		// The NeedTypes LoadMode bit sets this field for packages matching the
		// patterns; type information for dependencies may be missing or incomplete,
		// unless NeedDeps and NeedImports are also set.
		_ = pkg.Types

		info := pkg.TypesInfo // GO BACK TO go/types from golang.org/x/tools/go/packages ~
		files := pkg.Syntax   // GO BACK TO go/ast   from golang.org/x/tools/go/packages ~

		if err != nil {
			return nil, err
		}

		for _, file := range files {
			//fmt.Println("ind", ind, file)

			// 遍历所有的顶级声明。
			for _, decl := range file.Decls {
				// 确保声明是函数声明。
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				// 获取函数定义的对象。
				obj := info.Defs[fn.Name]
				if obj == nil {
					continue
				}

				// 确保对象是函数对象。
				fnObj, ok := obj.(*types.Func)
				if !ok {
					continue
				}

				// 打印函数的签名。
				typ := fnObj.Type()
				sig := typ.(*types.Signature)

				//fmt.Printf("Function %s: %s\n", fnObj.Name(), sig)
				// TODO user option to include/exclude not-exported funcs
				// TODO not-exported receiver may have exported method, but seems not in pkg.go.dev ……
				// p -> q == !p || q
				if recv := sig.Recv(); fnObj.Exported() && (recv == nil || recv.Exported()) {
					//sig = ranking.Anonymize(sig) // TODO !!! be elegant // 循环依赖
					sigs = append(sigs, u.T2(lo.T2(sig, fnObj)))
				}

				//sig.Params().Len()
			}
		}

	}

	return sigs, nil
}

//func LoadDirGPT(path string) {
//	cfg := &packages.Config{
//		Mode:  packages.LoadSyntax, // Load all syntax trees for the packages.
//		Tests: false,
//	}
//
//	// Load the current package using go.mod file.
//	pkgs, err := packages.Load(cfg, ".")
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "failed to load packages: %v", err)
//		os.Exit(1)
//	}
//
//	// Process each package and analyze it as needed.
//	for _, pkg := range pkgs {
//		fmt.Printf("Package: %s\n", pkg.PkgPath)
//
//		// Iterate through the package's syntax trees.
//		for _, file := range pkg.Syntax {
//			// Perform syntax analysis or type checking on the file.
//			// ...
//		}
//	}
//}

// patterns[i] can be:
// `github.com/samber/lo` (already being go get)
func LoadDirDoc(patterns ...string) ([]*packages.Package, error) {
	// Many tools pass their command-line arguments (after any flags)
	// uninterpreted to packages.Load so that it can interpret them
	// according to the conventions of the underlying build system.

	// 如果你使用的是 packages.LoadSyntax 或更严格的模式（例如 packages.LoadTypes 或 packages.LoadAllSyntax），加载的包将包含类型信息。
	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	//cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	//// Print the names of the source files
	//// for each package listed on the command line.
	//for _, pkg := range pkgs {
	//	fmt.Println(pkg.ID, pkg.GoFiles)
	//}

	return pkgs, nil
}

//func ParseGenDeclTypeSpecFuncSigsWithImport(src string) (sigs []*types.Signature, err error) {
//	// 创建一个新的 token 文件集，用于语法解析。
//	fset := token.NewFileSet()
//
//	// 解析 Go 源代码字符串。
//	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
//	if err != nil {
//		return nil, err
//	}
//
//	// 创建类型信息配置并初始化一个新的类型检查器。
//	config := packages.Config{Mode: packages.LoadAllSyntax}
//	info := &types.Info{
//		Types:      make(map[ast.Expr]types.TypeAndValue),
//		Instances:  make(map[*ast.Ident]types.Instance), // 感觉这个对 signature.tparams 比较重要……
//		Defs:       make(map[*ast.Ident]types.Object),
//		Uses:       make(map[*ast.Ident]types.Object),
//		Implicits:  make(map[ast.Node]types.Object),
//		Selections: make(map[*ast.SelectorExpr]*types.Selection),
//		Scopes:     make(map[ast.Node]*types.Scope),
//		InitOrder:  make([]*types.Initializer, 0),
//	}
//
//	// 对 AST 进行类型检查，填充 info。
//	_, err = config.Check("", fset, []*ast.File{file}, info)
//	if err != nil {
//		return nil, err
//	}
//
//	// 遍历所有的顶级声明。
//	for _, decl := range file.Decls {
//		//
//		genDecl, ok := decl.(*ast.GenDecl)
//		if !ok {
//			continue
//		}
//
//		// 遍历一般声明中的规格（Specs）。
//		for _, spec := range genDecl.Specs {
//			// 确保规格是类型规格。
//			typeSpec, ok := spec.(*ast.TypeSpec)
//			if !ok {
//				continue
//			}
//
//			// 查找类型名称对应的对象。
//			obj := info.Defs[typeSpec.Name]
//			if obj == nil {
//				continue
//			}
//
//			//typeObj, ok := obj.(*types.Named)
//			// 确保对象是类型对象。
//
//			typeObj, ok := obj.(*types.TypeName)
//			if !ok {
//				continue
//			}
//
//			//if typeParam, ok := typeObj.Type().Underlying().(*types.TypeParam); ok {
//			//	fmt.Println("typeParam", typeParam)
//			//
//			//}
//
//			// 获取类型对象的类型并断言为 *types.Signature。
//			if signature, ok := typeObj.Type().Underlying().(*types.Signature); ok {
//				//fmt.Printf("Function type %s: %s\n", typeObj.Name(), signature)
//
//				// todo support method
//				recvTypeParams0 := signature.RecvTypeParams()
//				recvTypeParams := make([]*types.TypeParam, recvTypeParams0.Len())
//				for i := 0; i < recvTypeParams0.Len(); i++ {
//					recvTypeParams[i] = recvTypeParams0.At(i)
//				}
//
//				typeParams := make([]*types.TypeParam, 0)
//				nameToTypeParam := make(map[string]*types.TypeParam)
//				// 填补 signature.TypeParams() 的空白
//				//fmt.Println("typeSpec.TypeParams", typeSpec.TypeParams)
//				if typeSpec.TypeParams != nil {
//					for _, tp := range typeSpec.TypeParams.List {
//						//fmt.Println(tp)
//						for _, tpn := range tp.Names {
//							obj := info.Defs[tpn]
//							if obj == nil {
//								continue
//							}
//
//							if tn, ok := obj.Type().(*types.TypeParam); ok {
//								//fmt.Println("tnnnnnnn", tn)
//								//typeParams = append(typeParams, tn) //panic: type parameter bound more than once
//								newTP := types.NewTypeParam(types.NewTypeName(
//									tn.Obj().Pos(), tn.Obj().Pkg(), tn.Obj().Name(), tn.Obj().Type()),
//									tn.Constraint())
//								typeParams = append(typeParams, newTP)
//								nameToTypeParam[newTP.Obj().Name()] = newTP
//							}
//
//						}
//
//					}
//				}
//
//				params, err := rebindVars(signature.Params(), nameToTypeParam, "params of "+signature.String())
//				if err != nil {
//					return sigs, fmt.Errorf("rebinding params of %s: %w", signature, err)
//				}
//				results, err := rebindVars(signature.Results(), nameToTypeParam, "results of "+signature.String())
//				if err != nil {
//					return sigs, fmt.Errorf("rebinding results of %s: %w", signature, err)
//				}
//
//				sigWithTypeParams := types.NewSignatureType(
//					signature.Recv(),
//					recvTypeParams,
//					typeParams, //!
//					types.NewTuple(params...),
//					types.NewTuple(results...),
//					signature.Variadic(),
//				)
//				//panic: type parameter bound more than once
//				//考虑让用户输入带函数名了。至少，泛型的情况。 // 我还得费劲加函数体啊，别了
//
//				sigs = append(sigs, sigWithTypeParams)
//			}
//		}
//
//	}
//
//	return sigs, nil
//}

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
				nameToTypeParam := make(map[string]*types.TypeParam)
				// 填补 signature.TypeParams() 的空白
				//fmt.Println("typeSpec.TypeParams", typeSpec.TypeParams)
				if typeSpec.TypeParams != nil {
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
									tn.Obj().Pos(), tn.Obj().Pkg(), tn.Obj().Name(), tn.Obj().Type()),
									tn.Constraint())
								typeParams = append(typeParams, newTP)
								nameToTypeParam[newTP.Obj().Name()] = newTP
							}

						}

					}
				}

				params, err := rebindVars(signature.Params(), nameToTypeParam, "params of "+signature.String())
				if err != nil {
					return sigs, fmt.Errorf("rebinding params of %s: %w", signature, err)
				}
				results, err := rebindVars(signature.Results(), nameToTypeParam, "results of "+signature.String())
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

// vars can be signature.Params() or signature.Results()
func rebindVars(vars *types.Tuple, nameToTypeParam map[string]*types.TypeParam, debug any) (varsBound []*types.Var, err error) {
	varsBound = make([]*types.Var, vars.Len())
	for i := 0; i < vars.Len(); i++ {
		var_ := vars.At(i)
		newParam, err := rebindVar(var_, nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindVars %d: %w", i, err)
		}
		varsBound[i] = newParam
	}
	return varsBound, nil
}

func rebindVar(var_ *types.Var, nameToTypeParam map[string]*types.TypeParam, debug any) (varBound *types.Var, err error) {
	if var_ == nil {
		return nil, nil
	}
	fdType, err := rebindType(var_.Type(), nameToTypeParam, debug)
	if err != nil {
		return nil, fmt.Errorf("rebindVar: %w", err)
	}
	return types.NewVar(var_.Pos(), var_.Pkg(), var_.Name(), fdType), nil
}

// [typBound] will have the same type as [typ].
// I tried to use [T types.Type](var_ T……), but failed.
// TODO maybe make it OO, hiding nameToTypeParam to field.
func rebindType(typ types.Type, nameToTypeParam map[string]*types.TypeParam, debug any) (typBound types.Type, err error) {
	// refers to types.IdenticalIgnoreTags
	switch typ.(type) {
	case *types.TypeParam:
		tp := typ.(*types.TypeParam)
		newTP, bound := nameToTypeParam[tp.Obj().Name()]
		if !bound {
			// Type checking should not let this happen?
			return nil, fmt.Errorf("[typ=%s] has a parameterized [type=%s] not declared in type typ list", typ, tp.Obj().Name())
		}
		return newTP, nil

	case *types.Array:
		ar := typ.(*types.Array)
		// I need a Maybe Monad
		el, err := rebindType(ar.Elem(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Array: %w", err)
		}
		return types.NewArray(el, ar.Len()), nil

	case *types.Slice:
		sl := typ.(*types.Slice)
		el, err := rebindType(sl.Elem(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Slice: %w", err)
		}
		return types.NewSlice(el), nil

	case *types.Struct: // TODO test
		st := typ.(*types.Struct)
		fields := make([]*types.Var, st.NumFields())
		tags := make([]string, st.NumFields())
		for i := 0; i < st.NumFields(); i++ {
			fd, err := rebindVar(st.Field(i), nameToTypeParam, debug)
			if err != nil {
				return nil, fmt.Errorf("rebindType Struct: %w", err)
			}
			fields[i] = fd
			tags[i] = st.Tag(i) // just copy
		}
		return types.NewStruct(fields, tags), nil

	case *types.Pointer:
		pt := typ.(*types.Pointer)
		el, err := rebindType(pt.Elem(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Pointer: %w", err)
		}
		return types.NewPointer(el), nil

	case *types.Tuple: // “simplified struct”
		tp := typ.(*types.Tuple)
		vars := make([]*types.Var, tp.Len())
		for i := 0; i < tp.Len(); i++ {
			var_, err := rebindVar(tp.At(i), nameToTypeParam, debug)
			if err != nil {
				return nil, fmt.Errorf("rebindType Tuple: %w", err)
			}
			vars[i] = var_
		}
		return types.NewTuple(vars...), nil

	case *types.Signature:
		// spec#FunctionType can not have type params, so I don’t need to handle nested type param scope. Lucky.
		// But that’s exactly the reason to rebindVars. Unlucky.

		sig := typ.(*types.Signature)
		params, err := rebindVars(sig.Params(), nameToTypeParam, "params of "+sig.String()+" <- "+debug.(string)) // mutual recursion
		if err != nil {
			return nil, fmt.Errorf("rebindType Signature Params: %w", err)
		}
		results, err := rebindVars(sig.Results(), nameToTypeParam, "results of "+sig.String()+" <- "+debug.(string))
		if err != nil {
			return nil, fmt.Errorf("rebindType Signature Results: %w", err)
		}
		recv, err := rebindVar(sig.Recv(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Signature Recv: %w", err)
		}

		// trivial
		recvTypeParams0 := sig.RecvTypeParams()
		recvTypeParams := make([]*types.TypeParam, recvTypeParams0.Len())
		for i := 0; i < recvTypeParams0.Len(); i++ {
			recvTypeParams[i] = recvTypeParams0.At(i)
		}
		typeParams0 := sig.TypeParams()
		typeParams := make([]*types.TypeParam, typeParams0.Len())
		for i := 0; i < typeParams0.Len(); i++ {
			typeParams[i] = typeParams0.At(i)
		}

		return types.NewSignatureType(recv, recvTypeParams, typeParams, types.NewTuple(params...), types.NewTuple(results...), sig.Variadic()), nil

	case *types.Union: // TODO: is there?
		un := typ.(*types.Union)
		terms := make([]*types.Term, un.Len())
		for i := 0; i < un.Len(); i++ {
			newType, err := rebindType(un.Term(i).Type(), nameToTypeParam, debug)
			if err != nil {
				return nil, fmt.Errorf("rebindType Union Recv: %w", err)
			}
			terms[i] = types.NewTerm(un.Term(i).Tilde(), newType)
		}
		return types.NewUnion(terms), nil

	case *types.Interface:
		// TODO impl
		//iface:=typ.(*types.Interface)
		return typ, nil

	case *types.Map:
		mp := typ.(*types.Map)
		key, err := rebindType(mp.Key(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Map Key: %w", err)
		}
		el, err := rebindType(mp.Elem(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Map Elem: %w", err)
		}
		return types.NewMap(key, el), nil

	case *types.Chan:
		ch := typ.(*types.Chan)
		// I need a Maybe Monad
		el, err := rebindType(ch.Elem(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Chan: %w", err)
		}
		return types.NewChan(ch.Dir(), el), nil

	case *types.Named: // func(A[T])
		nm := typ.(*types.Named)
		nm.TypeParams() // does it mean type args here?
		// TODO impl
		return typ, nil

	default: // nil, *types.Basic
		return typ, nil
	}
}

// ParseFuncSigsFromPackage 解析 ~~~ 里 Go 程序所有函数的签名，正确处理 import
// patterns[i] can be:
// `github.com/samber/lo` (already being go get)
func ParseGenDeclTypeSpecFuncSigsUsingPackages(patterns ...string) (sigs []u.T2, err error) {
	pkgs, _ := LoadDirDoc(patterns...)

	for _, pkg := range pkgs {
		// 在 packages.Package 类型中，
		//Types 字段是一个 *types.Package 对象，包含了类型信息，
		//TypesInfo 字字段是一个 *types.Info 对象，包含了关于包中每个语法节点的详细类型信息。

		// Types provides type information for the package.
		// The NeedTypes LoadMode bit sets this field for packages matching the
		// patterns; type information for dependencies may be missing or incomplete,
		// unless NeedDeps and NeedImports are also set.
		_ = pkg.Types

		info := pkg.TypesInfo
		files := pkg.Syntax // 从 golang.org/x/tools/go/packages 回到了 go/ast ~

		if err != nil {
			return nil, err
		}

		for ind, file := range files {
			ind = ind
			//fmt.Println("ind", ind, file)

			// 遍历所有的顶级声明。
			for _, decl := range file.Decls {
				// 确保声明是函数声明。
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				// 获取函数定义的对象。
				obj := info.Defs[fn.Name]
				if obj == nil {
					continue
				}

				// 确保对象是函数对象。
				fnObj, ok := obj.(*types.Func)
				if !ok {
					continue
				}

				// 打印函数的签名。
				typ := fnObj.Type()
				sig := typ.(*types.Signature)
				//sig.TypeParams()
				//sig.Params()
				//sig.Results()
				//similarity 开整！

				//fmt.Printf("Function %s: %s\n", fnObj.Name(), sig)
				sigs = append(sigs, u.T2(lo.T2(sig, fnObj)))
				//sig.Params().Len()
			}
		}

	}

	return sigs, nil
}
