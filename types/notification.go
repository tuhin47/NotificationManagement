package types

type Notification struct {
	To       string
	Subject  string
	Message  string
	Channels []string
}
