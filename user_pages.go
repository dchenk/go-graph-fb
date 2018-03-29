package fb

import "net/http"

// A UserPagesList response lists the pages belonging to a user.
type UserPagesList struct {
	Data []struct {
		AccessToken  string `json:"access_token"`
		ID           string `json:"id"`
		Category     string `json:"category"`
		CategoryList []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"category_list"`
		Name  string   `json:"name"`
		Perms []string `json:"perms"`
	} `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"` // nil if no error is given by FB
}

// ListUserPagesReq lists the pages belonging to a user using a user access token.
func ListUserPagesReq(accessToken string) *http.Request {
	return Req("GET", "me/accounts", accessToken, nil)
}
