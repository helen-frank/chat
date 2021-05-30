package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

//用户在线常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息
}

type LoginMes struct {
	UserId   int    `json:"userId"`   //用户id
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code int `json:"code"` // 返回状态码
	//200 表示登录成功
	//403 表示密码不正确
	//500 表示该用户未注册
	//505 表示服务器内部错误
	UsersId []int  `json:"usersId"` //保存用户id的切片
	Error   string `json:"error"`   // 返回错误信息
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	Code int `json:"code"`
	//200 表示注册成功
	//400 表示该用户已存在
	Error string `json:"error"`
}

//为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId     int `json:"userId"` //用户id
	UserStatus int `json:"status"` //用户的状态
}

// 发送的消息
type SmsMes struct {
	Content string `json:"content"`
	User           //匿名结构体
}
