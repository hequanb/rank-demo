package cc

import (
	"fmt"
	"log"
	"net/http"
)

// Engine httpserver入口
type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{map[string]HandlerFunc{}}
}

type HandlerFunc func(context *Context)

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	key := keyBuild(req.Method, req.URL.Path)
	handleFunc, ok := engine.router[key]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "404 Resource not found:%s \n", req.URL)
		return
	}
	c := newContext(writer, req)
	handleFunc(c)
}

func (engine *Engine) addRoute(method, path string, handle HandlerFunc) {
	key := keyBuild(method, path)
	log.Printf("serving method: %s, %s \n", method, path)
	engine.router[key] = handle
}

func keyBuild(method, path string) (key string) {
	key = method + "==" + path
	return
}

func (engine *Engine) GET(path string, handle HandlerFunc) {
	engine.addRoute("GET", path, handle)
}

func (engine *Engine) POST(path string, handle HandlerFunc) {
	engine.addRoute("POST", path, handle)
}

func (engine *Engine) Serve(address string) error {
	return http.ListenAndServe(address, engine)
}

type PProfFunc func(w http.ResponseWriter, r *http.Request)

func (engine *Engine) PProf(path string, handle PProfFunc) {
	engine.addRoute("GET", path, func(context *Context) {
		handle(context.Writer, context.Req)
	})
}
