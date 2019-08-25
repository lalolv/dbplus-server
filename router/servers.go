package router

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

// addServer 创建新的服务器
// @name @desc @addr @port
// @group key
func addServer(c *gin.Context) {
	// 参数
	groupKey := c.Query("group")
	var params Server
	c.BindJSON(&params)

	// 打开数据文件
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		fmt.Println("open err: ", err.Error())
	}
	defer db.Close()

	// 读取组信息
	g, _ := db.Get([]byte(groupKey), nil)
	var group gin.H
	json.Unmarshal(g, &group)
	// 服务器列表
	serveList := []interface{}{}
	if group["servers"] != nil {
		serveList = append(serveList, group["servers"].([]interface{})...)
	}

	// 新的key
	newKey := fmt.Sprintf("s-%s", bson.NewObjectId().Hex())
	// 批量写入
	batch := new(leveldb.Batch)
	// 保存服务器
	servVal, _ := json.Marshal(
		gin.H{
			"name": params.Name,
			"desc": params.Desc,
			"addr": params.Addr,
			"port": params.Port,
		})
	batch.Put([]byte(newKey), servVal)
	// 保存列表到组中
	serveList = append(serveList, newKey)
	groupVal, _ := json.Marshal(
		gin.H{
			"name":    group["name"],
			"desc":    group["desc"],
			"servers": serveList,
		})
	batch.Put([]byte(groupKey), groupVal)
	// 保存
	err = db.Write(batch, nil)
	if err != nil {
		fmt.Println("put err: ", err.Error())
	}

	// Resp
	respData(c, 0, "ok", newKey)
}

// updateServer 更新服务信息
func updateServer(c *gin.Context) {
	// 参数
	key := c.Query("key")
	var params Server
	c.BindJSON(&params)

	// 打开数据文件
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		fmt.Println("open err: ", err.Error())
	}
	defer db.Close()

	// 更新组信息
	servVal, _ := json.Marshal(
		gin.H{"name": params.Name, "desc": params.Desc, "addr": params.Addr, "port": params.Port})
	db.Put([]byte(key), servVal, nil)

	// Resp
	respData(c, 0, "ok", key)
}

// getServer 获取服务器信息
func getServer(c *gin.Context) {
	// 参数
	sKey := c.Query("key")

	// 打开数据文件
	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		fmt.Println("open err: ", err.Error())
	}
	defer db.Close()

	// get server data
	s, _ := db.Get([]byte(sKey), nil)
	var ss Server
	json.Unmarshal(s, &ss)
	ss.Key = sKey

	// 连接数据库服务器
	if mgos != nil {
		fmt.Println("close mongo session")
		mgos.Close()
	}
	// connStr := fmt.Sprintf("%s:%s", ss.Addr, ss.Port)
	// fmt.Println(connStr)
	mgos, err = mgo.Dial(fmt.Sprintf("%s:%s", ss.Addr, ss.Port))
	if err != nil {
		fmt.Println("MongoDB连接失败")
		respData(c, -1, "MongoDB连接失败", bson.M{})
		return
	}
	mgos.SetMode(mgo.Monotonic, true)
	fmt.Println("Connected to MongoDB!")

	// get database and coll
	sess := mgos.Clone()
	defer sess.Close()
	dbs, _ := sess.DatabaseNames()
	dbList := []bson.M{}
	for _, dbName := range dbs {
		collNames, _ := sess.DB(dbName).CollectionNames()
		dbList = append(dbList, bson.M{"db": dbName, "colls": collNames})
	}

	// Resp
	respData(c, 0, "ok", bson.M{"server": ss, "dbs": dbList})
}
