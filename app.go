package fb

import "net/http"

func ListAppSubscriptionsReq(appAccessToken, appID string) *http.Request {
	return Req("GET", appID+"/subscriptions", appAccessToken, nil)
}

type AppSubscriptionsList struct {
	Data []struct {
		Object      string `json:"object"`
		CallbackURL string `json:"callback_url"`
		Active      bool   `json:"active"`
		Fields      []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"fields"`
	} `json:"data"`
	Error *ErrResponse `json:"error"`
}
