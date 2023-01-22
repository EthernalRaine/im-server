package oscar

type OSCARContext struct {
	UIN            int
	ServerSequence uint16
	ClientSequence uint16
	Challenge      []byte
	BOSCookie      []byte
}

type FLAPPacket struct {
	Frame      byte
	Sequence   uint16
	DataLength uint16
	Data       []byte
}

type SNACMessage struct {
	Foodgroup uint16
	Subgroup  uint16
	Flags     uint16
	RequestID uint32
	Data      []byte
}

type TLV struct {
	Type   uint16
	Length uint16
	Value  []byte
}

// Foodgroups
const (
	FoodgroupOSERVICE = 0x0001 // Generic service controls
	FoodgroupBUDDY    = 0x0003 // Buddy List management service
	FoodgroupICBM     = 0x0004 // ICBM (messages) service
	FoodgroupFEEDBAG  = 0x0013 // Server Side Information (SSI) service
	FoodgroupBUCP     = 0x0017 // Authorization/registration service
)

// FLAP
const (
	FrameSignOn  = 0x01
	FrameData    = 0x02
	FrameError   = 0x03
	FrameSignOff = 0x04
)

// BUCP
const (
	BUCPLoginRequest      = 0x0002
	BUCPLoginResponse     = 0x0003
	BUCPChallengeRequest  = 0x0006
	BUCPChallengeResponse = 0x0007
)

// OSERVICE
const (
	OSERVICEHostOnline     = 0x0003
	OSERVICEClientVersions = 0x0017
)

// http://web.archive.org/web/20211023214802fw_/http://iserverd.khstu.ru/oscar/families.html
var supportedFoodgroups []uint16 = []uint16{
	FoodgroupOSERVICE,
	FoodgroupBUDDY,
	FoodgroupICBM,
	FoodgroupFEEDBAG,
}

var ClientContexts []*OSCARContext
