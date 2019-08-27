package router

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	conf bson.M
	mgos *mgo.Session
)

// User 用户信息
type User struct {
	Name   string `form:"name" json:"name"`
	Passwd string `form:"passwd" json:"passwd"`
}

// Server 数据服务器
type Server struct {
	Key  string `form:"key" json:"key"`
	Name string `form:"name" json:"name"`
	Desc string `form:"desc" json:"desc"`
	Addr string `form:"addr" json:"addr"`
	Port string `form:"port" json:"port"`
}

// Database 数据库，包含了集合
type Database struct {
	Server
	Dbs [][]string `form:"dbs" json:"dbs"`
}
