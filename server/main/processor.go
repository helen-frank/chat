package main

import (
	"chat/common/message"
	"chat/server/processS"
	"chat/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// serverProcessMes 函数
// 功能：根据客户端发送消息种类的不同，决定调用哪个函数来处理
func (pr *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:

		// 创建一个UserProcessLogin实例
		up := &processS.UserProcess{
			Conn: pr.Conn,
		}

		err = up.ServerProcessLogin(mes)
		if err != nil {
			fmt.Println("up.ServerProcessLogin(mes) fail , err = ", err)
		}
	case message.RegisterMesType:
		up := &processS.UserProcess{
			Conn: pr.Conn,
		}
		err = up.ServerProcessRegister(mes)
		if err != nil {
			fmt.Println("up.ServerProcessRegister(mes) fail , err = ", err)
		}
	case message.SmsMesType:

		smsProcess := &processS.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在，无法处理")
	}
	return
}

func (pr *Processor) process2() (err error) {
	for {

		// 这里将读取数据包，直接封装成一个函数 readPkg(), 返回 message,err
		//创建一个TransFer 实例完成读包任务
		tf := &utils.Transfer{
			Conn: pr.Conn,
		}
		mes, err := tf.ReadPkg()

		if err != nil {
			if err == io.EOF {
				return err
			}
			fmt.Println(err)
		}
		fmt.Println("mes = ", mes)
		err = pr.serverProcessMes(&mes)
		if err != nil {
			fmt.Println("serverProcessMes(this.Conn, &mes) fail , err = ", err)
		}
	}

}
