package router

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

// getServer 获取服务器信息
func getServer(c *gin.Context) {
	// Clone session
	ss := mgos.Clone()
	defer ss.Close()
	// All dbs
	dbs, _ := ss.DatabaseNames()
	dbList := []bson.M{}
	for _, dbName := range dbs {
		collNames, _ := ss.DB(dbName).CollectionNames()
		dbList = append(dbList, bson.M{"db": dbName, "colls": collNames})
	}

	// Resp
	respData(c, 0, "ok", dbList)
}
