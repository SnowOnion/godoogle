package main

import (
	"context"
	"github.com/SnowOnion/godoogle/server/model"
	"github.com/SnowOnion/godoogle/server/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/samber/lo"
)

func Home(ctx context.Context, c *app.RequestContext) {
	hlog.CtxInfof(ctx, "Home invoked~")

	c.HTML(consts.StatusOK, "search.html", nil)
}

func SearchH(ctx context.Context, c *app.RequestContext) {
	hlog.CtxInfof(ctx, "Search invoked~")

	req := model.SearchReq{}
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxErrorf(ctx, "b")
		//c.HTML(consts.StatusBadRequest, "search.html",)
	}

	if req.Query == "" {
		c.HTML(consts.StatusOK, "search.html", map[string]any{})
		return
	}

	resp, err := service.Search(ctx, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "service.Search err=%s", err)
		c.HTML(consts.StatusInternalServerError, "search.html",
			map[string]any{
				"q":     req.Query,
				"error": "Sorry, something wrong with server. It's not your fault.",
			})
		return
	}

	fellingLucky := lo.TernaryF(len(resp.Result) == 0, lo.Empty[string], func() string { return resp.Result[0].FullName })
	hlog.CtxInfof(ctx, "result len=%d [0]=%s", len(resp.Result), fellingLucky)
	c.HTML(consts.StatusOK, "search.html", map[string]any{
		"q":      req.Query,
		"result": resp.Result,
	})
}

func SearchJ(ctx context.Context, c *app.RequestContext) {
	hlog.CtxInfof(ctx, "Search invoked~")

	req := model.SearchReq{}
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, model.Resp{Code: 400000, Message: "Bad Request"})
		return
	}

	result, err := service.Search(nil, req)
	if err != nil {
		hlog.CtxErrorf(ctx, "service.Search err=%s", err)
		c.JSON(consts.StatusInternalServerError, model.Resp{Code: 500000, Message: "Server Error"})
		return
	}

	c.JSON(consts.StatusOK, model.Resp{Data: result})
	return
}
