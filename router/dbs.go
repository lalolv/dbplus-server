package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lalolv/goutil"
	"gopkg.in/mgo.v2/bson"
)

// getDatas 获取数据集
// 使用查询语句获取数据集 db.role.find()
// 返回头部字段信息
// 参数：@collName
func dataList(c *gin.Context) {
	// 参数
	dbName := c.Query("db_name")
	collName := c.Query("coll_name")

	// 目标集合
	ss := mgos.Clone()
	defer ss.Close()
	// databse
	db := ss.DB(dbName)

	// 数据列表
	header := []bson.M{}
	var list []bson.M

	// 获取数据集
	err := db.C(collName).Find(nil).Skip(0).Limit(20).All(&list)
	if err != nil {
		respData(c, -1, err.Error(), bson.M{"header": header, "list": list})
		return
	}

	// 获取头部信息
	for _, v := range list {
		for k := range v {
			// 判断重复
			var isContain bool
			for _, t := range header {
				if t["value"] == k {
					isContain = true
					break
				}
			}
			// 追加
			if !isContain {
				header = append(header, bson.M{"text": k, "value": k})
			}
		}
	}

	// Resp
	respData(c, 0, "ok", bson.M{"header": header, "list": list})
}

// updateData 修改数据
// 修改指定字段的值, 没有则增加新字段
// 获取行数据的主键
// @collName @columnName @val
func updateData(c *gin.Context) {
	// 参数
	dbName := c.Query("db_name")
	collName := c.Query("coll_name")
	// 参数体
	var params bson.M
	err := c.BindJSON(&params)
	if err != nil {
		respData(c, -1, err.Error(), "")
		return
	}

	columnName, _ := goutil.ToString(params["column_name"])
	updateVal := params["update_val"]
	id := c.Query("id")

	// 目标集合
	coll := mgos.DB(dbName).C(collName)
	// update
	err = coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{columnName: updateVal}})
	if err != nil {
		respData(c, -2, err.Error(), "")
		return
	}

	// Resp
	respData(c, 0, "ok", "")
}

// addData 新增数据
// @data json格式数据
func addData(c *gin.Context) {
	// 参数
	dbName := c.Query("db_name")
	collName := c.Query("coll_name")
	// 参数体
	var params bson.M
	err := c.BindJSON(&params)
	if err != nil {
		respData(c, -1, err.Error(), "")
		return
	}

	if len(params) > 0 {
		// 目标集合
		coll := mgos.DB(dbName).C(collName)
		// add
		err = coll.Insert(params)
		if err != nil {
			respData(c, -2, err.Error(), "")
			return
		}
	}

	// Resp
	respData(c, 0, "ok", "")
}

// removeData 删除数据
func removeData(c *gin.Context) {
	// 参数
	dbName := c.Query("db_name")
	collName := c.Query("coll_name")
	id := c.Query("id")
	// 参数体
	var params bson.M
	err := c.BindJSON(&params)
	if err != nil {
		respData(c, -1, err.Error(), "")
		return
	}

	// 目标集合
	coll := mgos.DB(dbName).C(collName)
	// remove
	err = coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		respData(c, -2, err.Error(), "")
		return
	}

	// Resp
	respData(c, 0, "ok", "")
}
