## 简介
***
搭建轻量级go Web框架gee,参考gin框架实现 

实现功能如下：
* 上下文Context
* 前缀树路由
* 分组控制
* 中间件
* 静态资源服务、HTML模板渲染
* 错误恢复

## Installation
***
```
$ go get -u github.com/wjh791072385/gee
```

## Quick start
***
```
package main

import (
	"net/http"

	"github.com/wjh791072385/gee"
)

func main() {
	r := gee.Default()
	r.GET("/hello", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"msg": "hello world",
		})
	})
	r.Run("localhost:8888")
}

```

## API Examples
***
**Using GET, POST, PUT, DELETE**
```
func main() {
	r := gee.Default()
	
	r.GET("/someGet", getting)
	r.POST("/somePost", posting)
	r.PUT("/somePut", putting)
	r.DELETE("/someDelete", deleting)

	r.Run("localhost:8888")
}

```

**Parameters in path**
```
func main() {
	r := gee.Default()
	
	r.GET("/hello/:name", func(c *Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Param("name"), c.Path)
	})

	r.GET("/ass/*filepath", func(c *Context) {
		c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
	})

	r.Run("localhost:8888")
}
	
```

**Querystring parameters**
```
func main() {
	r := gee.Default()
	
	r.GET("/hello/:name", func(c *Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Param("name"), c.Path)
	})

	r.GET("/ass/*filepath", func(c *Context) {
		c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
	})

	r.Run("localhost:8888")
}
```
**Querystring parameters**
```
func main() {
	r := gee.Default()
	
	r.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *Context) {
		c.JSON(http.StatusOK, H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run("localhost:8888")
}
```

**Grouping routes**
```
func main() {
	r := gee.Default()
	
	v1 := r.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	r.Run("localhost:8888")
}
```

**Blank Gee without middleware by default**
Use
```
r := gin.New()
```
instead of
```
// Default With the Logger and Recovery middleware already attached
r := gin.Default()
```

**Using middleware**
```
func main() {
	r := gee.New()
	
	r.Use(MyMiddleware())
	r.GET("/hello", func(c *gee.Context) {
		log.Println("execute...")
		c.JSON(http.StatusOK, gee.H{
			"msg" : "world",
		})
	})

	r.Run("localhost:8888")
}
```
**Load Static Source and Use HTML Templates**
```
func main() {
	r := gee.Default()
	
	//加载静态资源,访问localhost:8888/assets/file1.txt相当于访问./resource/file1.txt
	r.Static("/assets", "./resource")
	r.LoadHTMLGlob("templates/*")
	
	r.GET("/getHtml", func(c *Context) {
		c.HTML(http.StatusOK, "tes.tmpl", H{
			"msg":    "this is a template",
			"status": "ok",
		})
	})

	r.Run("localhost:8888")
}
```
