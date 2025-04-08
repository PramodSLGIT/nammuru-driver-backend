package jwtauth

import (
	"encoding/json"
	"nammuru-driver-backend/forms"
	"nammuru-driver-backend/models"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const PhoneNumber = "phone_number"

func InitJwt() *jwt.GinJWTMiddleware {
	var authMiddleware *jwt.GinJWTMiddleware
	var identityKey = "id"
	authMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("pramod"),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			var v = data.(forms.UserRegistrationForm)

			return jwt.MapClaims{
				identityKey: v.PhoneNumber,
				"data":      v,
			}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var user forms.UserRegistrationForm

			jsonString, _ := json.Marshal(claims["data"])
			json.Unmarshal(jsonString, &user)

			return &forms.UserRegistrationForm{
				PhoneNumber: claims[identityKey].(string),
			}

		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login forms.LoginModel
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			phoneNumber := login.PhoneNumber
			password := login.Password

			var loginModel = new(models.UserRegistration)
			user, _, err := loginModel.UserLogin(phoneNumber, password)

			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			c.Set("data", user)
			return user, nil
		},

		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.Header("token", token)
			c.Header("expire", expire.Format(time.RFC3339))
			data, exist := c.Get("data")
			if exist {
				c.JSON(http.StatusOK, gin.H{"data": data})
			}
		},

		TokenHeadName: "Bearer",
		TokenLookup:   "header:Authorization, query:token, cookie:jwt",
		TimeFunc:      time.Now,
	})

	return authMiddleware
}
