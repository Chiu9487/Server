package handler

import (
	"chatdemo/commont/message"
	"chatdemo/server/manager"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net"
)

func Chat(msg *message.ReqMessage ,conn net.Conn,currUserId int) int {
	data := msg.MsgData

	var chatData message.ChatData
	err := json.Unmarshal([]byte(data) ,&chatData)
	chatData.FromId = currUserId
	if err != nil{
		fmt.Println("反序列化聊天内容失败，err" ,err)
		return -1
	}
	if chatData.ChatType == 1{ //私聊
		fmt.Println("进行私聊")
		manager.PrivateChat(&chatData)
	} else { //公聊
		fmt.Println("进行共聊")
		manager.AddMsgTpChannel(&chatData)
	}
	return 1
}

