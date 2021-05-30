package model

type User struct {
	//为了序列化和反序列成功
	//用户信息的json字符串的key 需和结构体的字段对应的 tag名字一致
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
