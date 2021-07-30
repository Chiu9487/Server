package handler

import (
	"chatdemo/commont/message"
	"chatdemo/server/manager"
	"chatdemo/server/redisUtils"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func LoginHandler(msg *message.ReqMessage ,conn net.Conn , redisPools *redisUtils.Redispool)(int,string,int){
	//defer redisConn.Close()
	loginmsg := msg.MsgData
	users := strings.Split(loginmsg,"-")
	userId,_ := strconv.Atoi(users[0])
	pwd := users[1]

	user,err :=  redisPools.GetUser(userId);
	//user,err :=  redis.String(redisConn.Do("hget" ,"users" ,userId))
	if err != nil{
		if err.Error() =="redigo: nil returned"{
			redisPools.SaveUser(userId,pwd)
			manager.OnlineHandler(userId ,&conn)
			return 2,"没有用户，此号创建用户，并登陆成功",userId
		} else {
			fmt.Println("查找user失败 ，err: " ,err)
			return -2,"查询玩家失败",0
		}

	}
	//user := "1-1"
	if string(user) == ""{ //找到的用户不存在，保存用户，就用此号登陆
		redisPools.SaveUser(userId,pwd)
		//user := fmt.Sprintf("%d-%s",userId,pwd)
		//redisConn.Do("hset","users",userId,user)

		manager.OnlineHandler(userId ,&conn)
		return 2,"没有用户，此号创建用户，并登陆成功",userId
	}

	fmt.Println("比较两个是否一致",user,loginmsg)
	if string(loginmsg) != string(user) { //密码错误
		return -1,"密码错误",0
	}
	//登陆成功
	manager.OnlineHandler(userId ,&conn)
	return 1,"登陆成功",userId
}



