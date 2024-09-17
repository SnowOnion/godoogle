package u

import (
	"go/types"

	"github.com/samber/lo"
)

type T2 lo.Tuple2[*types.Signature, *types.Func]

//type T2T lo.Tuple2[*types.Signature, *ast.GenDecl]

// for better debugging
func (t T2) String() string {
	return t.B.FullName() + " :: " + t.A.String()
}
