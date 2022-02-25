package goWebGee

import (
	"log"
	"net/http"
	"testing"
)

func GroupMid() HandlerFunc {
	return func(c *Context) {
		log.Println("group middleware")
		c.Next()
	}
}

func TestDefault(t *testing.T) {
	//r := New()
	// Default方法默认使用Logger和Recovery中间件
	r := Default()

	//加载静态资源,访问localhost:8888/assets/file1.txt相当于访问./resource/file1.txt
	r.Static("/assets", "./resource")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *Context) {
		c.String(http.StatusOK, "<h1>This is Gee</h1>")
	})

	//定义组内中间件
	v1 := r.Group("/v1")
	v1.Use(GroupMid())
	{
		v1.GET("/", func(c *Context) {
			c.String(http.StatusOK, "<h1>This is Gee Group</h1>")
		})
		v1.GET("/hello", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"group": c.Path,
			})
		})
	}

	r.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Param("name"), c.Path)
	})

	r.GET("/ass/*filepath", func(c *Context) {
		c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *Context) {
		c.JSON(http.StatusOK, H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/getHtml", func(c *Context) {
		c.HTML(http.StatusOK, "tes.tmpl", H{
			"msg":    "this is a template",
			"status": "ok",
		})
	})

	log.Fatalln(r.Run(":8888"))
}
