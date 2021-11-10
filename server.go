package main

import (
	badWords2 "dosChat/badWords"
	"dosChat/cfg"
	connect2 "dosChat/connect"
	"dosChat/popularWords"
	"flag"
	"fmt"
	"net"
)

func main() {
	var address string
	flag.StringVar(&address, "add", "0.0.0.0:12312", "TCP listening address")
	flag.Parse()
	listener, err := net.Listen("tcp", address)

	if err != nil {
		panic("监听报错" + err.Error())
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Println("sendHistory err" + err.Error())
		}
	}()
	go popularWords.Timer()
	go connect2.BroadcastHandle()
	go connect2.SendHandle()
	badWords2.Init(cfg.BadWordsPath)
	fmt.Println(cfg.SysMsgTitle + cfg.SysMsgStartSuccess)
	fmt.Println(cfg.SysMsgTitle + cfg.SysMSGListenInfo + address)
	// 循环接收tcp连接请求
	for {
		connect, err := listener.Accept()
		if err != nil {
			fmt.Println("sendHistory err" + err.Error())
		} else {
			go connect2.Handle(connect)
		}
	}

}
