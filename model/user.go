package model

import (
	"time"

	"github.com/go-study-projs/vue-evernote-api/config"
	"github.com/golang-jwt/jwt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/gommon/log"
)

type User struct {
	Username string `json:"username" bson:"username" validate:"required,min=1,max=15"`
	Password string `json:"password,omitempty" bson:"password" validate:"required,min=6,max=16"`
}

var (
	prop config.Properties
)

func (u User) CreateToken() (string, error) {
	if err := cleanenv.ReadEnv(&prop); err != nil {
		log.Errorf("Configuration cannot be read : %v", err)
	}
	claims := jwt.MapClaims{}
	//claims["authorized"] = u.IsAdmin
	claims["user_name"] = u.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(prop.JwtTokenSecret))
	if err != nil {
		log.Errorf("Unable to generate the token :%v", err)
		return "", err
	}
	return token, nil
}
