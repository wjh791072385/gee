package gee

import (
	"errors"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node // 每个请求方法对应一个前缀树 "GET" "PUT"
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	// "/hello/make" 会解析为 ["","hello","jie"]
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			//要先添加在判断
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//添加路由之前先检测是否有路由冲突
	//只能解决简单冲突，比如/hello/info和/hello/:name, /hello/:name和/hello/:pwd
	//gin框架目前能区分/hello/info和/hello/:name
	n, _ := r.getRoute(method, pattern)
	if n != nil {
		err := errors.New(pattern + " conflicts with existing router " + n.pattern)
		panic(err)
	}

	parts := parsePattern(pattern)
	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 将url中的参数赋值到Context中，比如/hello/:name
		c.Params = params

		// 	这里调用的key注意是加上n.pattern，因为比如路由注册为/hello/:name
		//  若访问/hello/mike,此时c.Path = /hello/mike, 而n.pattern为正确的/hello/:name
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		//c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for idx, part := range parts {
			if part[0] == ':' {
				//比如/hello/:name,当请求为/hello/mike时，name匹配为mike
				params[part[1:]] = searchParts[idx]
			}
			if part[0] == '*' && len(part) > 1 {
				//比如/hello/*name,当请求为/hello/mike/jerry时，name匹配为mike/jerry
				params[part[1:]] = strings.Join(searchParts[idx:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
