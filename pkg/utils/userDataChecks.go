package utils

import (
    "strconv"
    "strings"
    "time"
    "unicode"

    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

func CheckUsername(s string) error {
    if len(s) < 3 || len(s) > 30 {
        return errors.ErrWrongUsernameLength
    }

    allowedChars := "_-."
    for _, c := range []rune(s) {
        if !unicode.IsDigit(c) &&
            !strings.ContainsRune(allowedChars, c) &&
            !IsLatin(c) {
            return errors.ErrWrongUsername
        }
    }

    return nil
}

func CheckPassword(s string) error {
    if len(s) < 8 {
        return errors.ErrShortPass
    }

    var contains_num, contains_letter, contains_sym bool

    allowedChars := "!@#$%^&*()-_=+{}[]\\|;:,<.>/?"
    for _, c := range []rune(s) {
        if unicode.IsDigit(c) {
            contains_num = true
        } else if strings.ContainsRune(allowedChars, c) {
            contains_sym = true
        } else if IsLatin(c) {
            contains_letter = true
        } else {
            // Сюда попадаем, когда появился запрещенный символ в пароле
            return errors.ErrSimplePass
        }
    }

    if !(contains_num && contains_letter && contains_sym) {
        // Сюда попадаем, если не соблюдены условия сложности пароля
        return errors.ErrSimplePass
    }

    return nil
}

func CheckName(s string) error {
    if len(s) > 60 {
        return errors.ErrWrongNameLength
    }

    for _, c := range []rune(s) {
        if !unicode.IsLetter(c) && !(c == ' ') {
            return errors.ErrWrongName
        }
    }

    return nil
}

func ParseDate(s string) (time.Time, error) {
    dateParts := strings.Split(s, "-")
    if len(dateParts) != 3 || len(dateParts[0]) != 4 || len(dateParts[1]) != 2 || len(dateParts[2]) != 2 {
        return time.Time{}, errors.ErrWrongDateFormat
    }

    year, err := strconv.Atoi(dateParts[0])
    if err != nil {
        return time.Time{}, errors.ErrWrongDateFormat
    }

    month, err := strconv.Atoi(dateParts[1])
    if err != nil {
        return time.Time{}, errors.ErrWrongDateFormat
    }

    day, err := strconv.Atoi(dateParts[2])
    if err != nil {
        return time.Time{}, errors.ErrWrongDateFormat
    }

    date := time.Date(
        year,
        time.Month(month),
        day,
        0, 0, 0, 0, nil)

    return date, nil
}

func IsLatin(c rune) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
