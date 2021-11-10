package connect

import (
	badWords2 "dosChat/badWords"
	"dosChat/cfg"
	"dosChat/popularWords"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)


// 用户结构
type userInfo struct {
	userName   string
	connect    net.Conn
	onlineTime int64
	groupList  []string
}

//聊天室结构
type groupInfo struct {
	groupName string
	member    []string
	chatList  []sendChanMsg
}

// 频道传送结构体
type sendChanMsg struct {
	receiverName string // “”表示广播给发送者的所有聊天室
	senderName   string
	msgType      int8
	msgContent   string
}

// 用户map、聊天室map
var userMaps = make(map[string]userInfo)
var groupMaps = make(map[string]groupInfo)

// 广播频道、私聊频道
var broadcastChan = make(chan sendChanMsg)
var sendChan = make(chan sendChanMsg)

// 处理每一个用户连接
func Handle(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		println("read err" + err.Error())
		return
	}
	name := string(buf[:n])
	isAlready := online(name, conn)
	if isAlready {
		_, err := conn.Write([]byte(cfg.SysMsgUserNameOccupy))

		if err != nil {
			fmt.Println("sendHistory err" + err.Error())
		}
		err = conn.Close()

		if err != nil {
			fmt.Println("sendHistory err" + err.Error())
		}
		return
	}
	defer onlineOff(name)
	broadcastChan <- sendChanMsg{"", name, cfg.MsgTypeSys, name + cfg.SysMsgOnline}
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		// 判断GM
		body := string(buf[:n])
		handleData(name, body)
	}
}

func handleData(name string, body string) {
	if !checkGM(name, body) {
		// 判断创建聊天室
		if !checkCreateGroup(name, body) {
			//	私聊
			if !checkSendOne(name, body) {
				// 聊天室
				if !checkSendAll(name, body) {
					// 发送错误
					sendChan <- sendChanMsg{name, name, cfg.MsgTypeSys, cfg.SysMsgSyntaxErr}
				}
			}
		}
	}
}



// 发送消息给指定用户
func send(memberName string, msgInfo sendChanMsg) {
	memberInfo := userMaps[memberName]
	if memberInfo.onlineTime > 0 {
		_, err := memberInfo.connect.Write([]byte(msgInfo.msgContent))
		if err != nil {
			fmt.Println("sendHistory err" + err.Error())
		}
	}
}

// 处理系统消息
func msgType(msgInfo *sendChanMsg) {
	switch msgInfo.msgType {
	case cfg.MsgTypeSys:
		msgInfo.msgContent = cfg.SysMsgTitle + msgInfo.msgContent
	case cfg.MsgTypeOne:
		msgInfo.msgContent = cfg.SysMsgFromFriends + msgInfo.senderName + "：" + msgInfo.msgContent
	case cfg.MsgTypeAll:
		msgInfo.msgContent = cfg.SysMsgFromGroup + msgInfo.receiverName + " - " + msgInfo.senderName + "：" + msgInfo.msgContent
	}
}


// 私聊同一个聊天室内的用户
func checkSendOne(sendName string, str string) bool {
	//@UserName Msg
	if str[0] == '@' {
		field := strings.Fields(str[1:])
		receiveName := field[0]
		if receiveInfo, ok := userMaps[receiveName]; ok {
			if receiveInfo.onlineTime == 0 {
				sendChan <- sendChanMsg{sendName, sendName, cfg.MsgTypeSys, cfg.SysMsgNotOnline}
				return true
			}
			msg := strings.Join(field[1:], " ")
			msg = badWords2.HandelBad(msg)
			if checkRelation(sendName, receiveName) {
				sendChan <- sendChanMsg{receiveName, sendName, cfg.MsgTypeOne, msg}
				return true
			} else {
				sendChan <- sendChanMsg{sendName, sendName, cfg.MsgTypeSys, cfg.SysMsgNotSameGroup}
				return true
			}
		} else {
			sendChan <- sendChanMsg{sendName, sendName, cfg.MsgTypeSys, cfg.SysMsgNotExistUser}
			return true
		}
	}
	return false
}

// 聊天室发消息
func checkSendAll(sendName string, str string) bool {
	//*GroupName Msg
	if str[0] == '*' {
		field := strings.Fields(str[1:])
		groupName := field[0]
		if isExistGroup(sendName, groupName) {
			msg := strings.Join(field[1:], " ")
			msg = badWords2.HandelBad(msg)
			popularWords.Statistic(field[1:])
			broadcastChan <- sendChanMsg{groupName, sendName, cfg.MsgTypeAll, msg}
		}
		return true
	}
	return false
}

// 检测是否存在聊天室
func isExistGroup(sendName string, groupName string) bool {
	if groupInfo, ok := groupMaps[groupName]; ok {
		for _, r := range groupInfo.member {
			if r == sendName {
				return true
			}
		}
		sendChan <- sendChanMsg{sendName, sendName, cfg.MsgTypeSys, cfg.SysMsgNotJoinGroup}
		return false
	}
	sendChan <- sendChanMsg{sendName, sendName, cfg.MsgTypeSys, groupName + cfg.SysMsgNonExistGroup}
	return false
}

// 检测关系
func checkRelation(sendName string, receiveName string) bool {
	sendInfo := userMaps[sendName]
	receiveInfo := userMaps[receiveName]
	for _, r := range sendInfo.groupList {
		for _, e := range receiveInfo.groupList {
			if r == e {
				return true
			}
		}
	}
	return false
}

// 创建聊天室
func checkCreateGroup(senderName string, str string) bool {
	//add#GroupName=UserName1+UserName2+UserName3
	reg := regexp.MustCompile("add#(.+)=")
	result := reg.FindStringSubmatch(str)
	if len(result) < 2 {
		return false
	}
	groupName := result[1]
	splitR := strings.Split(str, "=")
	if len(splitR) < 2 {
		return false
	}
	groupInfo, ok := groupMaps[groupName]
	msg := sendChanMsg{groupName, senderName, cfg.MsgTypeSys, groupName + cfg.SysMsgCreateGroupSuccess}
	if ok {
		msg.msgContent = groupName + cfg.SysMsgJoinSuccess
	}
	msgType(&msg)
	for _, e := range strings.Split(splitR[1], "+") {
		if !isExist(groupInfo.member, e) {
			memberInfo := userMaps[e]
			memberInfo.groupList = append(memberInfo.groupList, groupName)
			userMaps[e] = memberInfo
			groupInfo.member = append(groupInfo.member, e)
			//
			//broadcastChan <- sendChanMsg{groupName, senderName, msgTypeSys, groupName + sysMsgCreateGroupSuccess}
			send(memberInfo.userName, msg)
			sendHistory(memberInfo.connect, groupInfo)
		}
	}
	groupInfo.groupName = groupName
	groupMaps[groupName] = groupInfo

	return true
}

// 判断是否存在数组中
func isExist(strArr []string, ele string) bool {
	for _, e := range strArr {
		if e == ele {
			return true
		}
	}
	return false
}

// 检测是否GM指令
func checkGM(userName string, str string) bool {
	if str[0] == '/' {
		field := strings.Fields(str[1:])
		if len(field) != 2 {
			return false
		}
		gm := field[0]
		switch gm {
		case "popular":
			second, _ := strconv.Atoi(field[1])
			if second > popularWords.Expire {
				return false
			}
			popular, _ := popularWords.GetPopular(int64(second))
			if popular == "" {
				sendChan <- sendChanMsg{userName, userName, cfg.MsgTypeSys, cfg.SysMsgNotPopular}
			} else {
				sendChan <- sendChanMsg{userName, userName, cfg.MsgTypeSys, cfg.SysMsgPopular + popular}
			}
		case "stats":
			username := field[1]
			tips, err := getOnlineTime(username)
			if !err {
				sendChan <- sendChanMsg{userName, userName, cfg.MsgTypeSys, tips}
			} else {
				sendChan <- sendChanMsg{userName, userName, cfg.MsgTypeSys, tips}
			}
		default:
			return false
		}
		return true
	}
	return false
}

// 获取在线时长
func getOnlineTime(name string) (string, bool) {
	if userInfo, ok := userMaps[name]; ok {
		if userInfo.onlineTime == 0 {
			return cfg.SysMsgNotOnline, false
		}
		curTime := time.Now().Unix()
		diff := curTime - userInfo.onlineTime
		d := diff / 86400
		diff = diff % 86400
		h := diff / 3600
		diff = diff % 3600
		m := diff / 60
		s := diff % 60
		return fmt.Sprintf("%02dd %02dh %02dm %02ds ", d, h, m, s), true
	}
	return cfg.SysMsgNotExistUser, false
}


// 存储聊天记录
func saveChat(msgInfo sendChanMsg) {
	groupName := msgInfo.receiverName
	groupInfo := groupMaps[groupName]

	groupInfo.chatList = append(groupInfo.chatList, msgInfo)
	if len(groupInfo.chatList) > cfg.ChatCache {
		groupInfo.chatList = groupInfo.chatList[cfg.ChatCache-cfg.ChatLimit:]
	}
	groupMaps[groupName] = groupInfo
}

// 下发历史记录50条
func sendHistory(conn net.Conn, group groupInfo) {
	max := len(group.chatList)
	if max <= cfg.ChatLimit {
		for _, r := range group.chatList {
			_, err := conn.Write([]byte(r.msgContent))
			if err != nil {
				fmt.Println("sendHistory err" + err.Error())
			}
		}
	} else {
		for _, r := range group.chatList[max-cfg.ChatLimit:] {
			_, err := conn.Write([]byte(r.msgContent))
			if err != nil {
				fmt.Println("sendHistory err" + err.Error())
			}
		}
	}
}




// 接收群聊处理（该函数处于goroutine机制
func BroadcastHandle() {
	for {
		msgInfo := <-broadcastChan
		msgType(&msgInfo)
		if msgInfo.receiverName == "" {
			for _, userInfo := range userMaps {
				if userInfo.onlineTime > 0 {
					_, err := userInfo.connect.Write([]byte(msgInfo.msgContent))
					if err != nil {
						fmt.Println("sendHistory err" + err.Error())
					}
				}
			}
			//userInfo := userMaps[msgInfo.senderName]
			//for _, groupName := range userInfo.groupList {
			//	sendGroup(groupName, msgInfo)
			//}
		} else {
			sendGroup(msgInfo.receiverName, msgInfo)
		}
	}
}
// 发送消息给指定聊天室的所有成员
func sendGroup(groupName string, msgInfo sendChanMsg) {
	groupInfo := groupMaps[groupName]

	saveChat(msgInfo)

	for _, memberName := range groupInfo.member {
		if msgInfo.msgType == cfg.MsgTypeAll && msgInfo.senderName == memberName {
			continue
		} else
		{
			send(memberName, msgInfo)

		}
	}
}

// 接收私聊处理（该函数处于goroutine机制
func SendHandle() {
	for {
		msgInfo := <-sendChan
		msgType(&msgInfo)
		send(msgInfo.receiverName, msgInfo)
	}

}

// 上线逻辑
func online(name string, conn net.Conn) bool {
	userInfo, _ := userMaps[name]
	if userInfo.onlineTime > 0 {
		return true
	} else {
		userInfo.userName = name
		userInfo.connect = conn
		userInfo.onlineTime = time.Now().Unix()
		userMaps[name] = userInfo
		return false
	}
}

//离线逻辑
func onlineOff(name string) {
	userInfo, _ := userMaps[name]
	userInfo.onlineTime = 0
	userMaps[name] = userInfo
	broadcastChan <- sendChanMsg{"", name, cfg.MsgTypeSys, name + cfg.SysMsgOnlineOff}
}
