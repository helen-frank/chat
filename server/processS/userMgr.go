package processS

import "fmt"

//因为UserMgr实例在服务器端有且只有一个，且在很多地方都会用到
//因此，将其定义为全局变量

var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对 onlineUsers 添加
func (um *UserMgr) AddOnlineUser(up *UserProcess) {
	um.onlineUsers[up.UserId] = up
}

//删除
func (um *UserMgr) DelOnlineUser(userId int) {
	delete(um.onlineUsers, userId)
}

//返回当前所有在线用户
func (um *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return um.onlineUsers
}

//根据id返回对应的值
func (um *UserMgr) GetOnlineById(UserId int) (up *UserProcess, err error) {
	up, ok := um.onlineUsers[UserId]
	if !ok { //说明当前用户不在线
		err = fmt.Errorf("用户%d 不在线", UserId)
		return
	}
	return
}
