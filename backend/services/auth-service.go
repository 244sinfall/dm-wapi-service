package services

import (
	"bytes"
	"context"
	permissions "darkmoon-wapi-service/permissions"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	Id int `json:"id"`
}

type UserScope struct {
	Root          bool  `json:"root"`
	SecurityLevel int   `json:"securityLevel"`
	RBAC          []int `json:"rbac"`
}

type AuthenticatedUser struct {
	User  User      `json:"user"`
	Scope UserScope `json:"scope"`
}

type IntegrationAuthenticatedUser struct {
	*AuthenticatedUser
	IntegrationUserId string `json:"integrationUserId"`
	Permission        int    `json:"permission"`
}

func (a *AuthenticatedUser) GetPermission() int {
	if a.Scope.Root || a.Scope.SecurityLevel == 3 {
		return permissions.AdminPermission
	}
	var arbiter = false
	for _, v := range a.Scope.RBAC {
		if v == 1031 {
			arbiter = true
		}
	}
	if a.Scope.SecurityLevel == 3 || (a.Scope.SecurityLevel == 2 && arbiter) {
		return permissions.AdminPermission
	}
	if a.Scope.SecurityLevel == 1 {
		return permissions.GmPermission
	}
	return permissions.PlayerPermission
}

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

type NotConnectedError struct{}

func (e NotConnectedError) Error() string {
	return "Account is not connected to Darkmoon"
}

type RevokedError struct{}

func (e *RevokedError) Error() string {
	return "User integration seems to be revoked"
}

var authHttpClient = &http.Client{}

func authenticateByRefreshToken(token string) (string, error) {
	data := []byte(`{"Token": "` + token + `"}`)
	request, err := http.NewRequest("POST", os.Getenv("DM_API_AUTH_SERVICE_REFRESH"), bytes.NewBuffer(data))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Fail on request creation refresh token: %v", err)
		return "", err
	}
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("Fail on request execution refresh token: %v", err)
		return "", err
	}
	var pair = new(authServiceTokenKeyPair)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&pair)
	if err != nil {
		fmt.Printf("Fail on response decode refresh token: %v", err)
		return "", err
	}
	response.Body.Close()
	return pair.AccessToken, err
}

func authenticateByApiKey(api_key string) (*authServiceTokenKeyPair, error) {
	request, err := http.NewRequest("GET", os.Getenv("DM_API_AUTH_SERVICE_BY_API_KEY"), nil)
	request.Header.Add("x-integration-key", os.Getenv("DM_API_AUTH_SERVICE_INTEGRATION_KEY"))
	request.Header.Add("x-user-key", api_key)
	request.Close = true

	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	if response.StatusCode == 401 {
		return nil, &RevokedError{}
	}
	var pair = new(authServiceTokenKeyPair)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&pair)
	response.Body.Close()
	return pair, err
}

func authenticateByAccessToken(access_token string) (*AuthenticatedUser, error) {
	request, err := http.NewRequest("GET", os.Getenv("DM_API_AUTH_SERVICE_ME"), nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+access_token)
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	var authUser = new(AuthenticatedUser)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&authUser)
	response.Body.Close()
	return authUser, err
}

func ConnectToDarkmoon(fbaccess string, code string, a *auth.Client, f *firestore.Client, ctx context.Context) (*IntegrationAuthenticatedUser, error) {
	token, err := a.VerifyIDToken(ctx, fbaccess)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	data := []byte(`{"integrationUserId": "` + token.UID + `"}`)
	request, err := http.NewRequest("POST", os.Getenv("DM_API_AUTH_SERVICE_CONNECT"), bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	request.Header.Add("x-integration-key", os.Getenv("DM_API_AUTH_SERVICE_INTEGRATION_KEY"))
	request.Header.Add("x-user-key", code)
	request.Header.Add("Content-Type", "application/json")
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	var integrationStatus = new(authConnectResponse)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&integrationStatus)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	response.Body.Close()
	if !integrationStatus.Success {
		fmt.Printf("%v\n", integrationStatus.Reason)
		return nil, errors.New(integrationStatus.Reason)
	}
	pair, err := authenticateByApiKey(code)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	newDoc := f.Collection("authService").Doc(token.UID)
	var newDocData = authServiceDoc{}
	newDocData.ApiKey = code
	newDocData.RefreshToken = pair.RefreshToken
	newDoc.Create(ctx, newDocData)
	authUser, err := authenticateByAccessToken(pair.AccessToken)
	return &IntegrationAuthenticatedUser{AuthenticatedUser: authUser, IntegrationUserId: token.UID, Permission: authUser.GetPermission()}, err
}

func Authenticate(fbaccess string, a *auth.Client, f *firestore.Client, ctx context.Context) (*IntegrationAuthenticatedUser, error) {
	token, err := a.VerifyIDToken(ctx, fbaccess)
	if err != nil {
		return nil, err
	}
	docRef := f.Doc("authService/" + token.UID)
	permInfo, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, &NotConnectedError{}
		}
		return nil, err
	}
	data := permInfo.Data()
	refreshToken := data["RefreshToken"].(string)
	api_key := data["ApiKey"].(string)
	if api_key == "" {
		return nil, &NotConnectedError{}
	}
	access_token := ""
	if refreshToken != "" {
		access_token, err = authenticateByRefreshToken(refreshToken)
	}
	if access_token == "" || err != nil {
		pair, err := authenticateByApiKey(api_key)
		if err != nil {
			if (errors.Is(err, &RevokedError{})) {
				ref := f.Doc("authService/" + token.UID)
				ref.Delete(ctx)
			}
			return nil, err
		}
		access_token = pair.AccessToken
		var newDoc = authServiceDoc{}
		newDoc.ApiKey = api_key
		newDoc.RefreshToken = pair.RefreshToken
		docRef.Set(ctx, newDoc)
	}
	authUser, err := authenticateByAccessToken(access_token)
	return &IntegrationAuthenticatedUser{AuthenticatedUser: authUser, IntegrationUserId: token.UID, Permission: authUser.GetPermission()}, err
}
