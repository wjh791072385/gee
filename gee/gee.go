package gee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

// Engine implements the interface of ServeHTTP
type Engine struct {
	*RouterGroup //方便Engine作为最顶层的设计，可以调用RouterGroup的方法
	router       *router
	groups       []*RouterGroup // store all groups
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	//初始group，即相当于全局group
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	//all groups share the same Engine instance !
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	//路由是分组前缀 + 后缀
	newPattern := group.prefix + pattern

	//全局group共享一个engine实例，因此可以这样调用
	group.engine.router.addRoute(method, newPattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
