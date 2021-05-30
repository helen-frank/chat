package processS

import (
	"chat/common/message"
	"chat/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

//转发消息
func (sp *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的 onlineUsers map[int]*UserProcess     后期改为遍历数据库，或者存在数据库中，让所有人都能看到
	//将消息转发取出
	var smsMes message.SmsMes

	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("ProcessS | smsProcess.go | SendGroupMes | json.Unmarshal([]byte(mes.Data),&smsMes) fail , err = ", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("ProcessS | smsProcess.go | SendGroupMes | json.Marshal(mes) fail , err = ", err)
		return
	}
	for id, up := range userMgr.onlineUsers {

		if id == smsMes.UserId {
			continue
		}
		sp.SendMesToEachOnlineUser(data, up.Conn)
	}

}

func (sp *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建 Transfer 实例， 发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("ProcessS | smsProcess.go | SendMesToEachOnlineUser | tf.WritePkg(data) fail , err = ", err)
	}

}
