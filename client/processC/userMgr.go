package processC

import (
	"chat/client/model"
	"chat/common/message"

	"fmt"
)

// 客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser //在用户登录后，完成对CurUser的初始化

//显示当前在线的用户
func outputOnlineUser() {

	fmt.Println("当前在线用户列表")
	for id := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
}

// 处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	//适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok { //原来没有
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
		fmt.Println("用户", notifyUserStatusMes.UserId, "上线")
	}

	user.UserStatus = notifyUserStatusMes.UserStatus

	onlineUsers[notifyUserStatusMes.UserId] = user

}
