package redisUtils
import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Redispool struct{
	Pools *redis.Pool
}

func InitRedisPool() (*Redispool) {
	redispool := &Redispool{
			Pools : &redis.Pool{
			MaxIdle:     8,
			MaxActive:   0,
			IdleTimeout: 300,
			Dial:func () (redis.Conn, error){
				return redis.Dial("tcp", "localhost:6379")
				},
			},
	}
	fmt.Println("初始化redis链接池成功～～～")
	return redispool
}

//获取用户
func (this *Redispool)GetUser(userId int)(string,error){
	conn := this.Pools.Get()
	fmt.Println(conn)
	defer conn.Close()
	return redis.String(conn.Do("hget" ,"users" ,userId))
}

//保存用户
func (this *Redispool) SaveUser(userId int,pwd string){
	conn := this.Pools.Get()
	defer conn.Close()
	user := fmt.Sprintf("%d-%s",userId,pwd)
	conn.Do("hset","users",userId,user)
}
