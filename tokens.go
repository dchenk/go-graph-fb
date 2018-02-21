package fb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

// This is how: https://developers.facebook.com/docs/marketing-api/businessmanager/systemuser/#systemusertoken
func CreateSystemTokenReq(userToken, appSecretProof, appID string, scope []string) *http.Request {
	return ReqSetup("GET", "", userToken, nil,
		&ParamStrStr{"business_app", appID},
		&ParamStrStr{"appsecret_proof", appSecretProof},
		&ParamStrStr{"scope", strings.Join(scope, ",")})
}

// https://developers.facebook.com/docs/graph-api/reference/v2.12/debug_token
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
