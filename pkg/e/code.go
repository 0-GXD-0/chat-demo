package e

type Code int

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistInentifier    = 202
	ERROR                 = 500
	InvalidParams         = 400
	ErrorDatabase         = 40001

	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004
	WebsocketOfflineReply   = 50005
	WebsocketLimit          = 50006
)

func (c Code) Msg() string {
	return codeMsg[c]
}
