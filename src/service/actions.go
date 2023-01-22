package service

import (
	"chimera/network"
	"chimera/utility/database"
	"chimera/utility/logging"
)

func ServiceActionBroadcastSignOnStatus(msg *network.ServiceMessage) {
	service := ServiceTranslateGetInterOp(msg)

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != service.Client.ClientAccount.UIN { // make sure we dont fuck the client up by sending it to ourselves
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", service.Client.ClientAccount.UIN)

			if err != nil {
				logging.Error("Service/ServiceActionBroadcastSignOnStatus", "Failed to get contact list for uin: %d (%s)", service.Client.ClientAccount.UIN, err.Error())
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
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, service.Client.ClientAccount.UIN)

					if err != nil {
						logging.Error("Service/ServiceActionBroadcastSignOnStatus", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					sender, err := database.GetUserDetailsDataByUIN(service.Client.ClientAccount.UIN)

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
							ServiceMySpaceBroadcastSignOnToRecv(network.Clients[ix], service.Client.ClientAccount.UIN, status, message)
						case network.Service_OSCAR:
							//ServiceOscarBroadcastSignOn(network.Clients[ix]) // please actually implement this fruther
						default:
							logging.Warn("Service/ServiceActionBroadcastSignOnStatus", "unknown messenger, skipping....")
							return
						}

						switch msg.Service {
						case network.Service_MSIM:
							status, _ := ServiceTranslateToMsimStatus(recv.StatusCode, recv.StatusMessage)
							ServiceMySpaceBroadcastSignOnToSender(service.Client, network.Clients[ix].ClientAccount.UIN, status, recv.StatusMessage)
						case network.Service_OSCAR:
							// please actually implement this fruther
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

func ServiceActionBroadcastLogOffStatus(msg *network.ServiceMessage) {
	service := ServiceTranslateGetInterOp(msg)

	for ix := 0; ix < len(network.Clients); ix++ {
		if network.Clients[ix].ClientAccount.UIN != service.Client.ClientAccount.UIN {
			row, err := database.Query("SELECT * from contacts WHERE SenderUIN= ?", service.Client.ClientAccount.UIN)

			if err != nil {
				logging.Error("Service/ServiceActionBroadcastLogOffStatus", "Failed to get contact list for uin: %d (%s)", service.Client.ClientAccount.UIN, err.Error())
				return
			}

			for row.Next() {
				var contact network.Contact
				err = row.Scan(&contact.SenderUIN, &contact.FriendUIN, &contact.Reason)

				if err != nil {
					logging.Error("Service/ServiceActionBroadcastLogOffStatus", "Failed to scan contact lists (%s)", err.Error())
					row.Close()
					return
				}

				if network.Clients[ix].ClientAccount.UIN == contact.FriendUIN {
					var count int
					innerrow, err := database.Query("SELECT COUNT(*) from contacts WHERE SenderUIN= ? AND RecvUIN= ?", network.Clients[ix].ClientAccount.UIN, service.Client.ClientAccount.UIN)

					if err != nil {
						logging.Error("Service/ServiceActionBroadcastLogOffStatus", "Failed to count contact list shit (%s)", err.Error())
						row.Close()
						return
					}

					innerrow.Next()
					innerrow.Scan(&count)
					innerrow.Close()

					if count > 0 {
						switch network.Clients[ix].ClientInfo.Service {
						case network.Service_MSIM:
							_, message := ServiceTranslateToMsimStatus(StatusCode_Offline, service.MSIM_Context.Status.Message)
							ServiceMySpaceBroadcastLogOff(network.Clients[ix], service.Client.ClientAccount.UIN, message)
						case network.Service_OSCAR:
							// please actually implement this fruther
						default:
							logging.Warn("Service/ServiceActionBroadcastLogOffStatus", "unknown messenger, skipping....")
							return
						}

					}
				}
			}
			row.Close()
		}
	}
}

func ServiceActionDeliverOfflineIM(msg *network.ServiceMessage) {
	service := ServiceTranslateGetInterOp(msg)

	//(17, 15, 1669828262515, '<p><f f=\'Times\' h=\'16\'><c v=\'black\'><b v=\'white\'>test</1b></1c></1f></1p>'),

	row, err := database.Query("SELECT * from offlinemsgs WHERE RecvUIN= ?", service.Client.ClientAccount.UIN)

	if err != nil {
		logging.Error("Service/ServiceActionDeliverOfflineIM", "Failed to get offline messages list for uin: %d (%s)", service.Client.ClientAccount.UIN, err.Error())
		return
	}

	for row.Next() {
		var message network.OfflineMessage
		var discardedMessageAttributes string
		err = row.Scan(&message.SenderUIN, &message.RecvUIN, &message.MessageDate, &message.MessageContent, &discardedMessageAttributes)

		if err != nil {
			logging.Error("Service/ServiceActionDeliverOfflineIM", "Failed to scan offline messages list of uin: %d (%s)", service.Client.ClientAccount.UIN, err.Error())
			row.Close()
			return
		}

		/* for now i will just pass this through as is, please implment a tokenizer for MessageAttributes*/
		switch msg.Service {
		case network.Service_MSIM:
			ServiceMySpaceDeliverOfflineIM(service.Client, service.MSIM_Context, message)
		case network.Service_OSCAR:
			// please actually implement this fruther
		default:
			logging.Warn("Service/ServiceActionDeliverOfflineIM", "unknown messenger, skipping....")
			return
		}

	}
	row.Close()

	row, err = database.Query("DELETE from offlinemsgs WHERE RecvUIN= ?", service.Client.ClientAccount.UIN)

	if err != nil {
		logging.Error("Service/ServiceActionDeliverOfflineIM", "Failed to delete offline messages for uin: %d (%s)", service.Client.ClientAccount.UIN, err.Error())
		return
	}

	row.Close()
}
