package router

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lalolv/goutil"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/mgo.v2/bson"
)

// addGroup 创建新组
// @name @desc
// g-xxx: {@name, @desc}
// groups: [g-xxx]
func addGroup(c *gin.Context) {
	var params Group
	c.BindJSON(&params)

	// 打开数据文件
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		fmt.Println("open err: ", err.Error())
	}
	defer db.Close()

	// 读取组列表
	gList, _ := db.Get([]byte("groups"), nil)
	var nowGroups []string
	json.Unmarshal(gList, &nowGroups)

	// 生成key
	key := fmt.Sprintf("g-%s", bson.NewObjectId().Hex())

	// 批量写入
	batch := new(leveldb.Batch)
	// 添加组信息
	groupVal, _ := json.Marshal(gin.H{"name": params.Name, "desc": params.Desc})
	batch.Put([]byte(key), groupVal)
	// 添加组列表
	nowGroups = append(nowGroups, key)
	groupList, _ := json.Marshal(nowGroups)
	batch.Put([]byte("groups"), groupList)
	// 保存
	err = db.Write(batch, nil)
	if err != nil {
		fmt.Println("put err: ", err.Error())
	}

	// Resp
	respData(c, 0, "ok", key)
}

// updateGroup 更新组
// @name @desc
func updateGroup(c *gin.Context) {
	// 参数
	key := c.Query("key")
	var params Group
	c.BindJSON(&params)

	// 打开数据文件
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		fmt.Println("open err: ", err.Error())
	}
	defer db.Close()

	// 读取服务器列表
	group, _ := db.Get([]byte(key), nil)
	var nowGroup gin.H
	json.Unmarshal(group, &nowGroup)
	var servers []interface{}
	if nowGroup["servers"] != nil {
		servers = nowGroup["servers"].([]interface{})
	} else {
		servers = []interface{}{}
	}

	// 更新组信息
	groupVal, _ := json.Marshal(gin.H{"name": params.Name, "desc": params.Desc, "servers": servers})
	db.Put([]byte(key), groupVal, nil)

	// Resp
	respData(c, 0, "ok", key)
}

// 读取全部组的列表
// groups: [g1, g2, g3]
// g1: {name, desc, servers:[s1, s2, s3]}
// s1: {addr, port}
func groupList(c *gin.Context) {
	// 打开数据文件
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		fmt.Println("open err: ", err.Error())
	}
	defer db.Close()

	allList := []Group{}
	// 读取组列表
	groups, _ := db.Get([]byte("groups"), nil)
	var groupList []string
	json.Unmarshal(groups, &groupList)

	for _, g := range groupList {
		// 组信息
		strGroup, _ := db.Get([]byte(g), nil)
		var group gin.H
		json.Unmarshal(strGroup, &group)
		// 获取服务器列表
		gServers := []Server{}
		if gs, ok := group["servers"].([]interface{}); ok {
			for _, s := range gs {
				sKey, _ := goutil.ToString(s)
				if sKey != "" {
					strServer, _ := db.Get([]byte(sKey), nil)
					var server Server
					json.Unmarshal(strServer, &server)
					server.Key = sKey
					// 追加到列表
					gServers = append(gServers, server)
				}
			}
		}
		// 追加到列表
		allList = append(allList,
			Group{Key: g, Name: group["name"].(string), Desc: group["desc"].(string), Servers: gServers})
	}

	// Resp
	respData(c, 0, "ok", allList)
}
