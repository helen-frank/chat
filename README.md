#  简介

这是一个海量即时通讯系统，使用go net tcp 开发

数据库用到了redis的hash,传输数据为json序列化后的字符

具备客户端和服务器端

# 开发环境

| language          | go version go1.16.4            |
| ----------------- | ------------------------------ |
| os                | archlinux/amd64 5.12.8-arch1-1 |
| development tools | vscode                         |
| sql               | redis                          |

# 目前已完成

> 1. 登录注册
> 2. 登录显示当前在线用户
> 3. 其他人上线更新在线用户
> 4. 向在线用户群发消息

下一步将开发功能

> 1. 下线更新在线用户
> 2. 私聊
> 3. 消息存储在数据库中
> 4. 向所有用户发送消息
> 5. 显示自己已收到的消息
> 6. hash混淆加密信息流

