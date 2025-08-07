package conn

import (
	"NotificationManagement/config"
	"net/http"
	"time"
)

var emailClient *http.Client

func ConnectEmail() {
	config := config.Email()
	timeout := config.Timeout * time.Second
	emailClient = newHTTPClient(timeout, 50)
}

func EmailClient() *http.Client {
	return emailClient
}
