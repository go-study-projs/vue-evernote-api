package utils

import (
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/go-study-projs/vue-evernote-api/model"
	"github.com/labstack/echo/v4"
)

//type record struct {
//	httpStatusCode int
//	msg            string
//	data           interface{}
//}

func Json(c echo.Context, httpStatusCode int, args ...interface{} /* msg = ""  data = nil*/) error {

	if len(args) > 2 {
		log.Errorf("the number of args is wrong : %v.", args)
		return toJson(c, http.StatusInternalServerError, "the number of args is wrong", nil)
	}
	r := &model.Response{
		Msg:  "",
		Data: nil,
	}
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
