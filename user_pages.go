package fb

import "net/http"

// A UserPagesList response lists the pages belonging to a user.
type UserPagesList struct {
	Data   []UserPage   `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // nil if no error is given by FB
}

type UserPage struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AccessToken  string `json:"access_token"`
	Category     string `json:"category"`
	CategoryList []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"category_list"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	Perms []string `json:"perms"`
}

// ListUserPagesReq lists the pages belonging to a user using a user access token with the
// fields: id,name,access_token,category,perms,picture{url}
func ListUserPagesReq(accessToken string) *http.Request {
	return ListUserPagesFieldsReq(accessToken, listUserPagesFields)
}

var listUserPagesFields = []string{"id", "name", "access_token", "category", "perms", "picture{url}"}

// ListUserPagesFieldsReq lists the pages belonging to a user using a user access token with
// the specified fields.
func ListUserPagesFieldsReq(accessToken string, fields []string) *http.Request {
	return Req("GET", "me/accounts", accessToken, fields)
}

// SubscribeAppToPageReq returns a request that can be used to subscribe an app to a page. A page access token belonging
// to the page must be used for this.
func SubscribeAppToPageReq(pageAccessToken, pageID string) *http.Request {
	return Req("POST", pageID+"/subscribed_apps", pageAccessToken, nil)
}
