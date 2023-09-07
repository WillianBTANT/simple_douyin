package main

import (
	"simple-douyin/config"
	"strconv"
)

func main() {
	// 下面这一行是聊天服务、需要写再启动、并且把services文件夹拷贝过来
	//go service.RunMessageServer()

	//if err := db.InitMySQLDB(); err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//if err := model.CreateTables(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	r := createAndInitEngine()

	configLocal := config.GetConfig()
	host := configLocal.ServerHost
	port := configLocal.ServerPort
	if err := r.Run(host + ":" + strconv.Itoa(port)); err != nil {
		return
	}
}
