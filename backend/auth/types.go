package auth

type authServiceDoc struct {
	ApiKey       string `json:"api_key"`
	RefreshToken string `json:"refresh_token"`
}

type authServiceTokenKeyPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type authConnectResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

type notConnectedError struct{}

func (e notConnectedError) Error() string {
	return "Account is not connected to Darkmoon"
}

type revokedError struct{}

func (e *revokedError) Error() string {
	return "User integration seems to be revoked"
}

const (
	adminPermission   = 3
	arbiterPermission = 2
	gmPermission      = 1
	playerPermission  = 0
)


type user struct {
	Id int `json:"id"`
}

type userScope struct {
	Root          bool  `json:"root"`
	SecurityLevel int   `json:"securityLevel"`
	RBAC          []int `json:"rbac"`
}

type authenticatedUser struct {
	User  user      `json:"user"`
	Scope userScope `json:"scope"`
}

type connectRequestBody struct {
	Code string `json:"code"`
}

type resetRequestBody struct {
	Email string `json:"email"`
}

type WapiAuthenticatedUser struct {
	UserId            int    `json:"userId"`
	IntegrationUserId string `json:"integrationUserId"`
	Permission        int    `json:"permission"`
}

func (a *WapiAuthenticatedUser) IsAdmin() bool {
	return a.Permission >= adminPermission
}

func (a *WapiAuthenticatedUser) IsArbiter() bool {
	return a.Permission >= arbiterPermission
}

func (a *WapiAuthenticatedUser) IsGM() bool {
	return a.Permission >= gmPermission
}

func (a *authenticatedUser) GetPermission() int {
	if a.Scope.Root || a.Scope.SecurityLevel == 3 {
		return adminPermission
	}
	var arbiter = false
	for _, v := range a.Scope.RBAC {
		if v == 1031 {
			arbiter = true
		}
	}
	if a.Scope.SecurityLevel == 2 && arbiter {
		return arbiterPermission
	}
	if a.Scope.SecurityLevel == 1 {
		return gmPermission
	}
	return playerPermission
}