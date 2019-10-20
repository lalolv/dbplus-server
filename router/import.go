package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func importData(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	// file.Filename
	filename := filepath.Base(file.Filename)
	distFilePath := fmt.Sprintf("%s/%s", "temp", filename)
	if err := c.SaveUploadedFile(file, distFilePath); err != nil {
		fmt.Println("err", err.Error())
		respData(c, -1, "", err.Error())
		return
	}

	// read all content
	bytes, err := ioutil.ReadFile(distFilePath)
	if err != nil {
		respData(c, -5, err.Error(), "")
		return
	}
	// 解析json
	var uData []interface{}
	err = json.Unmarshal(bytes, &uData)
	if err != nil {
		fmt.Println("err", err.Error())
	}

	// 目标集合
	collName := c.PostForm("coll_name")
	dbName := c.PostForm("db_name")
	if collName != "" && dbName != "" {
		coll := mgos.DB(dbName).C(collName)
		// update
		err = coll.Insert(uData...)
		if err != nil {
			respData(c, -3, err.Error(), "")
			return
		}
	}

	// Resp
	respData(c, 0, "OK", "")
}
