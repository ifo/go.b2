package b2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type B2 struct {
	AccountID          string
	ApplicationKey     string
	AuthorizationToken string
	ApiUrl             string
	DownloadUrl        string
}

type authResponse struct {
	AccountID          string `json:"accountId"`
	AuthorizationToken string `json:"authorizationToken"`
	ApiUrl             string `json:"apiUrl"`
	DownloadUrl        string `json:"downloadUrl"`
}

type ErrorResponse struct {
	Status  int64  `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("Status: %d, Code: %s, Message: %s", e.Status, e.Code, e.Message)
}

func CreateB2(accountId, appKey string) (*B2, error) {
	req, err := CreateRequest("GET", "https://api.backblaze.com/b2api/v1/b2_authorize_account", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(accountId, appKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b2 := &B2{
		AccountID:      accountId,
		ApplicationKey: appKey,
	}
	return b2.parseCreateB2Response(resp)
}

func (b2 *B2) parseCreateB2Response(resp *http.Response) (*B2, error) {
	authResp := &authResponse{}
	err := ParseResponse(resp, authResp)
	if err != nil {
		return nil, err
	}

	b2.AuthorizationToken = authResp.AuthorizationToken
	b2.ApiUrl = authResp.ApiUrl
	b2.DownloadUrl = authResp.DownloadUrl

	return b2, nil
}

func CreateRequest(method, url string, request interface{}) (*http.Request, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, bytes.NewReader(reqBody))
}

func GetBzInfoHeaders(resp *http.Response) map[string]string {
	out := map[string]string{}
	for k, v := range resp.Header {
		if strings.HasPrefix(k, "X-Bz-Info-") {
			// strip Bz prefix and grab first header
			out[k[10:]] = v[0]
		}
	}
	return out
}

func ParseResponse(resp *http.Response, respBody interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return ParseResponseBody(resp, respBody)
	} else {
		return ParseErrorResponse(resp)
	}
}

func ParseResponseBody(resp *http.Response, respBody interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, respBody)
}

func ParseErrorResponse(resp *http.Response) error {
	errResp := &ErrorResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, errResp)
	if err != nil {
		return err
	}
	return errResp
}
