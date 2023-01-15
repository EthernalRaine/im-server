package bridge

import (
	"chimera/utility/encryption"
	"chimera/utility/logging"
	"fmt"
	"strconv"
	"strings"
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

func BridgeSplitDataPacket(datastream string) []string {
	decodedPacket := datastream
	splits := strings.Split(decodedPacket, "\x1c")

	return splits
}

func BridgeRetrieveValueFromActionDataPair(actionkey string, datastream string) string {
	decodedPacket := BridgeSplitDataPacket(datastream)[1] // index 1 is action data
	splits := strings.Split(decodedPacket, "\xc0\x80")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == actionkey {
			return splits[ix+1]
		}
	}

	return ""
}

func BridgeRetrieveTransportHeaderInformation(datastream string) BridgeMessageTransportHeader {
	decode := BridgeSplitDataPacket(datastream)[0] // index 0 is header data
	uins := strings.Split(decode, "|")

	s_uin, err := strconv.Atoi(strings.Trim(uins[0], "S:"))

	if err != nil {
		logging.Error("Bridge/BridgeRetrieveTransportHeaderInformation", "Failed to decode SenderUIN from Header (%s)", err.Error())
		return BridgeMessageTransportHeader{}
	}

	r_uin, err := strconv.Atoi(strings.Trim(uins[1], "R:"))

	if err != nil {
		logging.Error("Bridge/BridgeRetrieveTransportHeaderInformation", "Failed to decode RecvUIN from Header (%s)", err.Error())
		return BridgeMessageTransportHeader{}
	}

	header := BridgeMessageTransportHeader{
		SenderUIN: s_uin,
		RecvUIN:   r_uin,
	}
	return header
}

func BridgeRetrieveDataCRC(datastream string) uint32 {
	decode := BridgeSplitDataPacket(datastream)[2] // index 2 is crc
	crc, err := strconv.ParseUint(decode, 10, 32)

	if err != nil {
		logging.Error("Bridge/BridgeRetrieveDataCRC", "Failed to decode CRC32 from datastream (%s)", err.Error())
		return 0xDEADBEEF
	}

	return uint32(crc)
}

func BridgeFormatActionData(datastream string) string {
	splits := strings.Split(datastream, "\xc0\x80")

	str := ""
	for ix := 1; ix < len(splits); ix++ {
		str += fmt.Sprintf("[%s %s] ", splits[ix], splits[ix+1])
		ix++
	}

	return str
}
