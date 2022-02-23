package main

import (
	"goWebGee/gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, your path is %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":8888")

}
