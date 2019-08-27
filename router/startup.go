package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/mattes/go-asciibot"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Startup 启动路由
// @r 路由
func Startup(cf bson.M) *gin.Engine {
	conf = cf
	fmt.Println("Startup router ...")
	r := gin.New()
	// 跨域
	r.Use(cors())
	// 压缩 gzip.BestCompression gzip.BestSpeed gzip.NoCompression
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// Routers
	fmt.Println("Create routers ...")
	routers(r)
	// 连接数据库服务器
	connectMongoDB()
	// 随机显示机器人
	fmt.Println(asciibot.Random())
	// 监控信息
	// monitoring()

	return r
}

// 路由
func routers(r *gin.Engine) {
	r.GET("/test", test)
	// Login
	r.POST("/user/login", login)
	// dblist
	r.GET("/server/get", getServer)
	// 数据集
	r.GET("/data/list", dataList)
	r.PUT("/data/update", updateData)
	r.POST("/data/add", addData)
	r.POST("/data/remove", removeData)
}

// 身份密钥验证
func keyRequired(c *gin.Context) {

}

func connectMongoDB() {
	var err error
	mgos, err = mgo.Dial(fmt.Sprintf("%s:%s", "127.0.0.1", "27017"))
	if err != nil {
		fmt.Println("Can't connect MongoDB server!")
		panic(err.Error())
	}
	mgos.SetMode(mgo.Monotonic, true)
	fmt.Println("Connected to MongoDB!")
}

// 实时监控服务器性能
func monitoring() {
	ss := mgos.Clone()
	defer ss.Close()

	for {
		var result bson.M
		ss.Run(bson.M{"serverStatus": 1}, &result)

		network := result["network"].(bson.M)
		r1 := fmt.Sprintf("网络：传入%d字节 输出%d字节 请求总数%d", network["bytesIn"], network["bytesOut"], network["numRequests"])
		mem := result["mem"].(bson.M)
		r2 := fmt.Sprintf("内存占用：%d", mem["resident"])
		fmt.Println(r1, r2)

		time.Sleep(time.Second * 2)
	}
}

// 测试连接
func test(c *gin.Context) {
	// Resp
	respData(c, 0, "OK!", fmt.Sprintf("已成功连接到服务接口"))
}

// RespData 输出数据到客户端
// @code
// @data
func respData(c *gin.Context, code int, msg string, data interface{}) {
	// 关闭请求
	c.Request.Body.Close()
	// 头部
	c.Header("content-type", "application/json")
	// 跨域访问
	c.Header("Access-Control-Allow-Origin", "*")
	// 输出json格式数据
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"now":  time.Now().Unix(),
		"data": data,
	})
}
