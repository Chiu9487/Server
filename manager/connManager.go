package manager

import (
	"chatdemo/commont/message"
	liuyan2 "chatdemo/server/model/liuyan"
	"encoding/json"
	"fmt"
	"log"
	"net"
	_ "github.com/go-sql-driver/mysql"
)

var(
	connMap map[int]*net.Conn //链接map
	broadMess chan string // 广播通道
)

func init(){
	connMap = make(map[int]*net.Conn)
	broadMess = make(chan string)
}

//上线
func OnlineHandler(userId int ,conn *net.Conn){
	connMap[userId] = conn
}

//下线
func OfflineHandler(userId int){
	delete(connMap ,userId)
}

//私聊
func PrivateChat(data *message.ChatData){
	//正确判断方法
	conn, ok := connMap[data.ToId]
	if !ok{ //不存在就留言
		//todo redis留言系统
		liuyan := liuyan2.LiuYan{
			ReceiveId: data.ToId,
			SendId: data.FromId,
			Content: data.Content,
		}
		liuyan2.AddLiuYan(&liuyan)
	} else { //发送
		b,err := json.Marshal(data)
		if err != nil{
			fmt.Println("序列话聊天内容失败，err：" ,err)
			return
		}
		backMsg := message.CreateResMessage(2,string(b),2)

		backStr,err := json.Marshal(backMsg) //返回给客户端的消息
		if err != nil{
			fmt.Println("服务器序列话返回消息失败 ：err : ",err)
			return
		}
		_,err = (*conn).Write(backStr)
		if err != nil{
			fmt.Println("发送聊天失败,err:",err)
		}
	}
}

//将消息发到广播通道中，另起协程进行广播
func AddMsgTpChannel(chatData *message.ChatData){
	str, err := json.Marshal(chatData)
	if err != nil{
		fmt.Println("序列话消息失败 ，err: ",err)
	}
	res := message.CreateResMessage(2,string(str),2)
	s,err := json.Marshal(res)
	if err != nil{
		fmt.Println("序列话消息失败 ，err: ",err)
	}
	broadMess <- string(s)

	//fmt.Println()
	for id,_ := range connMap{
		fmt.Println("所有的链接id是：" ,id)
	}
}

//广播
func BroadMessage(){
	for {
		msg := <-broadMess
		for userId,conn := range connMap{

			fmt.Println("广播的msg是：" ,msg)

			_, err := (*conn).Write([]byte(msg))
			if err != nil {
				log.Printf("broad message to %s err: %v\n", userId, err)
				delete(connMap, userId)
			}
		}
	}

}


