package forms

import "time"

const (
	AccountSid = "ACbcef2267c53f7781e4259253d31adf9f"
	AuthToken  = "bb0982c09f76a5360e3dc1bd8be62d99"
	FromNumber = "+1 567 339 2523"
)

type UserModel struct {
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Gender      string `json:"gender" bson:"gender"`
}

type LoginUserModel struct {
	PhoneNumber string     `json:"phone_number" bson:"phone_number"`
	Otp         string     `json:"otp" bson:"otp"`
	Time        *time.Time `json:"time" bson:"time"`
}
