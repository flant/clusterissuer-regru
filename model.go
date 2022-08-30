package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

// APIResponse is the representation of an API response.
type APIResponse struct {
	Result string `json:"result"`

	Answer *Answer `json:"answer,omitempty"`

	ErrorCode string `json:"error_code,omitempty"`
	ErrorText string `json:"error_text,omitempty"`
}

func (a APIResponse) Error() string {
	return fmt.Sprintf("API %s: %s: %s", a.Result, a.ErrorCode, a.ErrorText)
}

// HasError returns an error is the response contains an error.
func (a APIResponse) HasError() error {
	if a.Result != successResult {
		return a
	}

	if a.Answer != nil {
		for _, domResp := range a.Answer.Domains {
			if domResp.Result != successResult {
				return domResp
			}
		}
	}

	return nil
}

// Answer is the representation of an API response answer.
type Answer struct {
	Domains []DomainResponse `json:"domains,omitempty"`
}

// DomainResponse is the representation of an API response answer domain.
type DomainResponse struct {
	Result string `json:"result"`

	DName string `json:"dname"`

	ErrorCode string `json:"error_code,omitempty"`
	ErrorText string `json:"error_text,omitempty"`
}

func (d DomainResponse) Error() string {
	return fmt.Sprintf("API %s: %s: %s", d.Result, d.ErrorCode, d.ErrorText)
}

// AddTxtRequest is the representation of the payload of a request to add a TXT record.
type AddTxtRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Domains           []Domain `json:"domains,omitempty"`
	SubDomain         string   `json:"subdomain,omitempty"`
	Text              string   `json:"text,omitempty"`
	OutputContentType string   `json:"output_content_type,omitempty"`
}

// RemoveRecordRequest is the representation of the payload of a request to remove a record.
type RemoveRecordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Domains           []Domain `json:"domains,omitempty"`
	SubDomain         string   `json:"subdomain,omitempty"`
	Content           string   `json:"content,omitempty"`
	RecordType        string   `json:"record_type,omitempty"`
	OutputContentType string   `json:"output_content_type,omitempty"`
}

// Domain is the representation of a Domain.
type Domain struct {
	DName string `json:"dname"`
}

func (c *RegruClient) doRequest(req *http.Request, readResponseBody bool) (int, []byte, error) {
	if c.dumpRequestResponse {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Printf("Request: %q\n", dump)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	if c.dumpRequestResponse {
		dump, _ := httputil.DumpResponse(res, true)
		fmt.Printf("Response: %q\n", dump)
	}

	if res.StatusCode == http.StatusOK && readResponseBody {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return 0, nil, err
		}
		return res.StatusCode, data, nil
	}

	return res.StatusCode, nil, nil
}
