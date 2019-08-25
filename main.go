package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lalolv/dbplus-server/router"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

func main() {
	// 读取配置文件
	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		panic(err.Error())
	}
	// 解析 YAML 格式
	var conf bson.M
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		fmt.Println("YAML 解析错误")
		panic(err.Error())
	}

	// 调试模式设置
	if conf["debug"].(int) == 1 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := router.Startup(conf)

	serveConf := conf["server"].(map[interface{}]interface{})
	s := &http.Server{
		Addr:           fmt.Sprintf(":%v", serveConf["port"]),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// stop
	go graceful(s, 5*time.Second)

	// 运行服务
	s.ListenAndServe()
}

// 停止服务
func graceful(hs *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Println("Shutdown with timeout")

	if err := hs.Shutdown(ctx); err != nil {
		fmt.Println("Error: %v\n" + err.Error())
	} else {
		fmt.Println("Server stopped")
	}
}
