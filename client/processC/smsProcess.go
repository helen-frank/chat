package processC

import (
	"chat/client/utils"
	"chat/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

//发送群聊消息
func (sp *SmsProcess) SendGroupMes(content string) (err error) {
	// 1. 创建一个Mes
	mes := message.Message{
		Type: message.SmsMesType,
	}

	// 2. 创建一个SmsMes
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	// 3. 序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("ProcessC | smsProcess.go | SendGroupMes | json.Marshal(smsMes) fail , err = ", err)
		return
	}

	mes.Data = string(data)

	// 4. 对mes再次序列化

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("ProcessC | smsProcess.go | SendGroupMes | json.Marshal(mes) fail , err = ", err)
		return
	}

	// 5. 将mes 发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	// 6. 发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("ProcessC | smsProcess.go | SendGroupMes | tf.WritePkg(data) fail , err = ", err)
		return
	}
	return
}
