package bridge

type BridgeClient struct {
	SerivceName string
	ServiceRev  string
}

type BridgeMessageTransportHeader struct {
	SenderUIN int
	RecvUIN   int
}

type BridgeMessageActionData struct {
	Key   string
	Value string
}

var clients []*BridgeClient
var messages []string
