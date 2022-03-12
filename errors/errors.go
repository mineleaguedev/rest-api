package errors

import "errors"

var (
	ErrMissingAuthValues       = errors.New("missing auth errors")
	ErrMissingRegValues        = errors.New("missing reg errors")
	ErrInvalidUsername         = errors.New("invalid username")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrInvalidCaptcha          = errors.New("invalid captcha")
	ErrUserAlreadyExists       = errors.New("username or email already exists")
	ErrUserDoesNotExist        = errors.New("user does not exist")
	ErrHashingPassword         = errors.New("error hashing password")
	ErrUnhashingPassword       = errors.New("error unhashing password")
	ErrWrongUsernameOrPassword = errors.New("wrong username or password")
	ErrRegUser                 = errors.New("error registering user")

	ErrFailedTokenCreation     = errors.New("failed to create jwt token")
	ErrSavingAuthSession       = errors.New("error saving auth session")
	ErrGettingAuthSession      = errors.New("error getting auth session")
	ErrSavingPassResetSession  = errors.New("error saving password reset session")
	ErrGettingPassResetSession = errors.New("error getting password reset session")
	ErrSavingRegSession        = errors.New("error saving reg session")
	ErrGettingRegSession       = errors.New("error getting reg session")
	ErrDeletingSession         = errors.New("error deleting session")

	ErrSendingEmail = errors.New("error sending email")

	ErrInvalidAccessToken          = errors.New("invalid access token")
	ErrExpiredAccessToken          = errors.New("access token is expired")
	ErrAccessTokenUuidNotExists    = errors.New("failed to get access token uuid")
	ErrAccessTokenUserIdNotExists  = errors.New("failed to get access token user id")
	ErrMissingRefreshToken         = errors.New("missing refresh token")
	ErrInvalidRefreshToken         = errors.New("invalid refresh token")
	ErrExpiredRefreshToken         = errors.New("refresh token is expired")
	ErrRefreshTokenUuidNotExists   = errors.New("failed to get refresh token uuid")
	ErrRefreshTokenUserIdNotExists = errors.New("failed to get refresh token user id")
)
