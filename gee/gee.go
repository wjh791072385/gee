package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

// Engine implements the interface of ServeHTTP
type Engine struct {
	*RouterGroup  //方便Engine作为最顶层的设计，可以调用RouterGroup的方法
	router        *router
	groups        []*RouterGroup // store all groups
	htmlTemplates *template.Template
	funcMap       template.FuncMap
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

func (engine *Engine) setFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Group part
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

// Use middleware part
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// Static part
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(relativePath, group.prefix)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	// 使用http.Dir是因为Dir实现了http.FileSystem接口
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// addRoute router part
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
	//所有context共享一个engine
	c.engine = engine
	c.handlers = middlewares
	engine.router.handle(c)
}
