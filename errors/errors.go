package errors

import "errors"

var (
	ErrMissingAuthValues        = errors.New("missing auth values")
	ErrMissingRegValues         = errors.New("missing reg values")
	ErrMissingPassResetValues   = errors.New("missing password reset values")
	ErrMissingChangePassValues  = errors.New("missing change password values")
	ErrMissingChangeSkinValues  = errors.New("missing change skin values")
	ErrMissingChangeCloakValues = errors.New("missing change cloak values")

	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidCaptcha  = errors.New("invalid captcha")
	ErrInvalidSkin     = errors.New("invalid skin")
	ErrInvalidCloak    = errors.New("invalid cloak")

	ErrUserAlreadyExists = errors.New("username or email already exists")
	ErrUserDoesNotExist  = errors.New("user does not exist")

	ErrHashingPassword   = errors.New("error hashing password")
	ErrUnhashingPassword = errors.New("error unhashing password")

	ErrWrongUsernameOrPassword = errors.New("wrong username or password")
	ErrWrongPassword           = errors.New("wrong password")

	ErrRegisteringUser      = errors.New("error registering user")
	ErrUpdatingUserPassword = errors.New("error updating user password")

	ErrSavingAuthSession       = errors.New("error saving auth session")
	ErrGettingAuthSession      = errors.New("error getting auth session")
	ErrSavingPassResetSession  = errors.New("error saving password reset session")
	ErrGettingPassResetSession = errors.New("error getting password reset session")
	ErrSavingRegSession        = errors.New("error saving reg session")
	ErrGettingRegSession       = errors.New("error getting reg session")
	ErrDeletingSession         = errors.New("error deleting session")

	ErrFailedTokenCreation         = errors.New("failed to create jwt token")
	ErrInvalidAccessToken          = errors.New("invalid access token")
	ErrExpiredAccessToken          = errors.New("access token is expired")
	ErrAccessTokenUuidNotExists    = errors.New("failed to get access token uuid")
	ErrAccessTokenUserIdNotExists  = errors.New("failed to get access token user id")
	ErrMissingRefreshToken         = errors.New("missing refresh token")
	ErrInvalidRefreshToken         = errors.New("invalid refresh token")
	ErrExpiredRefreshToken         = errors.New("refresh token is expired")
	ErrRefreshTokenUuidNotExists   = errors.New("failed to get refresh token uuid")
	ErrRefreshTokenUserIdNotExists = errors.New("failed to get refresh token user id")

	ErrSendingEmail = errors.New("error sending email")

	ErrSettingSkin   = errors.New("error setting skin")
	ErrDeletingSkin  = errors.New("error deleting skin")
	ErrSettingCloak  = errors.New("error setting cloak")
	ErrDeletingCloak = errors.New("error deleting cloak")
)
