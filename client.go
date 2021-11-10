package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var address = flag.String( "add", "127.0.0.1:12312", "TCP listening address")
	var userName = flag.String( "u", "default", "User Name")
	flag.Parse()
	conn, _ := net.Dial("tcp", *address)
	buf := make([]byte, 1024)
	_, err := conn.Write([]byte(*userName))
	if err != nil {
		fmt.Println("【警告】连接服务器失败，请检查连接地址！" + err.Error())
		os.Exit(0)
	}
	go scan(conn)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("【警告】连接已断开，获取数据失败！" + err.Error())
			os.Exit(0)
		}
		fmt.Println(string(buf[:n]))
	}
}

func scan(conn net.Conn) {
	for {
		newScan := bufio.NewReader(os.Stdin)
		buf, _, _ := newScan.ReadLine()
		if string(buf) == "quit" {
			os.Exit(0)
		}
		_, err := conn.Write(buf)
		if err != nil {
			fmt.Println("【警告】连接异常！" + err.Error())
			os.Exit(0)
		}
	}
}
