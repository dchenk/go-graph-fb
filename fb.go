// Package fb provides helpers for using the Facebook Graph API.
package fb

import (
	"bytes"
	"encoding/json"
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
			Value json.RawMessage `json:"value"` // Not set for some endpoints.
		} `json:"changes"`
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

// NextPage makes a request to nextURL using the given client. If Client is nil, then http.DefaultClient is used.
func NextPage(nextURL string, client *http.Client) (*http.Response, error) {
	u, err := url.Parse(nextURL)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       nil,
		Host:       u.Host,
	}
	if client == nil {
		client = http.DefaultClient
	}
	return client.Do(req)
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

// Error gives the main details of the error code and message.
func (er *ErrResponse) Error() string {
	return fmt.Sprintf("fb: error code %d, subcode %d; msg: %s", er.Code, er.ErrorSubcode, er.Message)
}

// UserErrMessage summarizes the error in a way that can be displayed to users.
func (er *ErrResponse) UserErrMessage() string {
	title := er.ErrorUserTitle
	if title == "" { // Facebook sometimes returns no user title.
		title = er.Type
	}
	msg := er.ErrorUserMessage
	if msg == "" { // Facebook sometimes returns no user message.
		msg = er.Message
	}
	return fmt.Sprintf("%s (error code %d, subcode %d): %s", title, er.Code, er.ErrorSubcode, msg)
}

// IsErrResponse says if the error is of type *ErrResponse, which Facebook can send as part of a response payload.
func IsErrResponse(err error) bool {
	_, v := err.(*ErrResponse)
	return v
}

// Req sets up a request to the Facebook API but does not run it. The method should one of GET, POST, or DELETE.
// The nodeEdge parameter should not have a leading slash or the Graph API version (currently set to 2.12).
// Leave the fields slice empty or nil to not specify a fields parameter.
func Req(method, nodeEdge string, accessToken string, fields []string, params ...Param) *http.Request {
	r := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme: "https",
			Host:   "graph.facebook.com",
			Path:   "/v2.12/" + nodeEdge,
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

// ReqDo uses Req to set up the request and then runs Do on it. The client request timeout is set to 12 seconds.
func ReqDo(method, nodeEdge string, accessToken string, fields []string, params ...Param) (*http.Response, error) {
	return (&http.Client{Timeout: time.Second * 12}).Do(Req(method, nodeEdge, accessToken, fields, params...))
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

func (psi *ParamStrInt) Val() string { return strconv.FormatInt(psi.V, 10) }

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
