package notifier

type Notification struct {
	To       string
	Subject  string
	Message  string
	Channels []string
}

type Notifier interface {
	Send(notification Notification) error
	Type() string
}
