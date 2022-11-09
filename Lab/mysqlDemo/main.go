package mysqlDemo

//**********初始化连接数据库
//**********crud
//**********预处理
//**********事务
//**********sqlx

//**********初始化连接数据库********************************************************************
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//)
////定义一个全局对象db
//var db *sql.DB //sql.DB是表示连接的数据库对象（结构体实例），它保存了连接数据库相关的所有信息。它内部维护着一个具有零到多个底层连接的连接池，它可以安全地被多个goroutine同时使用
////定义一个初始化数据库函数
//func InitDB() (db *sql.DB, err error) {
//	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True" //连接信息
//	db, err = sql.Open("mysql", dsn)                                                  //给全局变量赋值，不要用简短变量,一定要导入驱动，否则mysql识别不了
//	if err != nil {
//		return db, err
//	}
//	err = db.Ping() // 尝试与数据库建立连接（校验dsn是否正确）
//	if err != nil {
//		return db, err
//	}
//	return db, nil //没有错误的话说明连接成功，返回一个nil代表没有错误
//}
//func main() {
//	_, err := InitDB()
//	if err != nil {
//		fmt.Printf("init db failed,err:%v\n", err)
//		return
//	}
//}

//**********CRUD*************************************************************************
//import (
//	"database/sql"
//	"fmt"
//)
//type user struct {
//	id, age int
//	name    string
//}
//func main() {
//	//构造数据库连接对象
//	dsn := "user:password@tcp(127.0.0.1:3306)/gorm_test"
//	db, _ := sql.Open("mysql", dsn)
//	//单行查询   func (db *DB) QueryRow(query string, args ...interface{}) *Row
//	var u user                                                  //构造结构体对象，以供查询数据写入
//	sqlstr := "select id, name, age from user where id=?"       //单行查询信息表示，？代表queryrow(信息表示，参数)
//	err := db.QueryRow(sqlstr, 30).Scan(&u.id, &u.name, &u.age) // 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
//	if err != nil {
//		fmt.Printf("scan failed,err:%v\n", err)
//		return
//	}
//	fmt.Printf("id:%d,name:%s,age:%d\n", u.id, u.name, u.age) //id:30,name:curry,age:33
//	//多行查询  func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
//	sqlmultistr := "select id,name,age from user where id>?"
//	rows, err := db.Query(sqlmultistr, 10) //注意和单行查询返回值不一样，rows此时存的是查询结果集，同时也要调用scan方法
//	if err != nil {
//		fmt.Printf("query failed,err:\n", err)
//		return
//	}
//	defer rows.Close() // 非常重要：关闭rows释放持有的数据库链接
//	for rows.Next() {  // 循环读取结果集中的数据
//		var u user //构造结构体对象，以供查询数据写入
//		err := rows.Scan(&u.id, &u.name, &u.age)
//		if err != nil {
//			fmt.Printf("scan failed,err:%v\n", err)
//			return
//		}
//		fmt.Printf("id:%d,name:%s,age:%d\n", u.id, u.name, u.age)
//		// id:11,name:klay,age:29
//		// id:23,name:green,age:29
//		// id:24,name:kobe,age:42
//		// id:30,name:curry,age:33
//	}
//	// //插入数据  插入、更新和删除操作都使用Exec方法func (db *DB) Exec(query string, args ...interface{}) (Result, error)Result是对已执行的SQL命令的总结
//	sqlinsertsql := "insert into user(name,age)values(?,?)"
//	retis, err := db.Exec(sqlinsertsql, "durant", 31)
//	if err != nil {
//		fmt.Printf("insert failed,err:%v\n", err)
//		return
//	}
//	retisId, err := retis.LastInsertId() // 新插入数据的id
//	if err != nil {
//		fmt.Printf("get lastinsert id failed,err:%v\n", err)
//		return
//	}
//	fmt.Printf("insert success,the id is %d.\n", retisId)
//	// //更新数据
//	sqlupdatestr := "update user set id=? where name=?"
//	retud, err := db.Exec(sqlupdatestr, 7, "durant")
//	if err != nil {
//		fmt.Printf("insert failed,err:%v\n", err)
//		return
//	}
//	n, err := retud.RowsAffected() // 操作影响的行数
//	if err != nil {
//		fmt.Printf("get RowsAffected failed, err:%v\n", err)
//		return
//	}
//	fmt.Printf("update success, affected rows:%d\n", n)
//	//删除数据
//	sqldelstr := "delete from user where id = ?"
//	retdl, err := db.Exec(sqldelstr, 24)
//	if err != nil {
//		fmt.Printf("delete failed, err:%v\n", err)
//		return
//	}
//	n, err = retdl.RowsAffected() // 操作影响的行数
//	if err != nil {
//		fmt.Printf("get RowsAffected failed, err:%v\n", err)
//		return
//	}
//	fmt.Printf("delete success, affected rows:%d\n", n)
//}

//**********预处理*************************************************************************
//我们任何时候都不应该自己拼接SQL语句！否则会引起注入问题！！！！！
//import (
//	"database/sql"
//	"fmt"
//)
//type user struct {
//	id, age int
//	name    string
//}
////预处理查询
////查找
//func preQuery(db *sql.DB) {
//	sqlstr := "select id,name,age from user where id>?"
//	stmt, err := db.Prepare(sqlstr) // 先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
//	if err != nil {
//		fmt.Printf("prepare failed,err:%v\n", err)
//		return
//	}
//	defer stmt.Close() //一定要注意别忘记关闭数据库连接
//	rows, err := stmt.Query(0)
//	if err != nil {
//		fmt.Printf("query failed,err:&v\n", err)
//		return
//	}
//	defer rows.Close() // 非常重要：关闭rows释放持有的数据库链接
//
//	for rows.Next() { // 循环读取结果集中的数据
//		var u user
//		err := rows.Scan(&u.id, &u.name, &u.age)
//		if err != nil {
//			fmt.Printf("scan failed ,err:%v\n", err)
//			return
//		}
//		fmt.Printf("id:%d,name:%s,age:%d\n", u.id, u.name, u.age)
//	}
//
//}
////插入
//func preInsert(db *sql.DB) {
//	sqlstr := "insert into user values(?,?,?)"
//	stmt, err := db.Prepare(sqlstr)
//	if err != nil {
//		fmt.Printf("prepare failed ,err:%v\n", err)
//		return
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(33, "lilard", 33)
//	if err != nil {
//		fmt.Printf("insert failed,err:%v\n", err)
//		return
//	}
//	fmt.Println("Insert success")
//
//}
////删除
//func preDelete(db *sql.DB) {
//	sqlstr := "delete from user where id=?"
//	stmt, err := db.Prepare(sqlstr)
//	if err != nil {
//		fmt.Printf("prepare failed ,err:%v\n", err)
//		return
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(32)
//	if err != nil {
//		fmt.Printf("delete failed,err:%v\n", err)
//		return
//	}
//	fmt.Println("delete success")
//
//}
////更新
////删除
//func preUpdate(db *sql.DB) {
//	sqlstr := "update user set age=? where id=?"
//	stmt, err := db.Prepare(sqlstr)
//	if err != nil {
//		fmt.Printf("prepare failed ,err:%v\n", err)
//		return
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(66, 33)
//	if err != nil {
//		fmt.Printf("update failed,err:%v\n", err)
//		return
//	}
//	fmt.Println("update success")
//}
//func main() {
//	dsn := "user:password@tcp(127.0.0.1:3306)/gorm_test"
//	db, _ := sql.Open("mysql", dsn)
//	preQuery(db)
//	preInsert(db)
//	preDelete(db)
//	preUpdate(db)
//}

//**********事务**********************************************
//import (
//	"database/sql"
//	"fmt"
//
//	_ "github.com/go-sql-driver/mysql" //init()导入驱动，但不使用
//)
//
//// 开始事务func (db *DB) Begin() (*Tx, error)
//// 提交事务func (tx *Tx) Commit() error
//// 回滚事务func (tx *Tx) Rollback() error
//type user struct {
//	id, age int
//	name    string
//}
//
//func transaction(db *sql.DB) {
//	tx, err := db.Begin() //开启事务
//	if err != nil {
//		if tx != nil {
//			tx.Rollback() // 回滚
//		}
//		fmt.Printf("begin trans failed, err:%v\n", err)
//		return
//	}
//	sqlstrinsert := "insert into user values(?,?,?)"
//	_, err = tx.Exec(sqlstrinsert, 9, "lebron", 38)
//	if err != nil {
//		tx.Rollback() //回滚
//		fmt.Printf("insert failed,err:%v\n", err)
//		return
//	}
//	fmt.Println("Insert success")
//	sqlstrdelete := "delete from user where id=?"
//	_, err = tx.Exec(sqlstrdelete, 30)
//	if err != nil {
//		tx.Rollback() //回滚事务
//		fmt.Printf("delete failed,err:%v\n", err)
//		return
//	}
//	fmt.Println("delete success")
//	tx.Commit() //提交事务
//	fmt.Println("事务成功提交了")
//}
//func main() {
//	dsn := "user:password@tcp(127.0.0.1:3306)/gorm_test"
//	db, _ := sql.Open("mysql", dsn)
//	transaction(db)
//}

//**********sqlx********************************************
//import (
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/jmoiron/sqlx"
//)
//type user struct {
//	Id   int
//	Name string
//	Age  int
//}
//var db *sqlx.DB
////查询单行数据
//func queryRow(db *sqlx.DB) {
//	sqlstr := "select id,name,age from user where id=?" //查询语句信息str
//	var u user                                          //定义结构体对象
//	err := db.Get(&u, sqlstr, 6)                        //调用参数数据库对象，执行查询单行方法
//	// 结构体user的访问限制改为全局（即大写)。 这里我们可以理解了，当db.Get(&u,sql,id)执行时，
//	// sqlx包中会访问结构体u中的各字段，这时发现字段全部为小写，不可访问，即报错了。我们修改为大写即解决了问题。
//	if err != nil {
//		fmt.Printf("get failed ,err:%v\n", err)
//		return
//	}
//	fmt.Printf("users:%#v\n", u)
//	// fmt.Printf("id:%d,name:%s,age:%d\n", u.d, u.name, u.age) //输出查询信息
//
//}
//
////查询多行数据
//func queryMultiRow(db *sqlx.DB) {
//	sqlstr := "select id,name,age from user where id>?" //查询语句信息str
//	var users []user                                    //定义结构体对象切片
//	err := db.Select(&users, sqlstr, 11)                //调用参数数据库对象，执行查询单行方法
//	if err != nil {
//		fmt.Printf("get failed ,err:%v\n", err)
//		return
//	}
//	fmt.Printf("users:%3v\n", users) //输出查询信息
//}
//func main() {
//	dsn := "user:password@tcp(127.0.0.1:3306)/gorm_test"
//	db, _ := sqlx.Connect("mysql", dsn)\
//	db.SetMaxOpenConns(20)
//	db.SetMaxIdleConns(10)
//	queryRow(db)
//}

//****************sqlx.In*******************************************
//import (
//	"database/sql/driver"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/jmoiron/sqlx"
//)
//var db *sqlx.DB
//type User struct {
//	Name string `db:"name"`
//	Age  int    `db:"age"`
//}
//func (u User) Value() (driver.Value, error) {
//	return []interface{}{u.Name, u.Age}, nil
//}
//func initDB() (err error) {
//	dsn := "root:123456@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
//	// 也可以使用MustConnect连接不成功就panic
//	db, err = sqlx.Connect("mysql", dsn)
//	if err != nil {
//		fmt.Printf("connect DB failed, err:%v\n", err)
//		return
//	}
//	db.SetMaxOpenConns(20)
//	db.SetMaxIdleConns(10)
//	return
//}
//
//// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
//func BatchInsertUsers2(users []interface{}) error {
//	query, args, _ := sqlx.In(
//		"INSERT INTO user (name, age) VALUES (?), (?), (?)",
//		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
//	)
//	fmt.Println(query) // 查看生成的querystring INSERT INTO user (name, age) VALUES (?, ?), (?, ?), (?, ?)
//	fmt.Println(args)  // 查看生成的args  [七米 18 q1mi 28 小王子 38]
//	_, err := db.Exec(query, args...)
//	return err
//}
//func main() {
//	err := initDB()
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//	u1 := User{Name: "七米", Age: 18}
//	u2 := User{Name: "q1mi", Age: 28}
//	u3 := User{Name: "小王子", Age: 38}
//	// 方法2
//	users2 := []interface{}{u1, u2, u3}
//	err = BatchInsertUsers2(users2)
//	if err != nil {
//		fmt.Printf("BatchInsertUsers2 failed, err:%v\n", err)
//	}
//}
