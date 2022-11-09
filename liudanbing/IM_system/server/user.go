package main

import (
	"net"
	"strings"
)

// 3.1.1 定义一个在线用户的结构体
type User struct {
	Name    string
	Address string
	C       chan string
	conn    net.Conn
	server  *Server
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

// 3.1.2 通过conn信息来定义一个构造用户的API函数，同时要启动一个gorotine监听但前用户是否收到了服务端的广播
func NewUser(conn net.Conn, server *Server) *User {
	userAddress := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddress,
		Address: userAddress,
		C:       make(chan string),
		conn:    conn,
		server:  server,
	}
	go user.ListenMessage()
	return user
}

func (u *User) Online() {
	u.server.MapLock.Lock()
	u.server.OnlineUserMap[u.Name] = u
	u.server.MapLock.Unlock()
	u.server.BroadCast(u, "已上线！！！")
}
func (u *User) Offline() {
	u.server.MapLock.Lock()
	delete(u.server.OnlineUserMap, u.Name)
	u.server.MapLock.Unlock()
	u.server.BroadCast(u, "下线了！！！")
}
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}
func (u *User) DoMessage(msg string) {
	if msg == "look" { //查询当前在线用户都有哪些
		u.server.MapLock.Lock()
		for _, user := range u.server.OnlineUserMap {
			onlineMsg := "[" + user.Address + "]" + user.Name + ":" + "在线！！！\n"
			u.SendMsg(onlineMsg)
		}
		u.server.MapLock.Unlock()

	} else if msg[:7] == "rename|" { //改用户名字：消息格式: rename|张三
		newName := strings.Split(msg, "|")[1]
		_, ok := u.server.OnlineUserMap[newName]
		if ok {
			u.SendMsg("当前用户名已被使用\n")
		} else {
			u.server.MapLock.Lock()
			delete(u.server.OnlineUserMap, u.Name) //map删除键值对的方式
			u.server.OnlineUserMap[newName] = u    //增加一个新键值对只能这么操作
			u.server.MapLock.Unlock()
			u.Name = newName
			u.SendMsg("您已经更新用户名为：" + u.Name + "\n")
		}

	} else if msg[:3] == "to|" { //给某客户端特定私信
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("消息格式不正确，请使用\"to|张三|你好！\"格式。\n")
			return
		}
		remoteUser, ok := u.server.OnlineUserMap[remoteName]
		if !ok {
			u.SendMsg("该用户名不存在！！\n")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("消息内容为空，请重新发送\n")
			return
		}
		remoteUser.SendMsg(u.Name + "发消息说：" + content)
	} else { //给全体用户广播消息
		u.server.BroadCast(u, msg)
	}
}
