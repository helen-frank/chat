package utils

import (
	"chat/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 这里将这些方法封装关联到结构体中
type Transfer struct {
	// 分析应该有哪些字段

	Conn net.Conn
	//Buf  [8096]byte //传输时使用的缓冲
}

func (tf *Transfer) ReadPkg() (mes message.Message, err error) {
	Buf := make([]byte, 8192)
	fmt.Println("读取客户端发送的数据ing")
	_, err = tf.Conn.Read(Buf[:4])
	if err != nil {
		fmt.Println("this.Conn.Read(this.Buf[:4]) fail , err = ", err)
		return
	}
	//根据buf[0:4] 转成一个 uint32 类型
	pkgLen := binary.BigEndian.Uint32(Buf[0:4])

	//根据 pkgLen 读取消息内容

	n, err := tf.Conn.Read(Buf[:pkgLen])

	if n != int(pkgLen) || err != nil {
		fmt.Println("this.Conn.Read(this.Buf[:pkgLen]) fail , err = ", err)
		return
	}

	// 把 pkgLen 反序列成 message.Message

	err = json.Unmarshal(Buf[:pkgLen], &mes)

	if err != nil {
		fmt.Println("json.Unmarshal(this.Buf[:pkgLen]) fail , err = ", err)
		return
	}
	return
}

func (tf *Transfer) WritePkg(data []byte) (err error) {
	// 先发送一个长度给对方
	pkgLen := uint32(len(data))
	var Buf [4]byte

	binary.BigEndian.PutUint32(Buf[0:4], pkgLen)
	// 发送长度
	_, err = tf.Conn.Write(Buf[:])
	if err != nil {
		fmt.Println("this.Conn.Write(bytes) , err = ", err)
		return
	}

	n, err := tf.Conn.Write(data)

	if n != int(pkgLen) || err != nil {
		fmt.Println("this.Conn.Write(data) fail , err = ", err)
		return
	}
	return

}
