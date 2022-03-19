package errors

import "errors"

var (
	ErrMissingAuthValues          = errors.New("missing auth values")
	ErrMissingRegValues           = errors.New("missing reg values")
	ErrMissingPassResetValues     = errors.New("missing password reset values")
	ErrMissingRefreshToken        = errors.New("missing refresh token")
	ErrMissingChangePassValues    = errors.New("missing change password values")
	ErrMissingChangeSkinValues    = errors.New("missing change skin values")
	ErrMissingDeleteSkinValues    = errors.New("missing delete skin values")
	ErrMissingChangeCloakValues   = errors.New("missing change cloak values")
	ErrMissingDeleteCloakValues   = errors.New("missing delete cloak values")
	ErrMissingTransferMoneyValues = errors.New("missing transfer money values")

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
	ErrInvalidRefreshToken         = errors.New("invalid refresh token")
	ErrExpiredRefreshToken         = errors.New("refresh token is expired")
	ErrRefreshTokenUuidNotExists   = errors.New("failed to get refresh token uuid")
	ErrRefreshTokenUserIdNotExists = errors.New("failed to get refresh token user id")

	ErrNotEnoughMoney     = errors.New("error not enough money")
	ErrMoneySubtraction   = errors.New("error money subtraction")
	ErrMoneyAddition      = errors.New("error money addition")
	ErrSavingTransferInfo = errors.New("error saving transfer info")

	ErrDBQuery                = errors.New("error database query")
	ErrDBRegisteringUser      = errors.New("error adding user to database")
	ErrDBUpdatingUserPassword = errors.New("error updating password in database")

	ErrUploadingSkin  = errors.New("error uploading skin")
	ErrDeletingSkin   = errors.New("error deleting skin")
	ErrUploadingCloak = errors.New("error uploading cloak")
	ErrDeletingCloak  = errors.New("error deleting cloak")

	ErrSendingEmail = errors.New("error sending email")
)
