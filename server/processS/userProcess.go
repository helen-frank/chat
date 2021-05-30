package processS

import (
	"chat/common/message"
	"chat/server/model"
	"chat/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//表示该Conn 是哪个用户
	UserId int
}

//通知所有在线用户
//userId 这个用户要通知其他在线用户上线
func (up *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历 onlineUsers , 然后一个个发送 NotifyUserStatusMes
	for id, ou := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		//开始通知
		ou.NotifyMeOnline(userId)
	}
}

func (up *UserProcess) NotifyMeOnline(userId int) {
	//开始组装 NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId

	notifyUserStatusMes.UserStatus = message.UserOnline

	//将 notifyUserStatusMes 序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("userProcess.go > NotifyMeOnline > json.Marshal(notifyUserStatusMes) fail , err =", err)
		return
	}
	//将序列化后的notfiyUserStatusMes赋值给mes , Data
	mes.Data = string(data)

	//对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("userProcess.go > NotifyMeOnline > json.Marshal(mes) fail , err = ", err)
		return
	}

	//发送 ，创建 Transfer实例
	tf := &utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("userProcess.go > NotifyMeOnline > tf.WritePkg(data) fail , err = ", err)
		return
	}

}

// serverProcessLogin 函数，专门处理登录请求
func (up *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1. 先从 mes 中取出 mes.Data, 并直接反序列化成  LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), loginMes) fail , err = ", err)
		return
	}
	// 1.先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2.再声明一个 LoginResMes
	var LoginResMes message.LoginResMes

	//到redis数据库去完成验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ErrorUserNotexists {
			LoginResMes.Code = 500
			LoginResMes.Error = err.Error()
		} else if err == model.ErrorUserPwd {
			LoginResMes.Code = 403
			LoginResMes.Error = err.Error()
		} else {
			LoginResMes.Code = 505
			LoginResMes.Error = "服务器内部错误"
		}

	} else {
		LoginResMes.Code = 200
		//登录成功的用户的userId赋给up
		up.UserId = loginMes.UserId
		//把该登录成功的用户添加到userMgr中
		userMgr.AddOnlineUser(up)
		//通知其他在线用户，我已上线
		up.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户的id,放入到loginResMes.UserId
		//遍历 userMgr.onlineUsers

		for id := range userMgr.onlineUsers {
			LoginResMes.UsersId = append(LoginResMes.UsersId, id)
		}

		fmt.Println("用户", user.UserId, "登录成功")
	}

	data, err := json.Marshal(LoginResMes)
	if err != nil {
		fmt.Println("json.Marshal(LoginReMes) fail , err = ", err)
		return
	}

	// 4. 将 data 赋值给resMes
	resMes.Data = string(data)

	// 5. 将 resMes 序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail , err = ", err)
	}

	// 6. 发送
	// 使用分层模式（mvc)
	// 先创建一个 Trans 实例 ， 然后读取
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return
}

func (up *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 1. 先从 mes 中取出 mes.Data, 并直接反序列化成  LoginMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &registerMes) fail , err = ", err)
		return
	}
	// 1.先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	// 2.再声明一个 LoginResMes
	var registerResMes message.RegisterResMes

	//到redis数据库去完成注册
	_, err = model.MyUserDao.Register(registerMes.User.UserId, registerMes.User.UserPwd, registerMes.User.UserName)

	if err != nil {
		if err == model.ErrorUserExists {
			registerResMes.Code = 400
			registerResMes.Error = model.ErrorUserExists.Error()
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "服务器内部错误"
		}

	} else {
		registerResMes.Code = 200
		fmt.Println("用户", registerMes.User.UserId, "注册成功")
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal(registerResMes) fail , err = ", err)
		return
	}

	// 4. 将 data 赋值给resMes
	resMes.Data = string(data)

	// 5. 将 resMes 序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail , err = ", err)
	}

	// 6. 发送
	// 使用分层模式（mvc)
	// 先创建一个 Trans 实例 ， 然后读取
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return
}
