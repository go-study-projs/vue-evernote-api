package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-study-projs/vue-evernote-api/config"
	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	prop config.Properties
)

func init() {
	if err := cleanenv.ReadEnv(&prop); err != nil {
		log.Errorf("Configuration cannot be read : %v", err)
	}
}

func Json(c echo.Context, httpStatusCode int, args ...interface{} /* msg = ""  data = nil*/) error {

	if len(args) > 2 {
		log.Errorf("the number of args is wrong : %v.", args)
		return toJson(c, http.StatusInternalServerError, "the number of args is wrong", nil)
	}
	r := &model.Response{Msg: "", Data: nil}
	for i, v := range args {
		switch i {
		case 0: //msg
			msg, ok := v.(string)
			if !ok {
				log.Errorf("%v is not passed as string", v)
				return toJson(c, http.StatusInternalServerError, "msg is not passed as string", nil)
			}
			r.Msg = msg
		case 1: //data
			r.Data = v
		default:
			log.Errorf("unknown argument passed : %v.", v)
			return toJson(c, http.StatusInternalServerError, "unknown argument passed", nil)
		}
	}
	return toJson(c, httpStatusCode, r.Msg, r.Data)
}

func toJson(c echo.Context, httpStatusCode int, msg string, data interface{}) error {
	return c.JSON(httpStatusCode, model.Response{Msg: msg, Data: data})
}

func CreateToken(u model.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = u.ID.String()
	claims["user_name"] = u.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(prop.JwtTokenSecret))
	if err != nil {
		log.Errorf("Unable to generate the token :%v", err)
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStringWithBearer string) (userId primitive.ObjectID) {
	tokenString := strings.Split(tokenStringWithBearer, " ")[1]
	return parsePureToken(tokenString)
}

func parsePureToken(tokenString string) (userId primitive.ObjectID) {
	claims, err := parseToken(tokenString)
	if err != nil {
		log.Errorf("token is invalid :%s", tokenString)
		return primitive.NilObjectID
	}

	id, ok := claims["user_id"].(string)
	if !ok {
		log.Errorf("%v is not passed as string", id)
		return primitive.NilObjectID
	}
	// ObjectID("61fa1e3b67215eb9c7418145") => 61fa1e3b67215eb9c7418145
	id = id[10:34]
	uid, _ := primitive.ObjectIDFromHex(id)
	return uid

}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(prop.JwtTokenSecret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
