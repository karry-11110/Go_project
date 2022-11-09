package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

//1.1 客户端结构体对象及其构造API实现
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	// 1.1.2 在构造API里面实现与服务端的连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dail err:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	//一旦client.conn有数据，就直接copy到stdout标准输出上, 永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	//  3.1在客户端业务里面就实现了断开连接
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("输入合法选项")
		return false
	}
}
func (client *Client) SelectUsers() {
	sendMsg := "look\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容, exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>请输入消息内容, exit退出:")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
		fmt.Scanln(&remoteName)
	}
}
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}
func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return false
	}
	return true
}
func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {
		}
		switch client.flag {
		case 1:
			client.PublicChat()

		case 2:
			client.PrivateChat()

		case 3:
			client.UpdateName()

		}
	}
}

//      1.1.1 这里借助flag包和init函数实现命令行解析，就不用把服务器ip和端口定死了
var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认是8888）")
}
func main() {
	flag.Parse() //实现1.1.1的命令行解析
	//1.建立于服务端连接
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败？？？？？？？？？？？")
		return
	}
	//2.单独开一个goroutine去处理服务端的消息
	go client.DealResponse()
	fmt.Println("连接服务器成功！！！！！！！！！！！！！！！")

	//3.封装客户端业务
	client.Run()
}
