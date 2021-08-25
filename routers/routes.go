package routers

import (
	"boframe/controller/anchor"
	"boframe/controller/user"
	"cc"
	"net/http"
	"net/http/pprof"
)

func Setup() *cc.Engine {
	r := cc.New()

	r.GET("/", func(context *cc.Context) {
		context.String(http.StatusOK, "ok")
	})

	r.POST("/user/present", user.SendPresent)

	r.GET("/anchor/present/log", anchor.PresentLog)

	r.GET("/anchor/ranking", anchor.Ranking)

	r.PProf("/debug/pprof/", pprof.Index)
	r.PProf("/debug/pprof/heap", pprof.Index)
	r.PProf("/debug/pprof/allocs", pprof.Index)
	r.PProf("/debug/pprof/block", pprof.Index)
	r.PProf("/debug/pprof/goroutine", pprof.Index)
	r.PProf("/debug/pprof/mutex", pprof.Index)
	r.PProf("/debug/pprof/threadcreate", pprof.Index)
	r.PProf("/debug/pprof/cmdline", pprof.Cmdline)
	r.PProf("/debug/pprof/profile", pprof.Profile)
	r.PProf("/debug/pprof/symbol", pprof.Symbol)
	r.PProf("/debug/pprof/trace", pprof.Trace)

	return r
}
