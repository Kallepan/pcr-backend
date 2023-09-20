package pkg

import (
	"gitlab.com/kallepan/pcr-backend/app/constant"
	"gitlab.com/kallepan/pcr-backend/app/domain/dto"
)

func Null() interface{} {
	return nil
}

func BuildResponseWithMessage[T any](responseStatus constant.ResponseStatus, data T, message string) dto.APIResponse[T] {
	return BuildResponse_(responseStatus.GetResponseStatus(), message, data)
}

func BuildResponse[T any](responseStatus constant.ResponseStatus, data T) dto.APIResponse[T] {
	return BuildResponse_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data)
}

func BuildResponse_[T any](status int, message string, data T) dto.APIResponse[T] {
	return dto.APIResponse[T]{
		ResponseKey:     status,
		ResponseMessage: message,
		Data:            data,
	}
}
