package processC

import (
	"chat/client/utils"
	"chat/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面
func ShowMenu() {
	fmt.Println("-----------------恭喜xxx登录成功-----------------")
	fmt.Println("\t\t 1. 显示在线用户列表")
	fmt.Println("\t\t 2. 发送消息")
	fmt.Println("\t\t 3. 信息列表")
	fmt.Println("\t\t 4. 退出系统")
	fmt.Println("\t\t 请选择（1-4）")
	var key int
	var content string
	//因为常用smsProcess 实例
	smsProcess := &SmsProcess{}

	fmt.Scanln(&key)
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("请输入消息")
		fmt.Scanln(&content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("查看信息列表")
	case 4:
		fmt.Println("退出系统")
		//向服务器端发送退出消息
		os.Exit(0)
	default:
		fmt.Println("输入不正确")
	}
}

//和服务器端保持通讯
func serverProcessMes(conn net.Conn) {

	tf := &utils.Transfer{
		Conn: conn,
	}

	for {
		fmt.Println("客户端正在等待服务器端发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() fail , err = ", err)
			return
		}
		//fmt.Println("mes = ", mes)

		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线了
			// 1. 取出 NoitfyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			err = json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			if err != nil {
				fmt.Println("processC | server.go | json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes) fail , err = ", err)
			}
			// 2. 将把这个用户的状态保存到客户map中
			updateUserStatus(&notifyUserStatusMes)
			outputOnlineUser()
		case message.SmsMesType: // 有人群发消息
			smsMgr := SmsMgr{}
			err = smsMgr.outputGroupMes(&mes)
			if err != nil {
				fmt.Println("processC | server | serverProcessMes | smsMgr.outputGroupMes(&mes)", err)
			}
		default:
			fmt.Println("服务器端返回未知类型")
		}

	}

}
