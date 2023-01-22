package service

import (
	"chimera/network"
	"chimera/network/myspace"
	"chimera/network/oscar"
	"chimera/utility/database"
	"chimera/utility/logging"
)

func ServiceActionBroadcastSignOnStatus(msg *network.ServiceMessage) {
	var cli *network.Client
	var myspace_senderctx *myspace.MySpaceContext
	var oscar_ctx *oscar.OSCARContext

	for ix := 0; ix < len(network.Clients); ix++ {
		cli = network.Clients[ix]

		if cli == nil {
			logging.Error("Service/ServiceActionBroadcastSignOnStatus", "sender offline! skipping....")
			return
		}

		switch msg.Service {
		case network.Service_MSIM:
			if myspace.ClientContexts[ix].UIN == msg.Data.Sender {
				myspace_senderctx = myspace.ClientContexts[ix]

				if myspace_senderctx == nil {
					logging.Error("Service/ServiceActionBroadcastSignOnStatus", "sender contextless! skipping....")
					return
				}
			}

		case network.Service_OSCAR:
			if oscar.ClientContexts[ix].UIN == msg.Data.Sender {
				oscar_ctx = oscar.ClientContexts[ix]

				if oscar_ctx == nil {
					logging.Error("Service/ServiceActionBroadcastSignOnStatus", "sender contextless! skipping....")
					return
				}
			}

		default:
			logging.Warn("Service/ServiceActionBroadcastSignOnStatus", "unknown messenger, skipping....")
			return
		}
	}

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != cli.ClientAccount.UIN { // make sure we dont fuck the client up by sending it to ourselves
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", cli.ClientAccount.UIN)

			if err != nil {
				logging.Error("Service/ServiceActionBroadcastSignOnStatus", "Failed to get contact list for uin: %d (%s)", cli.ClientAccount.UIN, err.Error())
				return
			}

			for row.Next() {
				var contact network.Contact
				err = row.Scan(&contact.SenderUIN, &contact.FriendUIN, &contact.Reason)

				if err != nil {
					logging.Error("Service/ServiceActionBroadcastSignOnStatus", "Failed to scan contact lists (%s)", err.Error())
					row.Close()
					return
				}

				if network.Clients[ix].ClientAccount.UIN == contact.FriendUIN { // send the signon broadcast only to people on our friends list, otherwise the client will add them which is bad.
					var count int
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, cli.ClientAccount.UIN)

					if err != nil {
						logging.Error("Service/ServiceActionBroadcastSignOnStatus", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					sender, err := database.GetUserDetailsDataByUIN(cli.ClientAccount.UIN)

					if err != nil {
						logging.Error("Service/ServiceActionBroadcastSignOnStatus", "Unable to get Sender UserDetails! (%s)", err.Error())
						return
					}

					recv, err := database.GetUserDetailsDataByUIN(network.Clients[ix].ClientAccount.UIN)

					if err != nil {
						logging.Error("Service/ServiceActionBroadcastSignOnStatus", "Unable to get Recv UserDetails! (%s)", err.Error())
						return
					}

					if count > 0 {
						switch network.Clients[ix].ClientInfo.Service {
						case network.Service_MSIM:
							status, message := ServiceTranslateToMsimStatus(sender.StatusCode, sender.StatusMessage)
							ServiceMySpaceBroadcastSignOnToRecv(network.Clients[ix], cli.ClientAccount.UIN, status, message)
						case network.Service_OSCAR:
							//ServiceOscarBroadcastSignOn(network.Clients[ix]) // please actually implement this fruther
						default:
							logging.Warn("Service/ServiceActionBroadcastSignOnStatus", "unknown messenger, skipping....")
							return
						}

						switch msg.Service {
						case network.Service_MSIM:
							status, _ := ServiceTranslateToMsimStatus(recv.StatusCode, recv.StatusMessage)
							ServiceMySpaceBroadcastSignOnToSender(cli, network.Clients[ix].ClientAccount.UIN, status, recv.StatusMessage)
						case network.Service_OSCAR:
							//ServiceOscarBroadcastSignOn(network.Clients[ix]) // please actually implement this fruther
						default:
							logging.Warn("Service/ServiceActionBroadcastSignOnStatus", "unknown messenger, skipping....")
							return
						}

					}
				}
			}
			row.Close()

		}
	}
}
