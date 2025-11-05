package handler

import (
	"milestone2/internal/entity"
	"net/http"

	"github.com/sirupsen/logrus"
)

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	
	switch err {
	case entity.ErrInternalServerError:
		return http.StatusInternalServerError
	case entity.ErrBadParamInput:
		return http.StatusBadRequest
	case entity.ErrConflict:
		return http.StatusConflict
	case entity.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}