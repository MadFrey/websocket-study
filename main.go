package main

import (
	"websocket/router"
	"websocket/util"
)

func main()  {
	//dns:="root:sjk123456@tcp(localhost:3306)/?parseTime=true"
	//err := dao.InitMysql(dns)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	go util.NewHub.Start()
	r:= router.InitRouter()
	r.Run(":9090")
}
