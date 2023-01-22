package network

import (
	"chimera/utility/tcp"
)

type Client struct {
	Connection    tcp.TcpConnection
	ClientInfo    Details
	ClientAccount Account
	ClientUser    User
}

type Details struct {
	Service   int
	Messenger string
	Build     string
	Protocol  string
}

type Account struct {
	UIN         int
	DisplayName string
	Mail        string
	Password    string
}

type User struct {
	UIN             int
	AvatarBlob      string
	AvatarImageType string
	StatusCode      int
	StatusMessage   string
	LastLogin       int64
	SignupDate      int
}

type Contact struct {
	SenderUIN int
	FriendUIN int
	Reason    string
}

type OfflineMessage struct {
	SenderUIN      int
	RecvUIN        int
	MessageDate    int
	MessageContent string
}

type Meta struct {
	UIN         int
	UsageFlag   int
	AccountFlag int
}

type ServiceMessage struct {
	Service int
	Type    int
	Data    ServiceData
}

type ServiceData struct {
	Sender int
	Recv   int
	Raw    []byte
}

var Clients []*Client
var MessageCache []*ServiceMessage

// Service
const (
	Service_MSIM  = 0x0001
	Service_OSCAR = 0x0002
)

const (
	MessageType_Status    = 0x0001
	MessageType_SignOn    = 0x0002
	MessageType_LogOff    = 0x0003
	MessageType_IM        = 0x0004
	MessageType_OfflineIM = 0x0005
	MessageType_AddFriend = 0x0006
	MessageType_DelFriend = 0x0007
)

// Meta
const (
	UsageFlag_Normal    = 0x0001
	UsageFlag_Donator   = 0x0002
	UsageFlag_Beta      = 0x0003
	UsageFlag_Alpha     = 0x0004
	UsageFlag_Developer = 0x0005
)

const (
	AccountFlag_Normal   = 0x0001
	AccountFlag_Disabled = 0x0002
	AccountFlag_Banned   = 0x0003
	AccountFlag_Underage = 0x0004
	AccountFlag_Timeout  = 0x0005
)
