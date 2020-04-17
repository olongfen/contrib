package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olongfen/contrib/config"
)

var a = &c{
	Name: "dsdasd",
	GG: struct {
		Name string
		Age  int
	}{Name: "张三", Age: 19},
	Data: map[string]interface{}{
		"name": "sdsd",
	},
}

func main() {

	var err error
	if err = config.LoadConfigAndSave("test.yml", a, a); err != nil {
		panic(err)
	}
	go a.MonitorChange()
	r := gin.Default()
	r.GET("/", getConfig)
	r.Run("0.0.0.0:1563")

}

func getConfig(ctx *gin.Context) {
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
	Data map[string]interface{}
}
