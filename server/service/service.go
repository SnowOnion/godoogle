package service

import (
	"context"
	"fmt"
	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/server/model"
	"github.com/SnowOnion/godoogle/u"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/samber/lo"
)

func Search(ctx context.Context, req model.SearchReq) (model.SearchResp, error) {
	hlog.CtxInfof(ctx, "Search %s", req) // todo security concern
	inpSig, err := lo.T2(collect.Dummy(req.Query)).Unpack()
	if err != nil {
		hlog.CtxErrorf(ctx, "Dummy err=%s", err)
		// TODO on error, print request_id to user
		return model.SearchResp{}, fmt.Errorf("error parsing query")
	}

	result := lo.Map(collect.NaiveRanker.Rank(inpSig, collect.FuncDatabase),
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
				Signature: sigDecl.A.String(),
			}
		})

	return model.SearchResp{
		Result: result,
	}, err
}
