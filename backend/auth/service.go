package auth

import (
	"bytes"
	"darkmoon-wapi-service/globals"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var authHttpClient = &http.Client{}


func authenticateByRefreshToken(token string) (*authServiceTokenKeyPair, error) {
	data := []byte(`{"Token": "` + token + `"}`)
	request, err := http.NewRequest("POST", os.Getenv("DM_API_AUTH_SERVICE_REFRESH"), bytes.NewBuffer(data))
	request.Header.Add("Content-Type", "application/json")
	request.Close = true
	if err != nil {
		fmt.Printf("Fail on request creation refresh token: %v", err)
		return nil, err
	}
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("Fail on request execution refresh token: %v", err)
		return nil, err
	}
	var pair = new(authServiceTokenKeyPair)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&pair)
	if err != nil {
		fmt.Printf("Fail on response decode refresh token: %v", err)
		return nil, err
	}
	response.Body.Close()
	return pair, err
}

func authenticateByApiKey(api_key string) (*authServiceTokenKeyPair, error) {
	request, err := http.NewRequest("GET", os.Getenv("DM_API_AUTH_SERVICE_BY_API_KEY"), nil)
	request.Header.Add("x-integration-key", os.Getenv("DM_API_AUTH_SERVICE_INTEGRATION_KEY"))
	request.Header.Add("x-user-key", api_key)
	request.Close = true
	if err != nil {
		fmt.Printf("Fail on api key request creation: %v\n", err)
		return nil, err
	}
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("Fail on api key request execution: %v\n", err)
		return nil, err
	}
	if response.StatusCode == 401 {
		return nil, &revokedError{}
	}
	var pair = new(authServiceTokenKeyPair)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&pair)
	response.Body.Close()
	return pair, err
}

func getUser(access_token string) (*authenticatedUser, error) {
	request, err := http.NewRequest("GET", os.Getenv("DM_API_AUTH_SERVICE_ME"), nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+access_token)
	request.Close = true
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("Fail on user request execution:  %v\n", err)
		return nil, err
	}
	var authUser = new(authenticatedUser)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&authUser)
	response.Body.Close()
	return authUser, err
}

func connectToDarkmoon(integrationUserId string, code string) (*WapiAuthenticatedUser, error) {

	data := []byte(`{"integrationUserId": "` + integrationUserId + `"}`)
	request, err := http.NewRequest("POST", os.Getenv("DM_API_AUTH_SERVICE_CONNECT"), bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Fail on connect request creation: %v\n", err)
		return nil, err
	}
	request.Header.Add("x-integration-key", os.Getenv("DM_API_AUTH_SERVICE_INTEGRATION_KEY"))
	request.Header.Add("x-user-key", code)
	request.Header.Add("Content-Type", "application/json")
	request.Close = true
	response, err := authHttpClient.Do(request)
	if err != nil {
		fmt.Printf("Fail on connect request execution: %v\n", err)
		return nil, err
	}
	var integrationStatus = new(authConnectResponse)
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&integrationStatus)
	if err != nil {
		fmt.Printf("Fail on auth connect response decode: %v\n", err)
		return nil, err
	}
	response.Body.Close()
	if !integrationStatus.Success {
		fmt.Printf("%v\n", integrationStatus.Reason)
		return nil, errors.New(integrationStatus.Reason)
	}
	pair, err := authenticateByApiKey(code)
	if err != nil {
		fmt.Printf("Fail on using new api key: %v\n", err)
		return nil, err
	}

	newDoc := globals.GetFirestore().Collection("authService").Doc(integrationUserId)
	var newDocData = authServiceDoc{}
	newDocData.ApiKey = code
	newDocData.RefreshToken = pair.RefreshToken
	newDoc.Create(globals.GetGlobalContext(), newDocData)
	authUser, err := getUser(pair.AccessToken)
	return &WapiAuthenticatedUser{UserId: authUser.User.Id, IntegrationUserId: integrationUserId, Permission: authUser.GetPermission()}, err
}

func Authenticate(fbaccess string) (*WapiAuthenticatedUser, error) {
	token, err := globals.GetAuth().VerifyIDToken(globals.GetGlobalContext(), fbaccess)
	if err != nil {
		return nil, err
	}
	f := globals.GetFirestore()
	docRef := f.Doc("authService/" + token.UID)
	permInfo, err := docRef.Get(globals.GetGlobalContext())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, &notConnectedError{}
		}
		return nil, err
	}
	data := permInfo.Data()
	refreshToken := data["RefreshToken"].(string)
	api_key := data["ApiKey"].(string)
	if api_key == "" {
		docRef.Delete(globals.GetGlobalContext())
		return nil, &notConnectedError{}
	}
	keypair := new(authServiceTokenKeyPair)

	if refreshToken != "" {
		keypair, err = authenticateByRefreshToken(refreshToken)
	}
	if(keypair == nil || err != nil) {
		keypair, err = authenticateByApiKey(api_key)
		if err != nil {
			if (errors.Is(err, &revokedError{})) {
				docRef.Delete(globals.GetGlobalContext())
			}
			return nil, err
		}
		var newDoc = authServiceDoc{}
		newDoc.ApiKey = api_key
		newDoc.RefreshToken = keypair.RefreshToken
		docRef.Set(globals.GetGlobalContext(), newDoc)
	}

	authUser, err := getUser(keypair.AccessToken)
	return &WapiAuthenticatedUser{UserId: authUser.User.Id, IntegrationUserId: token.UID, Permission: authUser.GetPermission()}, err
}
