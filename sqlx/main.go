package main

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sqlx.DB

type user struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func initDB() (err error) {
	db, err = sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

func simpleQuery() {
	sqlStr := "select id, name, age from user where id = ?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.ID, u.Name, u.Age)
}

func simpleQueryMulti() {
	sqlStr := "select id, name, age from user where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

func simpleInsert() {
	sqlStr := "insert into user(name, age) values(?,?)"
	ret, err := db.Exec(sqlStr, "沙河小王子", 19)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theId, err := ret.LastInsertId() // 新插入的数据ID
	if err != nil {
		fmt.Printf("get lastInsertId failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theId)
}

func simpleUpdate() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 2)
	if err != nil {
		fmt.Printf("udpate failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed. err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

func namedExec() (err error) {
	_, err = db.NamedExec(`insert into user (name, age) values (:name, :age)`,
		map[string]interface{}{
			"name": "张三",
			"age":  33,
		})
	return
}

func namedQuery() {
	sqlStr := "select * from user where name=:name"
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "Paul"})
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		rows.StructScan(&u)
		fmt.Printf("user:%#v\n", u)
	}

	// 另一种使用
	u := user{
		Name: "张三",
	}
	rows, err = db.NamedQuery(sqlStr, &u)
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		rows.StructScan(&u)
		fmt.Printf("user:%#v\n", u)
	}
}

func transactionDemo() (err error) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
			fmt.Println("commit")
		}
	}()

	sqlStr1 := "update user set age = 20 where id = ?"
	rs, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}

	sqlStr2 := "update user set age=59 where id = ?"
	rs, err = tx.Exec(sqlStr2, 3)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr2 failed")
	}
	return err
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("init DB failed, err:%v\n", err)
	}
	fmt.Println("init DB success...")
	//simpleQuery()
	//simpleQueryMulti()
	//simpleInsert()
	//simpleUpdate()
	//namedExec()
	//namedQuery()
	transactionDemo()
}
