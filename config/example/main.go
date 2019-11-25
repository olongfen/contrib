package main

import (
	"github.com/gin-gonic/gin"
	"github.com/srlemon/go-middle/config"
)

var a = &c{
	Name: "dsdasd",
	GG: struct {
		Name string
		Age  int
	}{Name: "张三", Age: 19},
}

func main() {

	var err error
	if err = config.LoadConfiguration("test.yml", a, a); err != nil {
		panic(err)
	}
	if err = a.Save(nil); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/", getConfig)
	r.Run("0.0.0.0:1563")

}

func getConfig(ctx *gin.Context) {
	if err := config.LoadConfiguration("test.yml", a, a); err != nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"msg": err,
		})
		return
	}
	ctx.AbortWithStatusJSON(200, a)
	return
}

type c struct {
	config.Config `yaml:"-"`
	Name          string `yaml:"name" json:"name"`
	GG            struct {
		Name string
		Age  int
	} `yaml:"用户"`
}
