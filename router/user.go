package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lalolv/goutil"
)

func login(c *gin.Context) {
	var params User
	err := c.BindJSON(&params)
	if err != nil {
		respData(c, -2, err.Error(), "")
		return
	}

	// 获取默认用户名和密码
	cfUser := conf["user"].(map[interface{}]interface{})
	cfName, _ := goutil.ToString(cfUser["name"])
	cfPasswd, _ := goutil.ToString(cfUser["passwd"])

	// 身份验证
	var reCode int
	var reMsg string
	if params.Name == cfName && params.Passwd == goutil.MD5(cfPasswd) {
		reCode = 0
		reMsg = "OK"
	} else {
		reCode = -1
		reMsg = "用户名或密码不正确"
	}

	// Resp
	respData(c, reCode, reMsg, "")
}
