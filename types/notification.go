package types

type Notification struct {
	UserId   uint
	To       string
	Subject  string
	Message  string
	Channels []string
}
