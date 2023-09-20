package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/constant"
)

func PanicException_(key int, message string) {
	err := errors.New(message)
	err = fmt.Errorf("%d: %w", key, err)
	if err != nil {
		panic(err)
	}
}

func PanicException(responseKey constant.ResponseStatus) {
	PanicException_(responseKey.GetResponseStatus(), responseKey.GetResponseMessage())
}

func PanicExceptionWithMessage(responseKey constant.ResponseStatus, message string) {
	PanicException_(responseKey.GetResponseStatus(), message)
}

func PanicHandler(c *gin.Context) {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		strArr := strings.Split(str, ":")

		keyStr := strArr[0]
		key, _ := strconv.Atoi(keyStr)
		msg := strings.Trim(strArr[1], " ")

		switch key {
		case constant.InvalidRequest.GetResponseStatus():
			c.JSON(http.StatusBadRequest, BuildResponse_(key, msg, Null()))
		case
			constant.DataNotFound.GetResponseStatus():
			c.JSON(http.StatusNotFound, BuildResponse_(key, msg, Null()))
			c.Abort()
		case
			constant.Unauthorized.GetResponseStatus():
			c.JSON(http.StatusUnauthorized, BuildResponse_(key, msg, Null()))
			c.Abort()
		default:
			c.JSON(http.StatusInternalServerError, BuildResponse_(key, msg, Null()))
			c.Abort()
		}
	}
}
