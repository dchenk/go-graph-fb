package fb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

// CreateSystemTokenReq sets up an http.Request for getting a system user token.
// Info: https://developers.facebook.com/docs/marketing-api/businessmanager/systemuser/#systemusertoken
func CreateSystemTokenReq(userToken, systemUserID, appSecretProof, appID string, scope []string) *http.Request {
	return ReqSetup("POST", systemUserID+"/access_tokens", userToken, nil,
		&ParamStrStr{"business_app", appID},
		&ParamStrStr{"appsecret_proof", appSecretProof},
		&ParamStrStr{"scope", strings.Join(scope, ",")})
}

// A SystemToken is the format in which Facebook returns a system user token.
type SystemToken struct {
	AccessToken string       `json:"access_token"`
	Error       *ErrResponse `json:"error"` // nil if no error is given by FB
}

// A TokenDebug represents a Facebook response for the token debugging API.
// Info: https://developers.facebook.com/docs/graph-api/reference/v2.12/debug_token
type TokenDebug struct {
	Data struct {
		IsValid     bool     `json:"is_valid"`
		AppID       string   `json:"app_id"`
		Application string   `json:"application"`
		Type        string   `json:"type"`
		IssuedAt    int64    `json:"issued_at"`
		ExpiresAt   int64    `json:"expires_at"`
		Scopes      []string `json:"scopes"`
		UserID      string   `json:"user_id"`
		Error       struct { // Empty if no error.
			Code    int64  `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	Error *ErrResponse `json:"error"` // nil if no error is given by FB; if nil, no Data will be sent
}

// DebugToken sends a token debug request to Facebook and reads the response.
// If the client given is nil, then http.DefaultClient is used.
func DebugToken(accessToken, tokenToDebug string, client *http.Client) (*TokenDebug, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(DebugTokenReq(accessToken, tokenToDebug))
	if err != nil {
		return nil, err
	}
	info := new(TokenDebug)
	err = ReadResponse(resp, info)
	return info, err
}

// DebugTokenReq sets up an http.Request for debugging a token.
func DebugTokenReq(accessToken, tokenToDebug string) *http.Request {
	return ReqSetup("GET", "debug_token", accessToken, nil, &ParamStrStr{"input_token", tokenToDebug})
}

// AppsecretProof generates an app secret proof for an app. The userAccessToken must belong to an admin of the app.
// Info: https://developers.facebook.com/docs/graph-api/securing-requests/#appsecret_proof
func AppsecretProof(userAccessToken, appSecret string) (string, error) {
	sig := hmac.New(sha256.New, []byte(appSecret))
	if _, err := sig.Write([]byte(userAccessToken)); err != nil {
		return "", err
	}
	return hex.EncodeToString(sig.Sum(nil)), nil
}
