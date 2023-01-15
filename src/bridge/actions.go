package bridge

import (
	"chimera/network"
	"chimera/network/myspace"
	"chimera/utility/database"
	"chimera/utility/logging"
)

func BridgeHandleActionNewStatusMessage(actiondata string) {

}

func BridgeHandleActionSendInstantMessage(actiondata string) {

}

func BridgeHandleActionBroadcastSignon(header BridgeMessageTransportHeader, actiondata string) {
	var Client *network.Client
	var myspace_senderctx *myspace.MySpaceContext

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN == header.SenderUIN {
			Client = network.Clients[ix]

			if Client == nil {
				logging.Error("Bridge/BridgeHandleActionBroadcastSignon", "sender offline! skipping....")
				return
			}
		}

		service := network.Clients[ix].ClientInfo.Messenger
		switch service {
		case network.MessengerMySpace:

			if myspace.ClientContexts[ix].UIN == header.SenderUIN {
				myspace_senderctx = myspace.ClientContexts[ix]

				if myspace_senderctx == nil {
					logging.Error("Bridge/BridgeHandleActionBroadcastSignon", "sender offline! skipping....")
					return
				}
			}

		case network.MessengerAIM:
		case network.MessengerMSN:
		case network.MessengerYahoo:
		default:
			logging.Warn("Bridge/BridgeHandleActionBroadcastSignon", "unknown messenger, skipping....")
			return
		}
	}

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != Client.ClientAccount.UIN { // make sure we dont fuck the network.Client up by sending it to ourselves
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", Client.ClientAccount.UIN)

			if err != nil {
				logging.Error("Bridge/BridgeHandleActionBroadcastSignon", "Failed to get contact list for uin: %d (%s)", Client.ClientAccount.UIN, err.Error())
				return
			}

			for row.Next() {
				var contact network.Contact
				err = row.Scan(&contact.SenderUIN, &contact.FriendUIN, &contact.Reason)

				if err != nil {
					logging.Error("Bridge/BridgeHandleActionBroadcastSignon", "Failed to scan contact lists (%s)", err.Error())
					row.Close()
					return
				}

				if network.Clients[ix].ClientAccount.UIN == contact.FriendUIN { // send the signon broadcast only to people on our friends list, otherwise the network.Client will add them which is bad.
					var count int
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, Client.ClientAccount.UIN)

					if err != nil {
						logging.Error("Bridge/BridgeHandleActionBroadcastSignon", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					if count > 0 {
						service := network.Clients[ix].ClientInfo.Messenger

						switch service {
						case network.MessengerMySpace:
							BridgeHandleMySpaceSignOnStatus(Client, myspace_senderctx, network.Clients[ix], myspace.ClientContexts[ix])
						case network.MessengerAIM:
						case network.MessengerMSN:
						case network.MessengerYahoo:
						default:
							logging.Warn("Bridge/BridgeHandleActionBroadcastSignon", "unknown messenger, skipping....")
							return
						}

					}
				}
			}
			row.Close()

		}
	}
}

func BridgeHandleActionBroadcastLogoff(header BridgeMessageTransportHeader, actiondata string) {
	var Client *network.Client
	var myspace_senderctx *myspace.MySpaceContext

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN == header.SenderUIN {
			Client = network.Clients[ix]

			if Client == nil {
				logging.Error("Bridge/BridgeHandleActionBroadcastLogoff", "sender offline! skipping....")
				return
			}
		}

		service := network.Clients[ix].ClientInfo.Messenger
		switch service {
		case network.MessengerMySpace:

			if myspace.ClientContexts[ix].UIN == header.SenderUIN {
				myspace_senderctx = myspace.ClientContexts[ix]

				if myspace_senderctx == nil {
					logging.Error("Bridge/BridgeHandleActionBroadcastLogoff", "sender offline! skipping....")
					return
				}
			}

		case network.MessengerAIM:
		case network.MessengerMSN:
		case network.MessengerYahoo:
		default:
			logging.Warn("Bridge/BridgeHandleActionBroadcastLogoff", "unknown messenger, skipping....")
			return
		}
	}

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != Client.ClientAccount.UIN {
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", Client.ClientAccount.UIN)

			if err != nil {
				logging.Error("Bridge/BridgeHandleActionBroadcastLogoff", "Failed to get contact list for uin: %d (%s)", Client.ClientAccount.UIN, err.Error())
				return
			}

			for row.Next() {
				var contact network.Contact
				err = row.Scan(&contact.SenderUIN, &contact.FriendUIN, &contact.Reason)

				if err != nil {
					logging.Error("Bridge/BridgeHandleActionBroadcastLogoff", "Failed to scan contact lists (%s)", err.Error())
					row.Close()
					return
				}

				if network.Clients[ix].ClientAccount.UIN == contact.FriendUIN {
					var count int
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, Client.ClientAccount.UIN)

					if err != nil {
						logging.Error("Bridge/BridgeHandleActionBroadcastLogoff", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					if count > 0 {
						service := network.Clients[ix].ClientInfo.Messenger

						switch service {
						case network.MessengerMySpace:
							BridgeHandleMySpaceLogOffStatus(myspace_senderctx, network.Clients[ix])
						case network.MessengerAIM:
						case network.MessengerMSN:
						case network.MessengerYahoo:
						default:
							logging.Warn("Bridge/BridgeHandleActionBroadcastLogoff", "unknown messenger, skipping....")
							return
						}
					}
				}
			}
			row.Close()
		}
	}
}

func BridgeHandleActionAddContact(actiondata string) {

}

func BridgeHandleActionDeleteContact(actiondata string) {

}

func BridgeHandleActionBlockUser(actiondata string) {

}

func BridgeHandleActionUnblockUser(actiondata string) {

}

func BridgeHandleActionChangeServerInformation(actiondata string) {

}

func BridgeHandleActionFetchClientDetails(actiondata string) {

}

func BridgeHandleActionFetchUserDetails(actiondata string) {

}
