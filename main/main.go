package main

import (
	"chatdemo/commont/message"
	"chatdemo/server/handler"
	"chatdemo/server/manager"
	"chatdemo/server/model/liuyan"
	"chatdemo/server/redisUtils"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net"
)

var (
	redisPools *redisUtils.Redispool  //redis链接池

)

func init(){
	redisPools = redisUtils.InitRedisPool() //初始化链接池

}


func main()  {
	listen,err := net.Listen("tcp","127.0.0.1:8181")
	if err != nil{
		fmt.Println("服务器监听8181端口失败，检查~~  原因：" ,err)
		return;
	}

	go  manager.BroadMessage()
	// 循环监听8181 判断是否有连接过来
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("服务器接受客户端的连接失败，原因：" ,err)
			continue;
		}

		go clientHandler(conn ,redisPools)
	}
}





//处理客户端发来的请求
func clientHandler(conn net.Conn ,redisPools *redisUtils.Redispool){
	fmt.Println("客户端连接成功")
	var currUserId int
	//接收此连接的消息，如果没有就阻塞
	for{
		buf := make([]byte , 4096)
		n , err  := conn.Read(buf)

		if n == 0{ //读取字段尾0、客户端断开链接
			fmt.Println("客户端断开链接")
			if currUserId != 0{
				manager.OfflineHandler(currUserId)
			}
			return;
		}

		if (err != nil && err != io.EOF){ //读取失败
			fmt.Println("服务器读取客户端的消息失败 ！ err = " ,err)
			return
		}

		fmt.Println("接受到的消息是：" ,string(buf[:n]))

		var reqMsg message.ReqMessage

		//提取消息
		err = json.Unmarshal(buf[:n],&reqMsg)
		if err != nil{
			fmt.Println("接收数据反序列化成结构体失败 err:" ,err)
			continue;
		}

		//var resStr string
		var backMsg *message.ResMessage
		switch(reqMsg.Type){
		case 1://登陆
			loginCode,str,userId := handler.LoginHandler(&reqMsg ,conn ,redisPools)
			backMsg = message.CreateResMessage(1,str,loginCode)
			currUserId = userId
		case 2://日常交互
		fmt.Println("进行日常交互")
			code := handler.Chat(&reqMsg,conn,currUserId )
			backMsg = message.CreateResMessage(2,"",code)
		case 3://退出
			manager.OfflineHandler(currUserId)
			backMsg = message.CreateResMessage(3,"",1)
		default:
			fmt.Println("错误的选项")
		}
			backStr,err := json.Marshal(backMsg) //返回给客户端的消息
			if err != nil{
				fmt.Println("服务器序列话返回消息失败 ：err : ",err)
				continue
			}

			if reqMsg.Type != 2{
				_,err = conn.Write(backStr)
				if err != nil{
					fmt.Println("服务器回客户端消息失败 ：err : ",err)
					continue
				}
			}

			if(reqMsg.Type == 1 && backMsg.Code >= 1){
				checkLiuyan(currUserId )
			}
			fmt.Println("操作结束")
	}
}


func checkLiuyan(currId int){
	//fmt.Println(db.Ping())
	liuyans := liuyan.QueryLiuYanByReceiveId(currId)
	if len(liuyans) == 0{
		return
	}
	for _,l := range liuyans{
		chatData := message.ChatData{
			ChatType : 1, //1:私聊   2:群聊
			ToId : l.ReceiveId,
			Content : l.Content,
			FromId :l.SendId,
		}
		manager.PrivateChat(&chatData)
		liuyan.DelLiuYan(l.Id)
	}

}

