package connect

import (
	"dosChat/cfg"
	"regexp"
	"testing"
	"time"
)
// 初始化一个用户数据，用于后续测试
func TestOnline(t *testing.T) {
	name := "UserName"
	userInfo, _ := userMaps[name]
	userInfo.userName = name
	userInfo.onlineTime = time.Now().Unix() - 89731
	userMaps[name] = userInfo
}
// 测试获取在线时间，判断格式是否正确
func TestGetOnlineTime(t *testing.T) {
	timeStr, err := getOnlineTime("UserName")
	if !err {
		t.Error(`getOnlineTime("UserName")=false`)
	}
	reg, _ := regexp.Compile(`\d{2}d\s\d{2}h\s\d{2}m\s\d{2}s`)
	if !reg.MatchString(timeStr) {
		t.Error(`getOnlineTime("UserName")=false`)
	}
}

// 测试不同类型消息内容的拼接
func TestMsgType(t *testing.T) {
	oldChanMsg := sendChanMsg{"receiverName", "senderName", cfg.MsgTypeSys, "msgContent"}
	chanMsg1 := oldChanMsg
	msgType(&chanMsg1)
	if chanMsg1.msgContent != (cfg.SysMsgTitle+oldChanMsg.msgContent) {
		t.Error(`msgType(sendChanMsg{"receiverName", "senderName", cfg.MsgTypeSys, "msgContent"})`)
	}
	chanMsg2 := oldChanMsg
	chanMsg2.msgType = cfg.MsgTypeOne
	msgType(&chanMsg2)
	if chanMsg2.msgContent !=( cfg.SysMsgFromFriends+oldChanMsg.senderName+"："+oldChanMsg.msgContent) {
		t.Error(`msgType(sendChanMsg{"receiverName", "senderName", cfg.MsgTypeOne, "msgContent"})`)
	}
	chanMsg3 := oldChanMsg
	chanMsg3.msgType = cfg.MsgTypeAll
	msgType(&chanMsg3)
	if chanMsg3.msgContent != cfg.SysMsgFromGroup+oldChanMsg.receiverName+" - "+oldChanMsg.senderName+"："+oldChanMsg.msgContent {
		t.Error(`msgType(sendChanMsg{"receiverName", "senderName", cfg.MsgTypeAll, "msgContent"})`)
	}

}
