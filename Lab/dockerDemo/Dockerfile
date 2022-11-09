#******************************常规***********************
#FROM golang:latest
## 为我们的镜像设置必要的环境变量
#ENV GO111MODULE=on \
#    CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64
#
## 移动到工作目录：/go/dockerdemo
# #注意这里的地址不能就是/go,这是gopath路径，不能go build,但是可以go run
#WORKDIR /go/dockerdemo
## 将代码复制到容器中
#COPY . .
#
#
## 将我们的代码编译成二进制可执行文件app
#RUN go build -o app .
#
## 移动到用于存放生成的二进制文件的 /dist 目录
#WORKDIR /dist
#
## 将二进制文件从 /build 目录复制到这里
#RUN cp /go/dockerdemo/app .
#
## 声明服务端口
#EXPOSE 8888
#
## 启动容器时运行的命令
#CMD ["/dist/app"]


##**********************************分阶段构建镜像小很多************
#FROM golang:latest AS builder
#
## 为我们的镜像设置必要的环境变量
#ENV GO111MODULE=on \
#    CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64
#
## 移动到工作目录：/build
#WORKDIR /go/demo
#
## 将代码复制到容器中
#COPY . .
#
## 将我们的代码编译成二进制可执行文件 app
#RUN go build -o app .
#
####################
## 接下来创建一个小镜像
####################
#FROM scratch
#
## 从builder镜像中把/dist/app 拷贝到当前目录
#COPY --from=builder /go/demo/app /
#
## 需要运行的命令
#ENTRYPOINT ["/app"]

#**********************************分阶段构建镜像小很多************
FROM golang:latest AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /go/demo

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 app
RUN go build -o app .

###################
# 接下来创建一个小镜像
###################
FROM scratch

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /go/demo/app /

# 需要运行的命令
ENTRYPOINT ["/app"]