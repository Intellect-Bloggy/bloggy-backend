package errors

import (
    "errors"
    "net/http"
)

var (
    ErrTokenNotFound    = NewHTTPError(http.StatusForbidden, errors.New("token does not exist"))
    ErrTokenExpired     = NewHTTPError(http.StatusForbidden, errors.New("token expired"))
    ErrWrongCredentials = NewHTTPError(http.StatusForbidden, errors.New("wrong username or password"))

    ErrShortPass  = NewHTTPError(http.StatusBadRequest, errors.New("password must be at least 8 characters"))
    ErrSimplePass = NewHTTPError(http.StatusBadRequest, errors.New("password must contain a mix of letters, numbers, and symbols"))
)
