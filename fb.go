// Package fb provides helpers for using the Facebook Graph API.
package fb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GraphResponseMe struct {
	ID    string       `json:"id"`
	Name  string       `json:"name"`
	Email string       `json:"email"`
	Error *ErrResponse `json:"error"`
}

type GraphResponsePage struct {
	Data struct {
		ID      string `json:"id"` // numeric string
		Name    string `json:"name"`
		Article string `json:"article"`
		Type    string `json:"type"`
	} `json:"data"`
	Paging CursorPaging `json:"paging"`
	Error  *ErrResponse `json:"error"`
}

// A WebhookNotif represents any webhook notification payload.
// Info https://developers.facebook.com/docs/graph-api/webhooks
type WebhookNotif struct {
	Object string `json:"object"` // enum{user, page, permissions, payments}
	Entry  []struct {
		ID            string   `json:"id"`
		ChangedFields []string `json:"changed_fields"` // Fields include, e.g., for Page "leadgen", "location", "messages", etc.
		Changes       []struct {
			Field string          `json:"field"`
			Value json.RawMessage `json:"changes"` // Not set for some endpoints.
		}
		Time int `json:"time"` // A Unix timestamp.
	} `json:"entry"`
}

// A LeadGenEntry is a Page webhook notification value for the field "leadgen".
type LeadGenEntry struct {
	AdID        string `json:"ad_id"`
	FormID      string `json:"form_id"`
	LeadgenID   string `json:"leadgen_id"`
	PageID      string `json:"page_id"`
	AdgroupID   string `json:"adgroup_id"`
	CreatedTime int64  `json:"created_time"`
}

type CursorPaging struct {
	Cursors struct {
		Before string `json:"before"`
		After  string `json:"after"`
	} `json:"cursors"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Limit    string `json:"limit"`
}

type TimePaging struct {
	Until    int64  `json:"until"`
	Since    int64  `json:"since"`
	Limit    int64  `json:"limit"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type OffsetPaging struct {
	Offset   int64  `json:"offset"`
	Limit    int64  `json:"limit"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type ErrResponse struct {
	Message          string `json:"message"`
	Type             string `json:"type"`
	Code             int64  `json:"code"`
	ErrorSubcode     int64  `json:"error_subcode"`
	ErrorUserTitle   string `json:"error_user_title"`
	ErrorUserMessage string `json:"error_user_message"`
	FbTraceID        string `json:"fbtrace_id"`
}

// FillFromMap fills out the struct from the map extracted out of the JSON response sent by the Facebook API.
func (er *ErrResponse) fill(m map[string]interface{}) {
	if v, ok := m["message"].(string); ok {
		er.Message = v
	}
	if v, ok := m["type"].(string); ok {
		er.Type = v
	}
	if v, ok := m["code"].(int64); ok {
		er.Code = v
	}
	if v, ok := m["error_subcode"].(int64); ok {
		er.ErrorSubcode = v
	}
	if v, ok := m["error_user_title"].(string); ok {
		er.ErrorUserTitle = v
	}
	if v, ok := m["error_user_message"].(string); ok {
		er.ErrorUserMessage = v
	}
	if v, ok := m["fbtrace_id"].(string); ok {
		er.FbTraceID = v
	}
}

func (er *ErrResponse) Error() string {
	return fmt.Sprintf("fb: error code %d, subcode %d; msg: %s", er.Code, er.ErrorSubcode, er.Message)
}

// ReqSetup sets up a request to the Facebook API but does not begin it.
// Leave the fields slice empty or nil to not specify a fields parameter.
func ReqSetup(method, nodeEdge string, accessToken string, fields []string, params ...Param) *http.Request {
	r := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme: "https",
			Host:   "graph.facebook.com",
			Path:   "/v2.11/" + nodeEdge,
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}
	if len(fields) > 0 {
		params = append(params, &ParamStrStr{
			K: "fields",
			V: strings.Join(fields, ","),
		})
	}
	if method == "POST" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		b := new(bytes.Buffer)
		b.WriteString(encodeParams(accessToken, params))
		r.ContentLength = int64(b.Len())
		buf := b.Bytes()
		r.Body = ioutil.NopCloser(b)
		r.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(buf)), nil
		}
	} else {
		r.URL.RawQuery = encodeParams(accessToken, params)
	}
	return r
}

// Req uses ReqSetup to set up the request and then runs Do on it. The method parameter should be in all
// capitals: GET, POST, or DELETE. The nodeEdge parameter should not have a leading slash or the Graph API
// version (currently at 2.11). An access token must always be provided. The timeout is set to 10 seconds.
func Req(method, nodeEdge string, accessToken string, fields []string, params ...Param) (*http.Response, error) {
	return (&http.Client{Timeout: time.Second * 10}).Do(ReqSetup(method, nodeEdge, accessToken, fields, params...))
}

// A Param is a key => value pair to be sent in the request.
type Param interface {
	Key() string
	Val() string
}

type ParamStrStr struct {
	K, V string
}

func (pss *ParamStrStr) Key() string { return pss.K }

func (pss *ParamStrStr) Val() string { return pss.V }

type ParamStrInt struct {
	K string
	V int64
}

func (psi *ParamStrInt) Key() string { return psi.K }

func (psi *ParamStrInt) Val() string {
	return strconv.FormatInt(psi.V, 10)
}

// encodeParams builds url.Values from the given Param elements. This function sets the access token parameter.
func encodeParams(accessToken string, params []Param) string {
	v := make(url.Values, len(params)+1)
	v.Set("access_token", accessToken)
	for _, p := range params {
		v.Set(p.Key(), p.Val())
	}
	return v.Encode()
}

// ReadResponse simply reads the response and decodes it into v, which should be a non-nil pointer to a struct that can take
// an error response (in the Facebook Graph way) or the actual response expected. This function closes the http.Response body
// upon returning.
func ReadResponse(res *http.Response, v interface{}) error {
	err := json.NewDecoder(res.Body).Decode(v)
	res.Body.Close()
	return err
}

// ReadResponseFill reads an *http.Response from a request into rf. The error returned is not nil if an only if something goes
// wrong reading or decoding the response body (in which case the rf parameter passed in is not used at all). The *ErrResponse
// is not nil if and only if the API response was read without problems (though it may have the "error" property not as a JSON
// object and still set the error value to not nil) but there is an "error" property set in the JSON response. This function
// calls the Fill method of rf to fill it out with the response data if there were no errors reading the response or no error
// message returned by the Graph API body. This function closes the http.Response body upon returning.
func ReadResponseFill(res *http.Response, rf RespFiller) (error, *ErrResponse) {

	defer res.Body.Close()

	bod := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&bod); err != nil {
		return fmt.Errorf("error decoding response: %s", err), nil
	}

	// Try to fill out an ErrResponse struct only if the Facebook API returned an error message in their special way.
	if errResp, ok := bod["error"]; ok {
		er := new(ErrResponse)
		msi, ok := errResp.(map[string]interface{})
		if ok {
			er.fill(msi)
		} else {
			if str, ok := errResp.(string); ok {
				er.Message = str
			}
			return errors.New("got an error response, but error value is not a JSON object; check ErrResponse.Message"), er
		}
		return nil, er
	}

	if rf != nil {
		rf.Fill(bod)
	}

	return nil, nil

}

// A RespFiller is able to take any kind of Facebook Graph API response (as a map[string]interface{}) and extract the details it
// wants out of it. When a RespFiller is passed into ReadResponse, its Fill method is called only if both no error occurred reading
// and decoding the response body and if the response body does not have "error" as a property in the JSON object.
type RespFiller interface {
	Fill(map[string]interface{})
}
