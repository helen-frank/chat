package model

import (
	"chat/common/message"
	"net"
)

//因为在客户端，很多地方都

type CurUser struct {
	Conn net.Conn
	message.User
}
