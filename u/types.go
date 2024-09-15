package u

import (
	"fmt"
	"go/types"
)

func CopyTypeParams(src []*types.TypeParam) []*types.TypeParam {
	cp := make([]*types.TypeParam, len(src))
	for i := 0; i < len(src); i++ {
		tp := src[i]
		cp[i] = types.NewTypeParam(tp.Obj(), tp.Constraint())
	}
	return cp
}

// RebindVars also anonymize vars.
// vars can be signature.Params() or signature.Results()
func RebindVars(vars *types.Tuple, nameToTypeParam map[string]*types.TypeParam, debug any) (varsBound []*types.Var, err error) {
	varsBound = make([]*types.Var, vars.Len())
	for i := 0; i < vars.Len(); i++ {
		var_ := vars.At(i)
		newParam, err := rebindVar(var_, nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("RebindVars %d: %w", i, err)
		}
		varsBound[i] = newParam
	}
	return varsBound, nil
}

// rebindVar also anonymize var.
func rebindVar(var_ *types.Var, nameToTypeParam map[string]*types.TypeParam, debug any) (varBound *types.Var, err error) {
	if var_ == nil {
		return nil, nil
	}
	fdType, err := rebindType(var_.Type(), nameToTypeParam, debug)
	if err != nil {
		return nil, fmt.Errorf("rebindVar: %w", err)
	}
	// Anonymize early here!
	return types.NewVar(var_.Pos(), var_.Pkg(), "" /*var_.Name()*/, fdType), nil
}

// [typBound] will have the same type as [typ].
// I tried to use [T types.Type](var_ T……), but failed.
// TODO maybe make it OO, hiding nameToTypeParam to field.
func rebindType(typ types.Type, nameToTypeParam map[string]*types.TypeParam, debug any) (typBound types.Type, err error) {
	// refers to types.IdenticalIgnoreTags
	switch typ.(type) {
	case *types.TypeParam:
		tp := typ.(*types.TypeParam)
		// nameToTypeParam is finally used here.
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
		return typ, nil

		//// 先不管这个 case……（后果：struct + 使用类型参数，会不正常 TODO）
		//// internal/fuzz ReadCorpus
		//// func(dir string, types []reflect.Type) ([]struct{Parent string; Path string; Data []byte; Values []any; Generation int; IsSeed bool}, error)
		//// panic: multiple fields with the same name
		//
		//st := typ.(*types.Struct)
		//fields := make([]*types.Var, st.NumFields())
		//tags := make([]string, st.NumFields())
		//for i := 0; i < st.NumFields(); i++ {
		//	fd, err := rebindVar(st.Field(i), nameToTypeParam, debug)
		//	if err != nil {
		//		return nil, fmt.Errorf("rebindType Struct: %w", err)
		//	}
		//	fields[i] = fd
		//	tags[i] = st.Tag(i) // just copy
		//}
		//return types.NewStruct(fields, tags), nil

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
		// But that’s exactly the reason to RebindVars. Unlucky.

		sig := typ.(*types.Signature)
		params, err := RebindVars(sig.Params(), nameToTypeParam, "params of "+sig.String()+" <- "+debug.(string)) // mutual recursion
		if err != nil {
			return nil, fmt.Errorf("rebindType Signature Params: %w", err)
		}
		results, err := RebindVars(sig.Results(), nameToTypeParam, "results of "+sig.String()+" <- "+debug.(string))
		if err != nil {
			return nil, fmt.Errorf("rebindType Signature Results: %w", err)
		}
		recv, err := rebindVar(sig.Recv(), nameToTypeParam, debug)
		if err != nil {
			return nil, fmt.Errorf("rebindType Signature Recv: %w", err)
		}

		// trivial
		// todo use u.CopyTypeParams?
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

// Anonymize 先写一层吧……累了。显而易见的 badcase: lo.Map 模糊搜索不到 TODO recursive
// TODO 融进 RebindVars?
// 2024-09-13 09:42:31 先在 Anonymize 加入「类型参数换名 2」的成分。参考 collect.ParseGenDeclTypeSpecFuncSigs
// RebindVars Anonymize 类型参数换名 三者大概会融合起来：
// 大概把 Anonymize 去掉变量名的操作写进 RebindVars，然后把 Anonymize 的核心实现改为调用 RebindVars。
func Anonymize(sig *types.Signature) *types.Signature {
	//type S[T any] struct {}
	//func (S[T])f[F any](){} // Method cannot have type parameters（指 [F any]）

	//// 先只管 .TypeParams()；TODO .RecvTypeParams()
	typeParams := make([]*types.TypeParam, sig.TypeParams().Len())
	nameToTypeParam := make(map[string]*types.TypeParam)
	for i := 0; i < sig.TypeParams().Len(); i++ {
		tp := sig.TypeParams().At(i)
		newTP := types.NewTypeParam(types.NewTypeName(
			tp.Obj().Pos(), tp.Obj().Pkg(), fmt.Sprintf("_T%d", i), tp.Obj().Type()),
			tp.Constraint())
		typeParams[i] = newTP
		nameToTypeParam[tp.Obj().Name()] = newTP
	}

	params, err := RebindVars(sig.Params(), nameToTypeParam, "params of "+sig.String())
	if err != nil {
		// TODO
		panic(fmt.Errorf("rebinding params of %s: %w", sig, err))
	}
	results, err := RebindVars(sig.Results(), nameToTypeParam, "results of "+sig.String())
	if err != nil {
		// TODO
		panic(fmt.Errorf("rebinding results of %s: %w", sig, err))
	}
	recv, err := rebindVar(sig.Recv(), nameToTypeParam /*TODO 用 sig.RecvTypeParams() 对应物*/, "recv of "+sig.String())

	anonSig := types.NewSignatureType(
		recv,
		CopyTypeParams(TypeParamListToSliceOfTypeParam(sig.RecvTypeParams())),
		typeParams, // with new name!
		types.NewTuple(params...),
		types.NewTuple(results...),
		sig.Variadic(),
	)
	return anonSig

	//return types.NewSignatureType(
	//	anonymizeVar(sig.Recv()),
	//	CopyTypeParams(TypeParamListToSliceOfTypeParam(sig.RecvTypeParams())),
	//	CopyTypeParams(TypeParamListToSliceOfTypeParam(sig.TypeParams())),
	//	types.NewTuple(lo.Map(TupleToSliceOfVar(sig.Params()), anonymizeVarI)...),
	//	types.NewTuple(lo.Map(TupleToSliceOfVar(sig.Results()), anonymizeVarI)...),
	//	sig.Variadic(),
	//)
}
func anonymizeVarI(v *types.Var, _ int) *types.Var {
	return anonymizeVar(v)
}
func anonymizeVar(v *types.Var) *types.Var {
	if v == nil {
		return v
	}
	return types.NewVar(v.Pos(), nil, "", v.Type())
}

func TypeParamListToSliceOfTypeParam(inp *types.TypeParamList) []*types.TypeParam {
	out := make([]*types.TypeParam, inp.Len())
	for i := 0; i < inp.Len(); i++ {
		out[i] = inp.At(i)
	}
	return out
}

func TupleToSliceOfVar(inp *types.Tuple) []*types.Var {
	out := make([]*types.Var, inp.Len())
	for i := 0; i < inp.Len(); i++ {
		out[i] = inp.At(i)
	}
	return out
}

// todo see also https://pkg.go.dev/golang.org/x/exp/slices#Delete
func TupleToSliceOfVarExcept(inp *types.Tuple, except int) []*types.Var {
	out := make([]*types.Var, 0, inp.Len()-1)
	for i := 0; i < inp.Len(); i++ {
		if i != except {
			out = append(out, inp.At(i))
		}
	}
	return out
}
