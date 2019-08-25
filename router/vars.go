package router

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	groups []Group
	conf   bson.M
	mgos   *mgo.Session
)

// dataPath 数据保存路径
const dataPath string = "data"

// User 用户信息
type User struct {
	Name   string `form:"name" json:"name"`
	Passwd string `form:"passwd" json:"passwd"`
}

// Group 服务器组
type Group struct {
	Key     string   `form:"key" json:"key"`
	Name    string   `form:"name" json:"name"`
	Desc    string   `form:"desc" json:"desc"`
	Servers []Server `form:"servers" json:"servers"`
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
