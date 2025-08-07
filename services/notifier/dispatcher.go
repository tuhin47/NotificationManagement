package notifier

type Dispatcher struct {
	Notifiers []Notifier
}

func NewDispatcher(email *EmailNotifier, sms *SMSNotifier, telegram *TelegramNotifier) *Dispatcher {
	return &Dispatcher{
		Notifiers: []Notifier{email, sms, telegram},
	}
}

func (d *Dispatcher) NotifyAll(n Notification) error {
	for _, notifier := range d.Notifiers {
		for _, channel := range n.Channels {
			if notifier.Type() == channel {
				if err := notifier.Send(n); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
