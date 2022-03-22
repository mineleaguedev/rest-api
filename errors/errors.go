package errors

import "errors"

// DATABASE
var (
	ErrDBGettingLastInsertId = errors.New("error getting last insert id from insert database query")
	ErrDBGettingRowsAffected = errors.New("error getting rows affected from update database query")
)

// GENERAL
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

	ErrNotEnoughMoney       = errors.New("error not enough money")
	ErrDBMoneySubtraction   = errors.New("error money subtraction in general database")
	ErrDBMoneyAddition      = errors.New("error money addition in general database")
	ErrDBSavingTransferInfo = errors.New("error saving transfer info into general database")

	ErrDBRegisteringUser      = errors.New("error inserting user into general database")
	ErrDBGettingUser          = errors.New("error getting user from general database")
	ErrDBUpdatingUserPassword = errors.New("error updating password in general database")

	ErrS3UploadingSkin  = errors.New("error uploading skin into skins bucket")
	ErrS3DeletingSkin   = errors.New("error deleting skin from skins bucket")
	ErrS3UploadingCloak = errors.New("error uploading cloak into cloaks bucket")
	ErrS3DeletingCloak  = errors.New("error deleting cloak from cloaks bucket")

	ErrSendingEmail = errors.New("error sending email")
)

// MINIGAMES
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

	ErrMiniGamesDBCreatingPlayer         = errors.New("error inserting player into minigames database")
	ErrMiniGamesDBGettingPlayer          = errors.New("error getting player from minigames database")
	ErrMiniGamesDBUpdatingPlayerExp      = errors.New("error updating player exp in minigames database")
	ErrMiniGamesDBUpdatingPlayerRank     = errors.New("error updating player rank in minigames database")
	ErrMiniGamesDBUpdatingPlayerCoins    = errors.New("error updating player coins in minigames database")
	ErrMiniGamesDBUpdatingPlayerPlaytime = errors.New("error updating player playtime in minigames database")
	ErrMiniGamesDBUpdatingPlayerLastSeen = errors.New("error updating player last seen in minigames database")

	ErrMiniGamesDBGettingPlayerBanInfo = errors.New("error getting player's ban info from minigame database")
	ErrMiniGamesDBBanningPlayer        = errors.New("error inserting player ban info into minigames database")
	ErrMiniGamesDBUnbanningPlayer      = errors.New("error updating player ban status in minigames database")

	ErrMiniGamesDBGettingPlayerMuteInfo = errors.New("error getting player's mute info from minigame database")
	ErrMiniGamesDBMutingPlayer          = errors.New("error inserting player mute info into minigames database")
	ErrMiniGamesDBUnmutingPlayer        = errors.New("error updating player mute status in minigames database")

	ErrMissingMapCreateValues = errors.New("missing map create values")

	ErrInvalidMapWorldFile  = errors.New("invalid map world file")
	ErrInvalidMapConfigFile = errors.New("invalid map config file")

	ErrS3GettingMapsList                      = errors.New("error getting maps list from s3 maps bucket")
	ErrS3GettingMiniGameMapsList              = errors.New("error getting minigame maps list from s3 maps bucket")
	ErrS3EmptyMiniGameMapsList                = errors.New("empty list of minigame maps in s3 maps bucket")
	ErrS3GettingMiniGameFormatMapsList        = errors.New("error getting minigame format maps list from s3 maps bucket")
	ErrS3EmptyMiniGameFormatMapsList          = errors.New("empty list of minigame format maps in s3 maps bucket")
	ErrS3GettingMiniGameFormatMapVersionsList = errors.New("error getting minigame format map versions list from s3 maps bucket")
	ErrS3EmptyMiniGameFormatMapVersionsList   = errors.New("empty list of minigame format map versions in s3 maps bucket")
	ErrS3CreatingMap                          = errors.New("error creating map in s3 maps bucket")
	ErrS3DownloadingMapWorld                  = errors.New("error downloading map world from s3 maps bucket")
	ErrS3DownloadingMapConfig                 = errors.New("error downloading map config from s3 maps bucket")

	ErrS3GettingPluginsList        = errors.New("error getting plugins list from s3 plugins bucket")
	ErrS3GettingPluginVersionsList = errors.New("error getting plugin versions list from s3 plugins bucket")
	ErrS3EmptyPluginVersionsList   = errors.New("empty list of plugin versions in s3 plugins bucket")
)
