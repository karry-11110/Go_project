package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//1.1 由于服务端要与很多客户端通信，因此这里定义一个服务端对象，包括一个在线用户的map
type Server struct {
	Ip            string
	Port          int
	OnlineUserMap map[string]*User
	MapLock       sync.RWMutex
	Message       chan string
}

// 1.2 定义一个构造函数API来构造server对象
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:            ip,
		Port:          port,
		OnlineUserMap: make(map[string]*User),
		Message:       make(chan string),
	}
}
func (s *Server) ListenMessager() {
	msg := <-s.Message
	s.MapLock.Lock()
	for _, client := range s.OnlineUserMap {
		client.C <- msg
	}
	s.MapLock.Unlock()

}

//3.2.2 实现服务端向服务端通道广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// 3.1 通过链接conn信息，以及与客户端绑定的server对象来实现这个客户端对象的构造,来构建这个conn客户端与服务端的通信
	user := NewUser(conn, s)
	//3.2 实现用户的上线业务：就是将用户加入到服务端的在线用户map中，并由服务端广播消息到其他客户端
	user.Online()
	//监听用户是否活跃的channel
	isLive := make(chan bool)
	//3.3 循环接收客户端发送的消息：针对消息进行处理
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			// 3.3.1 如果接受的消息n数量为0，则用户下线
			if n == 0 {
				user.Offline()
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			//3.3.2 提取消息进行处理,Domessage
			msg := string(buf[:n-1]) //提取用户的消息(去除'\n')
			user.DoMessage(msg)
			isLive <- true ////用户的任意消息，代表当前用户是一个活跃的
		}
	}()
	//3.4 判断用户是否活跃，这里要在3.3处进行一个通道的flag通信,注意此时Handler协程是阻塞的，因为islive要在另外一个goroutine有值
	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 300):
			//3.4.1 已经超时
			//将当前的User强制的关闭

			user.SendMsg("由于您长时间潜水，现在系统强制下线您的帐号！！！")

			//销毁用的资源
			close(user.C)

			//关闭连接
			conn.Close()

			//退出当前Handler
			return //runtime.Goexit()

		}
	}
}

//1.3 写一个start()函数来启动server对象，启动监听
func (s *Server) Start() {
	//1.3.1 监听 启动listen socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	// 1.3.2 defer 关闭listen socket
	defer listener.Close()
	//1.3.3 创建goroutine监听server的通道是否有消息，有就立刻发送给全部的在线user
	go s.ListenMessager()
	//2.接收客户端请求链接,在一个无线循环里接受链接得到conn对象可以接受很多个客户端
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		fmt.Println(conn.RemoteAddr().String())
		// 3.创建goroutine处理链接conn
		go s.Handler(conn)

	}
}

func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()

}
