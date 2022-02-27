package systems

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidCaptcha    = errors.New("Invalid captcha")
	ErrUserDoesNotExist  = errors.New("User does not exist")
	ErrUnhashingPassword = errors.New("Error unhashing password")
	ErrUnknown           = errors.New("Unknown error")
)

var GeneralDB *sql.DB

func System(generalDB *sql.DB) {
	GeneralDB = generalDB
}
