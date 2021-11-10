# 说明文档



## 设计思路

1. 采用c/s架构，包括服务端、客户端。
2. 用户操作，均在命令行中进行。
3. 网络协议
   1. 传输协议采用了TCP协议
   2. 应用协议由字符判断识别
4. 客户端，拥有收发消息功能。分别由两个goroutine协程进行处理。
5. 服务端，拥有多个goroutine协程，如下：
   1. 接收tcp连接协程
      1. 无限循环接收客户端的连接请求
   2. 连接成功后的用户处理协程
      1. 处理客户端发送的所有请求，（例如：发消息、创建聊天室、GM指令等
   3. 消息广播协程
   4. 私聊消息协程
   5. 流行词协程
      1. 定时清理过期的流行词

## 关键算法

1. 屏蔽字检测
   1. 采用的字典树结构
   2. 原理：
      1. 节点结构：{node：{isEnd: false,  nextNode: NODE}}
      2. 初始化字典树，将所有屏蔽词依次放入字典树中
      3. 检测字符串，例如：”abc“
         1. 首先检索树根下是否存在字符 a 节点
         2. 若有，再依次检索b、c
         3. 最终只检索到了a、b两个节点时，判断b节点的isEnd标识是否真，若为真则“abc” 字符中“ab”是屏蔽字，替换为“**c"
   3. 优点：共享公共的前缀，占用空间少



## 部署

#### 服务器端

1. 安装golang环境
2. 编译server.go文件
   1. 命令：`go build server.go`
3. 启动server.exe程序
  1. 命令：`server.exe -add 127.0.0.1:12312`

#### 客户端

1. 安装golang环境
2. 编译client.go文件
   1. 命令：`go build client.go`
3. 启动client.exe程序
   1. 命令：`client.exe -add 127.0.0.1:12312 -u USERNAME`

## 使用说明

1. 登录程序
   1. 命令：`client.exe -add ADD -u USERNAME`  
      1. ADD: 服务器监听地址；
      2. USERNAME: 用户名称;
2. 退出程序
   1. 命令： `quit`
3. 创建聊天室
   1. 命令：`add#GROUPNAME=USERNAME1+USERNAME2`
      1. GROUPNAME: 聊天室名称
      2. USERNAME1: 加入聊天室的用户1
      3. USERNAME2: 加入聊天室的用户2
4. 加入聊天室
   1. 命令：`add#GROUPNAME=USERNAME3+USERNAME4`
      1. 表示用户3和用户4加入GROUPNAME聊天室
5. 发送消息至聊天室
   1. 命令：`*GROUPNAME MESSAGE` 
      1. GROUPNAME: 聊天室名称
      2. MESSAGE: 消息内容
   2. 注意：
      1. 命令第一个字符是星号
      2. GOURPNAME字段后有一个空格
6. 发送消息至好友
   1. 命令：`@USERNAME MESSAGE`
      1. USERNAME：用户名称
      2. MESSAGE：消息内容
7. GM指令
   1. 获取流行词
      1. 命令：`/popular N`
         1. 打印出最近N秒内发送频率最高的词，N小于等于60
   2. 获取用户在线时长
      1. 命令：`/stats USERNAME`
         1. 打印出USERNAME用户在线时长。
         2. 时间格式：00d 00h 00m 00s



