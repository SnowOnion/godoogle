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
		"std",

		//// Copy to imports.go
		//// Some has few functions but lots of methods.

		// basic data types & data structures
		"github.com/samber/lo",
		"github.com/samber/mo",
		"github.com/thoas/go-funk",
		"github.com/dominikbraun/graph",
		// Date and Time
		"github.com/jinzhu/now",
		// concurrent & parallel
		"github.com/sourcegraph/conc",
		// HTTP
		"github.com/go-resty/resty/v2",
		// ORM
		"gorm.io/gorm",
		"xorm.io/xorm",
		// Testing
		"github.com/stretchr/testify/assert",

		//"sort",
		// when siggraph has no depth limit: |V|=60509; |E|=351739
		// depthTTL=2: |V|=7950; |E|=13607
		// depthTTL=1: |V|=4465; |E|=4187

		//"strconv",
		//"slices",
		//"strings",
		//"maps",

		//`golang.org/x/exp/slices`,
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
					sig = u.Anonymize(sig) // TODO 不循环依赖了；去掉后续冗余的 Anonymize
					sigs = append(sigs, u.T2(lo.T2(sig, fnObj)))
				}

				//sig.Params().Len()
			}
		}

	}

	return sigs, nil
}

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

// TODO use or del
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
