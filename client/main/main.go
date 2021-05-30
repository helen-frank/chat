package main

import (
	"chat/client/processC"
	"fmt"
)

//定义 用户 id ,密码
var userId int
var userPwd string
var userName string

func main() {
	//接收用户的选择
	var key int
	//判断是否还继续显示菜单
	//	var loop = true

	for {
		fmt.Println("-----------------欢迎登录多人聊天系统-----------------")
		fmt.Println("\t\t 1 登录聊天室")
		fmt.Println("\t\t 2 注册用户")
		fmt.Println("\t\t 3 退出系统")
		fmt.Println("\t\t 请选择（1-3）")
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanln(&userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanln(&userPwd)
			// 完成登录
			// 1. 创建UserProcess实例
			up := &processC.UserProcess{}
			err := up.Login(userId, userPwd)
			if err != nil {
				fmt.Println("up.Login(userId, userPwd) fail, err = ", err)
			}

			//loop = false
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户的id")
			fmt.Scanln(&userId)
			fmt.Println("请输入用户的名字(nickname)")
			fmt.Scanln(&userName)
			fmt.Println("请输入用户的密码")
			fmt.Scanln(&userPwd)
			up := &processC.UserProcess{}
			err := up.Register(userId, userPwd, userName)
			if err != nil {
				fmt.Println("up.Register(userId, userPwd, userName) fail, err = ", err)
			}

		case 3:
			fmt.Println("退出系统")

			//loop = false
		default:
			fmt.Println("你的输入有误，请重新输入")
		}
	}
}
