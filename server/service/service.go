package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/samber/lo"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/ranking"
	"github.com/SnowOnion/godoogle/server/model"
	"github.com/SnowOnion/godoogle/u"
)

func Search(ctx context.Context, req model.SearchReq) (model.SearchResp, error) {
	hlog.CtxInfof(ctx, "Search %s", req.Query)
	inpSig, err := u.Dummy(req.Query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Dummy err=%s", err)
		var err2 error
		inpSig, err2 = u.Dummy2(req.Query)
		if err2 != nil {
			hlog.CtxErrorf(ctx, "Dummy2 err=%s", err2)
			return model.SearchResp{}, fmt.Errorf("error parsing query: %w", errors.Join(err, err2))
		}
	}
	hlog.CtxInfof(ctx, "DummyQ %s", inpSig.String())

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
	}, nil
}
