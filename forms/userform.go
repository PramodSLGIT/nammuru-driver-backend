package forms

import "time"

const (
	AccountSid       = "ACbcef2267c53f7781e4259253d31adf9f"
	AuthToken        = "bb0982c09f76a5360e3dc1bd8be62d99"
	FromNumber       = "+1 567 339 2523"
	DB               = "driver"
	DriverCollection = "driver_info"
	RideCollection   = "ride_info"
)

type UserModel struct {
	Name         string   `json:"name" bson:"name"`
	Email        string   `json:"email" bson:"email"`
	PhoneNumber  string   `json:"phone_number" bson:"phone_number"`
	Gender       string   `json:"gender" bson:"gender"`
	KYCData      KYC      `json:"kyc" bson:"kyc"`
	VechicleData Vechicle `json:"vechicle" bson:"vehicle"`
}

type LoginUserModel struct {
	PhoneNumber string     `json:"phone_number" bson:"phone_number"`
	Otp         string     `json:"otp" bson:"otp"`
	Time        *time.Time `json:"time" bson:"time"`
}

type KYC struct {
	Aadhaar  string `json:"aadhaar" bson:"aadhaar"`
	DLNumber string `json:"dl_number" bson"dl_number"`
}

type Vechicle struct {
	Registration string `json:"registration" bson:"registration"`
	Model        string `json:"model" bson:"model"`
}
