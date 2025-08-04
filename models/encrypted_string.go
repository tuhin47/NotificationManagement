package models

import (
	"NotificationManagement/config"
	"NotificationManagement/security"
	"database/sql/driver"
)

type EncryptedString string

func (es *EncryptedString) Scan(value interface{}) error {
	if val, ok := value.(string); ok {
		plain, err := security.DecryptAES(val, []byte(config.App().Encryption))
		if err != nil {
			return err
		}
		*es = EncryptedString(plain)
	}
	return nil
}

func (es EncryptedString) Value() (driver.Value, error) {
	return security.EncryptAES([]byte(es), []byte(config.App().Encryption))

}
