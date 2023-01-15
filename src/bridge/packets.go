package bridge

import (
	"chimera/network"
	"chimera/network/myspace"
	"fmt"
)

/* MySpace */

func BridgeHandleMySpaceSignOnStatus(sender_client *network.Client, sender_ctx *myspace.MySpaceContext, reciever_client *network.Client, reciever_ctx *myspace.MySpaceContext) {
	reciever_client.Connection.WriteTraffic(myspace.MySpaceBuildPackage([]myspace.MySpaceDataPair{
		myspace.MySpaceNewDataInt("bm", 100),
		myspace.MySpaceNewDataInt("f", sender_client.ClientAccount.UIN),
		myspace.MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", sender_ctx.Status.Code, sender_ctx.Status.Message)),
	}))
	sender_client.Connection.WriteTraffic(myspace.MySpaceBuildPackage([]myspace.MySpaceDataPair{
		myspace.MySpaceNewDataInt("bm", 100),
		myspace.MySpaceNewDataInt("f", reciever_client.ClientAccount.UIN),
		myspace.MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|%d|ss|%s", reciever_ctx.Status.Code, reciever_ctx.Status.Message)),
	}))
}

func BridgeHandleMySpaceLogOffStatus(sender_ctx *myspace.MySpaceContext, reciever_client *network.Client) {
	reciever_client.Connection.WriteTraffic(myspace.MySpaceBuildPackage([]myspace.MySpaceDataPair{
		myspace.MySpaceNewDataInt("bm", 100),
		myspace.MySpaceNewDataInt("f", sender_ctx.UIN),
		myspace.MySpaceNewDataGeneric("msg", fmt.Sprintf("|s|0|ss|%s", sender_ctx.Status.Message)),
	}))
}

/* OSCAR */
