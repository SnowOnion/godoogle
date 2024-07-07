package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/requestid"
	"github.com/sirupsen/logrus"

	"github.com/SnowOnion/godoogle/collect"
	"github.com/SnowOnion/godoogle/ranking"
)

func main() {
	// https://github.com/cloudwego/hertz-examples/blob/main/middleware/requestid/log_with_hertzlogrus/main.go
	logger := hertzlogrus.NewLogger(hertzlogrus.WithHook(&RequestIdHook{}))
	hlog.SetLogger(logger)
	hlog.Info("Logger set!")

	// server.Default() creates a Hertz with recovery middleware.
	// If you need a pure hertz, you can use server.New()
	h := server.Default(server.WithHostPorts("[::]:8888"))
	h.Use(requestid.New())

	pprof.Register(h)

	h.LoadHTMLGlob("res/views/*")
	h.Static("/", "./res/assets")
	h.GET("/", Home)

	h.GET("/search", SearchH)

	// todo elegantly initialize; OO
	hlog.Info("Start initializing FuncDatabase and ranker...")
	collect.InitFuncDatabase()
	ranking.DefaultRanker = ranking.NewHooglyRanker(collect.FuncDatabase) // = =、TODO be elegant!
	hlog.Info("End initializing FuncDatabase and ranker!")

	hlog.Info("Start serving...")
	h.Spin()
}

type RequestIdHook struct{}

func (h *RequestIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *RequestIdHook) Fire(e *logrus.Entry) error {
	ctx := e.Context
	if ctx == nil {
		return nil
	}
	value := ctx.Value("X-Request-ID")
	if value != nil {
		e.Data["request_id"] = value
	}
	return nil
}
