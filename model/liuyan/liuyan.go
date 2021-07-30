package liuyan

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var(
	Db *sql.DB //mysql链接
)

func init(){
	var err error
	Db, err = sql.Open("mysql", "root:qwer1234@tcp(127.0.0.1:3306)/chat?charset=utf8")
	if err != nil {
		log.Fatal(err)
		return
	}
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(5)
}


type LiuYan struct {
	Id int            `db:"id"` //主键id
	ReceiveId int     `db:"receiveId"`//接收者id
	SendId int        `db:"sendId"`//发送者id
	Content string    `db:"content"`//内容
}

// 根据id查询是否有留言
func QueryLiuYanByReceiveId( receiveId int)( []*LiuYan){
	fmt.Println("开始查询mysql")
	stmtOut, err := Db.Prepare("SELECT * FROM `liuyan` WHERE receiveId = ?")
	if err != nil {
		fmt.Println("创建命令失败 ，err: ",err)
		return nil
	}
	defer stmtOut.Close()

	//rows := stmtOut.QueryRow()
	rows, err := stmtOut.Query(receiveId)
	if err != nil {
		panic(err.Error())
		return nil

	}
	var liuyans []*LiuYan = make([]*LiuYan ,0,10)

	for rows.Next() {
		liuyan := new(LiuYan)
		err = rows.Scan(&liuyan.Id, &liuyan.ReceiveId, &liuyan.SendId, &liuyan.Content)
		if err != nil {
			panic(err.Error())
			return nil
		}
		liuyans = append(liuyans, liuyan)
	}
	return liuyans
}

//添加留言
func AddLiuYan(liuyan *LiuYan ){
	stmtOut, err := Db.Prepare("INSERT INTO `liuyan` (`receiveId`, `sendId`,`content`) values (?, ?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	result, err := stmtOut.Exec(liuyan.ReceiveId, liuyan.SendId,liuyan.Content)
	if err != nil {
		panic(err.Error())
	}
	_ , err = result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
}

//删除留言
func DelLiuYan(recevieId int)error{
	stmtOut, err := Db.Prepare("DELETE FROM `liuyan` WHERE receiveId = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	result, err := stmtOut.Exec(recevieId)
	if err != nil {
		panic(err.Error())
	}

	rowNum, err := result.RowsAffected();
	if err != nil  {
		//panic("delete error")
		return  err
	}
	fmt.Printf("删除了 %d 条留言，\n",rowNum)
	return nil
}

//通过id删除
func DeletById( id int)error{
	stmtOut, err := Db.Prepare("DELETE FROM `liuyan` WHERE Id= ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	result, err := stmtOut.Exec(id)
	if err != nil {
		panic(err.Error())
	}

	rowNum, err := result.RowsAffected();
	if err != nil  {
		//panic("delete error")
		return  err
	}
	fmt.Printf("删除了 %d 条留言，\n",rowNum)
	return nil
}