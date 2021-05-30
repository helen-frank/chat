package processC

import (
	"chat/common/message"
	"chat/server/utils"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type UserProcess struct {
}

func (up *UserProcess) Register(userId int, userPwd, userName string) (err error) {
	// 1. 链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial , err = ", err)
		return
	}
	defer conn.Close()
	// 2. 准备通过conn 发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType

	// 3. 创建一个LoginMes 结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. 将registerMes 序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal(loginMes) , err = ", err)
		return
	}
	// 5. 把data 赋给mes.Data 字段
	mes.Data = string(data)

	// 6. 将mes 进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) , err = ", err)
		return
	}

	// 7. data就是要发送的消息
	// 7.1 先把data 的长度发送给服务器
	// 先获取到data 的长度 转成一个表示长度的byte 切片

	pkgLen := uint32(len(data))
	var buf [4]byte

	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	_, err = conn.Write(buf[:])
	if err != nil {
		fmt.Println("conn.Write(bytes) , err = ", err)
		return
	}
	//fmt.Printf("客户端已经发送信息长度 = %d, 内容 = %s", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail , err = ", err)
		return
	}
	fmt.Println(string(data))

	//还需要处理服务器返回消息

	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) fail , err = ", err)
		return
	}

	// 将 mes 的Data 部分反序列化成 LoginResMes
	var registerResMes message.RegisterResMes

	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &loginResMes) fail , err = ", err)
	}

	if registerResMes.Code == 200 {
		fmt.Println("注册成功,请重新登录")
	} else {
		fmt.Println(registerResMes.Error)
	}
	return
}

func (up *UserProcess) Login(userId int, userPwd string) (err error) {
	//制定协议
	// 1. 链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial , err = ", err)
		return
	}
	defer conn.Close()
	// 2. 准备通过conn 发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType

	// 3. 创建一个LoginMes 结构体
	loginMes := message.LoginMes{
		UserId:  userId,
		UserPwd: userPwd,
	}

	// 4. 将loginMes 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal(loginMes) , err = ", err)
		return
	}
	// 5. 把data 赋给mes.Data 字段
	mes.Data = string(data)

	// 6. 将mes 进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) , err = ", err)
		return
	}

	// 7. data就是要发送的消息
	// 7.1 先把data 的长度发送给服务器
	// 先获取到data 的长度 转成一个表示长度的byte 切片

	pkgLen := uint32(len(data))
	var buf [4]byte

	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	_, err = conn.Write(buf[:])
	if err != nil {
		fmt.Println("conn.Write(bytes) , err = ", err)
		return
	}
	//fmt.Printf("客户端已经发送信息长度 = %d, 内容 = %s", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail , err = ", err)
		return
	}
	fmt.Println(string(data))

	//还需要处理服务器返回消息

	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) fail , err = ", err)
		return
	}

	// 将 mes 的Data 部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes

	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &loginResMes) fail , err = ", err)
	}

	if loginResMes.Code == 200 {
		//fmt.Println("登录成功")
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		//显示当前在线用户在线列表，遍历loginResMes.UsersId
		fmt.Println("当前在线用户列表如下")
		for _, v := range loginResMes.UsersId {
			if v == userId { //不显示自己的 id
				continue
			}
			fmt.Println("用户id:\t", v)
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		//客户端启动一个协程
		//该协程保持和服务器端的通讯，如果服务器有数据推送给客户端
		//则接收并显示在客户端的界面

		// 1.显示登录成功后的菜单
		go serverProcessMes(conn)
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
		err = errors.New(fmt.Sprint(loginResMes.Error))
	}
	return nil
}
