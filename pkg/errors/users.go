package errors

import (
    "errors"
    "net/http"
)

var (
    ErrUserNotFound     = NewHTTPError(http.StatusNotFound, errors.New("user does not exits"))
    ErrUsernameNotFound = NewHTTPError(http.StatusNotFound, errors.New("user with this username does not exist"))
    ErrIdNotFound       = NewHTTPError(http.StatusNotFound, errors.New("id does not exist"))

    ErrWrongUsername       = NewHTTPError(http.StatusBadRequest, errors.New("username must contain only letters, numbers and some special characters"))
    ErrWrongUsernameLength = NewHTTPError(http.StatusBadRequest, errors.New("username must be 3 to 30 characters long"))
    ErrWrongName           = NewHTTPError(http.StatusBadRequest, errors.New("name must contain only letters and spaces"))
    ErrWrongNameLength     = NewHTTPError(http.StatusBadRequest, errors.New("name must be no longer than 60 characters"))
    ErrWrongId             = NewHTTPError(http.StatusBadRequest, errors.New("id must be non negative integer"))
    ErrWrongDateFormat     = NewHTTPError(http.StatusBadRequest, errors.New(`date must be "2000-12-31" format`))

    ErrTakenUsername = NewHTTPError(http.StatusConflict, errors.New("user with this username already exists"))
    ErrTakenEmail    = NewHTTPError(http.StatusConflict, errors.New("user with this email already exists"))
    ErrTakenPhone    = NewHTTPError(http.StatusConflict, errors.New("user with this phone already exists"))
)
