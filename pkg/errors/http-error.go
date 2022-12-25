package errors

import (
    "net/http"
)

// EtoHe конвертирует неизвестный тип error в HTTPError.
// При ошибке конвертации создает новый объект HTTPError с кодом 500.
func EtoHe(err error) HTTPError {
    httpErr, ok := err.(HTTPError)
    if ok {
        return httpErr
    }

    return NewHTTPError(http.StatusInternalServerError, err)
}

func NewHTTPError(code int, err error) HTTPError {
    return &httpError{
        error: err,
        code:  code,
    }
}

type HTTPError interface {
    Error() string
    Status() int
}

type httpError struct {
    error
    code int
}

// Status Возвращает HTTP статус ответа
func (e httpError) Status() int {
    return e.code
}
