# go-contrib  [![GoDoc](https://godoc.org/github.com/olongfen/contrib?status.svg)](https://godoc.org/github.com/olongfen/contrib)
 a some middleware of glang frame
 
# Install
` go get github.com/olongfen/contrib`
 
# USAGE

- config 一个可以热加载yaml配置文件的配置管理
 
 ```golang 
   
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

  ```
  
- log 日志管理，设置日志等级
```golang
 package main

import "github.com/olongfen/contrib/log"

func main() {
	l := log.NewLogFile("demo.log")
	l1, _ := log.NewLog(l)
	l1.Println("aaaaaaaaaa")
}
  ```
 
