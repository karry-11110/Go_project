服务端实现：

    1.监听端口
        1.1 由于服务端要与很多客户端通信，因此这里定义一个服务端对象，包括一个在线用户的map
        1.2 定义一个构造函数API来构造server对象
        1.3 写一个start()函数来启动server对象，启动监听
            1.3.1 监听 启动listen socket
            1.3.2 defer 关闭listen socket
            1.3.3 创建goroutine监听server的通道是否有消息，有就立刻发送给全部的在线user
    2.接收客户端请求链接: 这里在一个无线循环里接受链接得到conn对象，可以接受很多个客户端
    3.创建goroutine处理链接conn(核心就是这个conn，围绕它来对服务端和客户端进行一切操作)
        3.1 通过链接conn信息，以及与客户端绑定的server对象来实现这个客户端对象的构造，来构建这个conn客户端与服务端的通信
            user.go:
            3.1.1 定义一个在线用户的结构体
            3.1.2 通过conn信息来定义一个构造用户的API函数，同时要启动一个gorotine监听但前用户是否收到了服务端的广播
        3.2 实现用户的上线业务：就是将用户加入到服务端的在线用户map中，并由服务端广播消息到其他客户端
            3.2.2 实现服务端向服务端通道广播消息的方法
        3.3 循环接收客户端发送的消息：针对消息进行处理
            3.3.1 如果接受的消息n数量为0，则用户下线
            3.3.2 提取消息进行处理,Domessage
        3.4 判断用户是否活跃，这里要在3.3处进行一个通道的flag通信,注意此时Handler协程是阻塞的，因为islive要在另外一个goroutine有值
            3.4.1 如果超时，用户要给当前对应的客户端发送消息为：由于您长时间潜水，现在系统强制下线您的帐号。
客户端实现：
    1.建立于服务端连接
        1.1 客户端结构体对象及其构造API实现
            1.1.1 这里借助flag包和init函数实现命令行解析，就不用把服务器ip和端口定死了
            1.1.2 在构造API里面实现与服务端的连接
    2.单独开一个goroutine去处理服务端的消息
    3.封装客户端业务
        3.1在客户端业务里面就实现了断开连接
        3.2 设置一个菜单实现功能选择