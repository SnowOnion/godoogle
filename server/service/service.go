package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/samber/lo"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/ranking"
	"github.com/SnowOnion/godoogle/server/model"
	"github.com/SnowOnion/godoogle/u"
)

func Search(ctx context.Context, req model.SearchReq) (model.SearchResp, error) {
	hlog.CtxInfof(ctx, "Search %s", req) // todo security concern
	inpSig, err := lo.T2(u.Dummy(req.Query)).Unpack()
	if err != nil {
		hlog.CtxErrorf(ctx, "Dummy err=%s", err)
		return model.SearchResp{}, fmt.Errorf("error parsing query")
	}

	ranked := ranking.DefaultRanker.Rank(inpSig, collect.FuncDatabase)
	result := lo.Map(ranked,
		func(sigDecl u.T2, ind int) model.ResultItem {
			name := sigDecl.B.Name()
			pkg := sigDecl.B.Pkg().Path()

			url := fmt.Sprintf("https://pkg.go.dev/%s#%s", pkg, name)
			if recv := sigDecl.A.Recv(); recv != nil {
				url = fmt.Sprintf("https://pkg.go.dev/%s#%s.%s", pkg, recv.Type().String(), name)
			}
			return model.ResultItem{
				ID:        ind,
				IDDisplay: fmt.Sprintf("%4d", ind),
				Name:      name,
				FullName:  sigDecl.B.FullName(),
				Pkg:       pkg,
				URL:       url,
				Signature: sigDecl.B.Signature().String(), // show signature before param anonymizing and type param renaming
			}
		})

	return model.SearchResp{
		Result: result,
	}, err
}
