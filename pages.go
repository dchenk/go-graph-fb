package fb

import "net/http"

// A UserPagesList response lists the pages belonging to a user.
type UserPagesList struct {
	Data   []UserPage   `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // nil if no error is given
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

// A SubscribeAppResponse represents the format in which a response indicates if an app successfully subscribed to a page.
type SubscribeAppResponse struct {
	Success bool         `json:"success"`
	Error   *ErrResponse `json:"error"` // nil if no error is given
}

// ListPageSubscribedAppsReq returns a request to query the Facebook apps that are subscribed to a page's events.
func ListPageSubscribedAppsReq(pageAccessToken, pageID string) *http.Request {
	return Req("GET", pageID+"/subscribed_apps", pageAccessToken, nil)
}

// A SubscribedAppsList response represents the list of apps subscribed to a page.
type SubscribedAppsList struct {
	Data []struct {
		Category string `json:"category"`
		Link     string `json:"link"`
		Name     string `json:"name"`
		ID       string `json:"id"`
	} `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // nil if no error is given
}

// PageLeadgenSetupReq returns a request to query the basic settings concerning a page's leadgen setup.
// Fields retrieved: id,name,leadgen_has_crm_integration,leadgen_forms{id,name,status}
func PageLeadgenSetupReq(pageAccessToken, pageID string) *http.Request {
	return Req("GET", pageID, pageAccessToken, leadgenSetupFields)
}

var leadgenSetupFields = []string{"id", "name", "leadgen_has_crm_integration", "leadgen_forms{id,name,status}"}

type PageLeadgenSetup struct {
	ID                       string `json:"id"` // The page ID
	Name                     string `json:"name"`
	LeadgenHasCrmIntegration bool   `json:"leadgen_has_crm_integration"`
	LeadgenForms             struct {
		Data   []PageLeadgenForm `json:"data"`
		Paging CursorPaging      `json:"paging"`
	} `json:"leadgen_forms"`
	Error *ErrResponse `json:"error"` // nil if no error is given
}

type PageLeadgenForm struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// A PageLeadgenFormList is used to read the response from a "next" page URL given to page through the list of
// lead ads forms belonging to a page.
type PageLeadgenFormList struct {
	Data   []PageLeadgenForm `json:"data"`
	Paging CursorPaging      `json:"paging"`
	Error  *ErrResponse      `json:"error"`
}
