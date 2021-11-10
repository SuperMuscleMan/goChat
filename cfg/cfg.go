package cfg

const (
	BadWordsPath = "./cfg/list.txt"

	MsgTypeSys int8 = 1
	MsgTypeOne int8 = 2
	MsgTypeAll int8 = 3

	ChatLimit = 3
	ChatCache = 6

	SysMSGListenInfo         = "服务器监听地址："
	SysMsgTitle              = "【系统消息】 "
	SysMsgStartSuccess       = "聊天服务器启动成功！"
	SysMsgOnline             = " 上线了"
	SysMsgOnlineOff          = " 离线了"
	SysMsgSyntaxErr          = "语法错误，消息发送失败，请查阅使用说明！"
	SysMsgNotSameGroup       = "抱歉，非同一个聊天室，无法聊天！"
	SysMsgCreateGroupSuccess = " 聊天室创建成功！"
	SysMsgUserNameOccupy     = "抱歉，用户名已被占用"
	SysMsgFromGroup          = "来自聊天室 "
	SysMsgFromFriends        = "来自好友 "
	SysMsgNonExistGroup      = " 聊天室不存！"
	SysMsgNotJoinGroup       = "未加入聊天室！"
	SysMsgNotPopular         = "暂无复合条件的流行词！"
	SysMsgPopular            = "当前流行词："
	SysMsgNotExistUser       = "用户名不存在"
	SysMsgNotOnline          = "当前用户不在线"
	SysMsgJoinSuccess        = " 聊天室加入成功！"
)
