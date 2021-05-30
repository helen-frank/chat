package model

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

//在服务器启动后，就初始化一个userDao实例
//把它做成全局的变量，在需要和redis操作时，就直接使用即可

var (
	MyUserDao *UserDao
)

//定义一个UserDao结构体
//完成对User 结构体的各种操作

type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式，创建一个UserDao 实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 1. 根据用户 id 返回一个User实例 + error

func (ud *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	//通过给定 id 去 redis 查询这个用户
	err = conn.Send("auth", "123456")
	if err != nil {
		fmt.Println("redis密码错误,err = ", err)
	}
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		//表示在 user 哈希中，没有找到对应id
		if err == redis.ErrNil {
			err = ErrorUserNotexists
		}
		return
	}

	//需要把res 反序列化程User 实例
	user = &User{}
	err = json.Unmarshal([]byte(res), user)

	if err != nil {
		fmt.Println("json.Unmarshal([]byte(res), user) fail , err = ", err)
		return
	}

	return
}

func (ur *UserDao) addUser(conn redis.Conn, userId int, userPwd, userName string) (err error) {
	//通过给定 id 去 redis 查询这个用户
	user := User{
		UserId:   userId,
		UserPwd:  userPwd,
		UserName: userName,
	}
	err = conn.Send("auth", "123456")
	if err != nil {
		fmt.Println("redis密码错误,err = ", err)
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json.Marshal(user) fail , err = ", err)
	}

	_, err = conn.Do("hset", "users", userId, string(data))

	if err != nil {
		fmt.Println("conn.Do(\"hset\", \"users\", userId, string(data)) fail , err = ", err)
		return
	}
	return
}

//完成登录的校验 Login
// 1. Login 完成对用户的验证
// 2. 如果用户的id和pwd都正确，则返回一个user实例
// 3. 如果用户的id和pwd有错误，则返回对应的错误信息

func (ud *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//先从 UserDao 的连接池中取出一个连接
	conn := ud.pool.Get()
	defer conn.Close()

	user, err = ud.getUserById(conn, userId)
	if err != nil {
		return
	}

	// 这里证明这个用户是获取到的
	if user.UserPwd != userPwd {
		err = ErrorUserPwd
		return
	}
	return

}

func (ud *UserDao) Register(userId int, userPwd, userName string) (user *User, err error) {
	//先从 UserDao 的连接池中取出一个连接
	conn := ud.pool.Get()
	defer conn.Close()

	user, err = ud.getUserById(conn, userId)
	if err == nil {
		err = ErrorUserExists
		return
	}
	err = ud.addUser(conn, userId, userPwd, userName)
	if err != nil {
		fmt.Println("ud.addUser(conn, userId, userPwd, userName) fail , err = ", err)
	}
	return
}
