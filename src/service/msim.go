package service

import (
	"chimera/network"
	"chimera/network/myspace"
	"fmt"
)

func ServiceMySpaceBroadcastSignOnToRecv(recv_cli *network.Client, sender_uin int, sender_statuscode int, sender_statusmsg string) {
	if sender_statuscode != 0x0000 {
		recv_cli.Connection.WriteTraffic(myspace.MySpaceBuildPackage([]myspace.MySpaceDataPair{
			myspace.MySpaceNewDataInt("bm", 100),
			myspace.MySpaceNewDataInt("f", sender_uin),
			myspace.MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", sender_statuscode, sender_statusmsg)),
		}))
	}
}

func ServiceMySpaceBroadcastSignOnToSender(sender_cli *network.Client, recv_uin int, recv_statuscode int, recv_statusmsg string) {
	if recv_statuscode != 0x0000 {
		sender_cli.Connection.WriteTraffic(myspace.MySpaceBuildPackage([]myspace.MySpaceDataPair{
			myspace.MySpaceNewDataInt("bm", 100),
			myspace.MySpaceNewDataInt("f", recv_uin),
			myspace.MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", recv_statuscode, recv_statusmsg)),
		}))
	}
}
