package bridge

import (
	"chimera/utility/encryption"
	"fmt"
	"strconv"
	"sync"
)

func BridgeNewTransportHeader(sender int, recv int) BridgeMessageTransportHeader {
	return BridgeMessageTransportHeader{SenderUIN: sender, RecvUIN: recv}
}

func BridgeNewActionDataCommand(key string, value string) BridgeMessageActionData {
	return BridgeMessageActionData{Key: key, Value: value}
}

func BridgeBuildDataPackage(header BridgeMessageTransportHeader, actiondata []BridgeMessageActionData) string {
	data := ""

	/*
		x1c is used as the seperator between the 3 data parts
		xc0x80 is used as the split for actiondata pairs
		all of those are bytes that should not be in use since the protocol is text based
	*/

	// build headerdata
	data += fmt.Sprintf("S:%d|R:%d\x1c", header.SenderUIN, header.RecvUIN)

	// build actiondata
	for ix := 0; ix < len(actiondata); ix++ {
		data += fmt.Sprintf("\xc0\x80%s\xc0\x80%s", actiondata[ix].Key, actiondata[ix].Value)
	}
	data += "\x1c"

	// build checksum
	crc := encryption.GetCRCHash([]byte(data))
	data += strconv.FormatUint(uint64(crc), 10)

	return data
}

var sendDataMutex sync.RWMutex

func BridgeSendData(data string) {
	sendDataMutex.Lock()
	messages = append(messages, data)
	sendDataMutex.Unlock()
}
