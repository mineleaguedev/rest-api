package errors

import "errors"

var (
	ErrDBGettingLastInsertId = errors.New("error getting last insert id from insert database query")
	ErrDBGettingRowsAffected = errors.New("error getting rows affected from update database query")
)

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

	ErrUploadingSkin  = errors.New("error uploading skin")
	ErrDeletingSkin   = errors.New("error deleting skin")
	ErrUploadingCloak = errors.New("error uploading cloak")
	ErrDeletingCloak  = errors.New("error deleting cloak")

	ErrNotEnoughMoney     = errors.New("error not enough money")
	ErrMoneySubtraction   = errors.New("error money subtraction")
	ErrMoneyAddition      = errors.New("error money addition")
	ErrSavingTransferInfo = errors.New("error saving transfer info")

	ErrDBRegisteringUser      = errors.New("error inserting user to general database")
	ErrDBGettingUser          = errors.New("error getting user from general database")
	ErrDBUpdatingUserPassword = errors.New("error updating password in general database")

	ErrSendingEmail = errors.New("error sending email")
)

var (
	ErrMissingPlayerCreateValues         = errors.New("missing player create values")
	ErrMissingPlayerUpdateExpValues      = errors.New("missing player update exp values")
	ErrMissingPlayerUpdateRankValues     = errors.New("missing player update rank values")
	ErrMissingPlayerUpdateCoinsValues    = errors.New("missing player update coins values")
	ErrMissingPlayerUpdatePlaytimeValues = errors.New("missing player update playtime values")
	ErrMissingPlayerUpdateLastSeenValues = errors.New("missing player update last seen values")
	ErrMissingPlayerBanValues            = errors.New("missing player ban values")
	ErrMissingPlayerUnbanValues          = errors.New("missing player unban values")
	ErrMissingPlayerMuteValues           = errors.New("missing player mute values")
	ErrMissingPlayerUnmuteValues         = errors.New("missing player unmute values")

	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrPlayerDoesNotExist  = errors.New("player does not exist")

	ErrPlayerIsNotBanned = errors.New("player is not banned")
	ErrPlayerIsNotMuted  = errors.New("player is not muted")

	ErrMiniGamesDBCreatingPlayer         = errors.New("error inserting player to minigames database")
	ErrMiniGamesDBGettingPlayer          = errors.New("error getting player from minigames database")
	ErrMiniGamesDBUpdatingPlayerExp      = errors.New("error updating player exp in minigames database")
	ErrMiniGamesDBUpdatingPlayerRank     = errors.New("error updating player rank in minigames database")
	ErrMiniGamesDBUpdatingPlayerCoins    = errors.New("error updating player coins in minigames database")
	ErrMiniGamesDBUpdatingPlayerPlaytime = errors.New("error updating player playtime in minigames database")
	ErrMiniGamesDBUpdatingPlayerLastSeen = errors.New("error updating player last seen in minigames database")

	ErrMiniGamesDBGettingPlayerBanInfo = errors.New("error getting player's ban info")
	ErrMiniGamesDBBanningPlayer        = errors.New("error inserting player ban info to minigames database")
	ErrMiniGamesDBUnbanningPlayer      = errors.New("error updating player ban status in minigames database")

	ErrMiniGamesDBGettingPlayerMuteInfo = errors.New("error getting player's mute info")
	ErrMiniGamesDBMutingPlayer          = errors.New("error inserting player mute info to minigames database")
	ErrMiniGamesDBUnmutingPlayer        = errors.New("error updating player mute status in minigames database")
)
