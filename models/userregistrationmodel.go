package models

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"
	"nammuru-driver-backend/database"
	"nammuru-driver-backend/forms"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mdbConn, _ = database.NewDbConn()

type UserRegistration struct {
}

func (u *UserRegistration) Register(data forms.UserModel) (resultId interface{}, err error) {
	ctx := context.TODO()
	collection := mdbConn.Use(forms.DB, forms.DriverCollection)

	filter := bson.M{"phone_number": data.PhoneNumber}
	update := bson.M{"$set": data}
	opts := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result.UpsertedID, err
}

func (u *UserRegistration) Login(phoneNumber string, otp string) (userData forms.UserModel, err error) {
	ctx := context.TODO()
	collection := mdbConn.Use(forms.DB, forms.DriverCollection)

	filter := bson.M{"phone_number": phoneNumber}

	var result forms.LoginUserModel
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return userData, err
	}

	currentTime := time.Now()
	duration := currentTime.Sub(*result.Time)

	if duration < 3*time.Minute {
		if result.Otp == otp {
			err1 := collection.FindOne(ctx, filter).Decode(&userData)
			if err1 != nil {
				log.Println(err1)
			}
			return userData, err
		} else {
			err = errors.New("Otp is mismatch")
			return userData, err
		}
	}
	err = errors.New("otp has been expired")
	return userData, err
}

func (u *UserRegistration) GenerateOtp(data string) (done string, err error) {
	ctx := context.TODO()
	collection := mdbConn.Use(forms.DB, forms.DriverCollection)

	var user forms.UserModel
	err = collection.FindOne(ctx, bson.M{"phone_number": data}).Decode(&user)
	if err == nil {
		otp, err := u.generate()
		if err != nil {
			return done, err
		}
		err = u.sendSMS(data, otp)
		if err != nil {
			return done, err
		}
		fields := bson.M{
			"phone_number": user.PhoneNumber,
			"otp":          otp,
			"time":         time.Now(),
		}
		update := bson.M{"$set": fields}
		opts := options.Update().SetUpsert(true)
		_, err = collection.UpdateOne(ctx, bson.M{"phone_number": data}, update, opts)
		fmt.Println("success")
	} else {
		err = errors.New("Not a registered Mobile Number. Please use Registered Mobile Number")
	}

	return done, err
}

func (u *UserRegistration) generate() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}

func (u *UserRegistration) sendSMS(phoneNumber, otp string) error {

	fromNumber := forms.FromNumber
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: forms.AccountSid,
		Password: forms.AuthToken,
	})

	// Create the message with the desired format
	messageBody := fmt.Sprintf("Your OTP is %s. It is valid for 5 minutes. - From Pram", otp)

	// Send SMS
	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(fromNumber)
	params.SetBody(messageBody)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	log.Printf("Message sent! SID: %s", *resp.Sid)
	return nil
}

func (u *UserRegistration) ExampleGenerateOtp(data string) (done string, err error) {
	ctx := context.TODO()
	collection := mdbConn.Use(forms.DB, forms.DriverCollection)

	var user forms.UserModel
	err = collection.FindOne(ctx, bson.M{"phone_number": data}).Decode(&user)
	if err == nil {
		otp, err := u.generate()
		if err != nil {
			return done, err
		}

		fields := bson.M{
			"phone_number": user.PhoneNumber,
			"otp":          otp,
			"time":         time.Now(),
		}
		update := bson.M{"$set": fields}
		opts := options.Update().SetUpsert(true)
		_, err = collection.UpdateOne(ctx, bson.M{"phone_number": data}, update, opts)
		fmt.Println("success")
	} else {
		err = errors.New("Not a registered Mobile Number. Please use Registered Mobile Number")
	}

	return done, err
}

func (u *UserRegistration) AddProfileImage(data forms.UserModel) (updateId interface{}, err error) {
	ctx := context.TODO()
	collection := mdbConn.Use(forms.DB, forms.DriverCollection)

	filter := bson.M{
		"phone_number": data.PhoneNumber,
	}

	update := bson.M{
		"$set": bson.M{
			"profile_image": data.ProfileImage,
		},
	}

	updateId, err = collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return updateId, err
	}
	return updateId, err
}

func (u *UserRegistration) UserLogin(phoneNumber string, password string) (user forms.UserRegistrationForm, success bool, err error) {
	ctx := context.TODO()
	collection := mdbConn.Use(forms.DB, forms.DriverCollection)
	var userm forms.UserModel
	filter := bson.M{
		"phone_number": phoneNumber,
	}
	err = collection.FindOne(ctx, filter).Decode(&userm)
	if err != nil {
		return user, false, err
	}

	// Password check (should use bcrypt ideally)
	if userm.PhoneNumber != phoneNumber || userm.Password != password {
		return user, false, errors.New("invalid credentials")
	}

	base64Image := ""
	if len(userm.ProfileImage) > 0 {
		base64Image = base64.StdEncoding.EncodeToString(userm.ProfileImage)
	}

	// Map UserModel to UserRegistrationForm
	user = forms.UserRegistrationForm{
		ProfileImage: base64Image,
		Name:         userm.Name,
		Email:        userm.Email,
		PhoneNumber:  userm.PhoneNumber,
		Gender:       userm.Gender,
		KYCData:      userm.KYCData,
		Password:     userm.Password, // Consider omitting or masking
	}
	return user, true, err
}
