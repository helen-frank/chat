package main

import (
	"chat/server/model"
	"fmt"
	"net"
	"time"
)

// func readPkg(conn net.Conn) (mes message.Message, err error) {
// 	buf := make([]byte, 8096)
// 	fmt.Println("读取客户端发送的数据ing")
// 	_, err = conn.Read(buf[:4])
// 	if err != nil {
// 		fmt.Println("conn.Read(buf[:4]) fail , err = ", err)
// 		return
// 	}
// 	//根据buf[0:4] 转成一个 uint32 类型
// 	pkgLen := binary.BigEndian.Uint32(buf[0:4])

// 	//根据 pkgLen 读取消息内容

// 	n, err := conn.Read(buf[:pkgLen])

// 	if n != int(pkgLen) || err != nil {
// 		fmt.Println("conn.Read(buf[:pkgLen]) fail , err = ", err)
// 		return
// 	}

// 	// 把 pkgLen 反序列成 message.Message

// 	err = json.Unmarshal(buf[:pkgLen], &mes)

// 	if err != nil {
// 		fmt.Println("json.Unmarshal(buf[:pkgLen]) fail , err = ", err)
// 		return
// 	}
// 	return
// }

// func writePkg(conn net.Conn, data []byte) (err error) {
// 	// 先发送一个长度给对方
// 	pkgLen := uint32(len(data))
// 	var buf [4]byte

// 	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
// 	// 发送长度
// 	_, err = conn.Write(buf[:])
// 	if err != nil {
// 		fmt.Println("conn.Write(bytes) , err = ", err)
// 		return
// 	}

// 	n, err := conn.Write(data)

// 	if n != int(pkgLen) || err != nil {
// 		fmt.Println("conn.Write(data) fail , err = ", err)
// 		return
// 	}
// 	return

// }

// // serverProcessLogin 函数，专门处理登录请求
// func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
// 	// 1. 先从 mes 中取出 mes.Data, 并直接反序列化成  LoginMes
// 	var loginMes message.LoginMes
// 	err = json.Unmarshal([]byte(mes.Data), &loginMes)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal([]byte(mes.Data), loginMes) fail , err = ", err)
// 		return
// 	}

// 	// 1.先声明一个 resMes
// 	var resMes message.Message
// 	resMes.Type = message.LoginResMesType

// 	// 2.再声明一个 LoginResMes
// 	var LoginReMes message.LoginResMes

// 	// 如果 用户id = 100 , 密码=123456 ， 认为合法， 否则不合法

// 	if loginMes.UserId == 123 && loginMes.UserPwd == "123456" {
// 		//合法
// 		LoginReMes.Code = 200
// 	} else {
// 		//不合法
// 		LoginReMes.Code = 500
// 		LoginReMes.Error = "用户不存在"
// 	}

// 	// 3. 将loginResMes序列化

// 	data, err := json.Marshal(LoginReMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal(LoginReMes) fail , err = ", err)
// 		return
// 	}

// 	// 4. 将 data 赋值给resMes
// 	resMes.Data = string(data)

// 	// 5. 将 resMes 序列化
// 	data, err = json.Marshal(resMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal(resMes) fail , err = ", err)
// 	}

// 	// 6. 发送
// 	err = writePkg(conn, data)
// 	return
// }

// // serverProcessMes 函数
// // 功能：根据客户端发送消息种类的不同，决定调用哪个函数来处理
// func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
// 	switch mes.Type {
// 	case message.LoginMesType:
// 		err = serverProcessLogin(conn, mes)
// 	case message.RegisterMesType:
// 		//处理注册
// 	default:
// 		fmt.Println("消息类型不存在，无法处理")
// 	}
// 	return
// }

//处理客户端通讯
func process(conn net.Conn) {
	defer conn.Close()
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误 ，err = ", err)
		return
	}
}

// 编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool本身就是一个全局的变量
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	//当服务器启动时，就初始化redis的连接池
	initPool("127.0.0.1:20001", 16, 0, 300*time.Second)
	initUserDao()

	fmt.Println("服务器在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen, err = ", err)
		return
	}
	defer listen.Close()

	//一旦监听，就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端来连接服务器ing")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept, err = ", err)
		}

		//一旦连接成功，则启动一个协程和客户端保持通讯
		go process(conn)

	}
}
