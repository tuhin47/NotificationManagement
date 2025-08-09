package types

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type VerifyOtpRequest struct {
	OTP string `json:"otp"`
}

func (r *VerifyOtpRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.OTP, validation.Required, validation.Length(6, 6)),
	)
}
