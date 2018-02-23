package fb

import (
	"net/http"
)

// ListSystemUsersReq lists the system users and admin system users for the business (adminToken must belong
// to an admin of the business or to an admin system user). The ID of each user returned is an app-scoped
// user ID. Use the SystemUserList type for responses.
func ListSystemUsersReq(adminToken, businessID string) *http.Request {
	return ReqSetup("GET", businessID+"/system_users", adminToken, nil)
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
		ID                 string `json:"id"`
		Name               string `json:"name"`
		AssignedAdAccounts struct {
			Data []struct {
				ID        string `json:"id"`
				AccountID string `json:"account_id"`
				Role      string `json:"role"`
			} `json:"data"`
		} `json:"assigned_ad_accounts"`
		AssignedPages struct {
			Data []struct {
				ID   string `json:"id"`
				Role string `json:"role"`
			} `json:"data"`
		} `json:"assigned_pages"`
	} `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // nil if no error is given by FB
}

// InstallSystemUserAppReq installs an app for a system user. The appUserID must be an app-scoped system user iD,
// which you can get with ListSystemUsersReq (adminToken must belong to an admin of the business or to an admin system user)
func InstallSystemUserAppReq(adminToken, appID, appUserID string) *http.Request {
	return ReqSetup("POST", appUserID+"/applications", adminToken, nil,
		&ParamStrStr{"business_app", appID})
}

type InstallSystemUserResponse bool      // TODO: correct? or wrapped somehow?
type InstallSystemUserResponse2 struct { // TODO: or this?
	Data  bool         `json:"data"`
	Error *ErrResponse `json:"error"` // nil if no error is given by FB
}
