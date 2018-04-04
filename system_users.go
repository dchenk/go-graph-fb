package fb

import (
	"net/http"
)

// ListSystemUsersReq lists the system users and admin system users for the business (adminToken must belong to
// an admin of the business or to an admin system user). The ID of each user returned is an app-scoped user ID.
// Use the SystemUserList type for responses. The fields parameter specifies which fields to show for the users;
// the default fields (if given nil) are given in the SystemUserList struct.
func ListSystemUsersReq(adminToken, businessID string, fields []string) *http.Request {
	if fields == nil {
		fields = []string{"id", "name", "assigned_ad_accounts{name,account_id,role}", "assigned_pages{id,name,role,picture{url}}"}
	}
	return Req(http.MethodGet, businessID+"/system_users", adminToken, fields)
}

// SystemUserList sample payload:
//	{
//		"data":[
//			{
//				"id":"1000081799813",
//				"name":"Reporting server",
//				"assigned_ad_accounts": {
//					"data": [
//						{
//							"id":"act_XXXXX",
//							"account_id":"XXXXXXXXX",
//							"role":"ADMIN"
//						}
//					]
//				},
//				"assigned_pages": {
//					"data": [
//						{
//							"id":"1750248626186",
//							"role":"INSIGHTS_ANALYST"
//						}
//					]
//				}
//			}
//		]
//	}
type SystemUserList struct {
	Data []struct {
		ID                 string                 `json:"id"`
		Name               string                 `json:"name"`
		AssignedAdAccounts AssignedAdAccountsList `json:"assigned_ad_accounts"`
		AssignedPages      AssignedPagesList      `json:"assigned_pages"`
	} `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // nil if no error is given by FB
}

type AssignedAdAccountsList struct {
	Data []struct {
		AccountID string `json:"account_id"`
		Name      string `json:"name"`
		Role      string `json:"role"`
	} `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // useful only when requesting list alone
}

type AssignedPagesList struct {
	Data []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Role    string `json:"role"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	} `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // useful only when requesting list alone
}

// InstallSystemUserAppReq installs an app for a system user. The appUserID must be an app-scoped system user iD, which
// you can get with ListSystemUsersReq (adminToken must belong to an admin of the business or to an admin system user).
func InstallSystemUserAppReq(adminAccessToken, appID, appUserID string) *http.Request {
	return Req(http.MethodPost, appUserID+"/applications", adminAccessToken, nil, &ParamStrStr{"business_app", appID})
}
