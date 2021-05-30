package processC

import (
	"chat/common/message"
	"encoding/json"
	"fmt"
)

type SmsMgr struct {
}

func (sm *SmsMgr) outputGroupMes(mes *message.Message) (err error) {
	// 1. 反序列化 mes.Data
	var smsMes message.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("processC | smsMgr.go | outputGroupMes | json.Unmarshal([]byte(mes.Data), &smsMes) fail , err = ", err)
		return
	}

	// 2. 显示信息
	info := fmt.Sprintf("用户id： %d\n用户名：%s 对所人说：\n\t%s\n", smsMes.UserId, smsMes.UserName, smsMes.Content)
	fmt.Println(info)
	return
}
