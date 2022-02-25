package main

import (
	"goWebGee/gee"
	"log"
	"net/http"
	"time"
)

func Logger() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func GroupMid() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Println("group middleware")
		c.Next()
	}
}

func main() {
	r := gee.New()

	//加载静态资源,访问localhost:8888/assets/file1.txt相当于访问./resource/file1.txt
	r.Static("/assets", "./resource")
	r.LoadHTMLGlob("templates/*")

	//定义全局中间件
	r.Use(Logger())

	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "<h1>This is Gee</h1>")
	})

	//定义组内中间件
	v1 := r.Group("/v1")
	v1.Use(GroupMid())
	{
		v1.GET("/", func(c *gee.Context) {
			c.String(http.StatusOK, "<h1>This is Gee Group</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"group": c.Path,
			})
		})
	}

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Param("name"), c.Path)
	})

	r.GET("/ass/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/getHtml", func(c *gee.Context) {
		c.HTML(http.StatusOK, "tes.tmpl", gee.H{
			"msg": "this is a template",
		})
	})

	r.Run(":8888")

}
