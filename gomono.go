// Copyright 2020 Joseph Cobhams. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE.md file.
//
package gomono

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type (
	Gomono interface {
		ExchangeToken(code string) (string, error)
		Information(id string) (*InformationResponse, error)
		Statement(id, period, output string) (*StatementResponse, error)
		PdfStatementJobStatus(id, jobId string) (*StatementResponsePdf, error)
		Transactions(id, start, end, narration, tnxType string, paginate bool) (*TransactionsResponse, error)
		CreditTransactions(id string) (*TransactionByTypeResponse, error)
		DebitTransactions(id string) (*TransactionByTypeResponse, error)
		Income(id string) (*IncomeResponse, error)
		Identity(id string) (*IdentityResponse, error)
		Institutions() (*InstitutionsResponse, error)
		LookupBVN(bvn string) (*IdentityResponse, error)
	}

	gomono struct {
		secretKey string
		client    *http.Client
		apiUrl    string
	}

	Error struct {
		Code     int
		Body     string
		Endpoint string
	}

	Config struct {
		SecretKey  string
		HttpClient *http.Client
		ApiUrl     string
	}

	header struct {
		Key   string
		Value string
	}
)

func New(cfg Config) (Gomono, error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	g := &gomono{
		secretKey: cfg.SecretKey,
		client:    cfg.HttpClient,
		apiUrl:    cfg.ApiUrl,
	}

	return g, nil
}

func NewDefaultConfig(secretKey string) Config {
	return Config{
		SecretKey: secretKey,
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		ApiUrl: "https://api.withmono.com",
	}
}

func validateConfig(cfg *Config) error {
	if cfg.SecretKey == "" {
		return errors.New("gomono: Missing Secret Key")
	}

	if cfg.HttpClient == nil {
		return errors.New("gomono: HTTP Client Cannot Be Nil")
	}

	if cfg.ApiUrl == "" {
		return errors.New("gomono: Missing API Url")
	}

	return nil
}

func (g *gomono) preparePayload(body interface{}) (io.Reader, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (g *gomono) makeRequest(method, url string, body io.Reader, headers []header, responseTarget interface{}) error {
	if reflect.TypeOf(responseTarget).Kind() != reflect.Ptr {
		return errors.New("gomono: responseTarget must be a pointer to a struct for JSON unmarshalling")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value)
	}
	req.Header.Set("mono-sec-key", g.secretKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Lib", "GoMono | v1 | github.com/jcobhams/gomono")

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		err = json.Unmarshal(b, responseTarget)
		if err != nil {
			return err
		}
		return nil
	}

	err = Error{
		Code:     resp.StatusCode,
		Body:     string(b),
		Endpoint: req.URL.String(),
	}
	return err
}

func (e Error) Error() string {
	return fmt.Sprintf("Request To %v Endpoint Failed With Status Code %v | Body: %v", e.Endpoint, e.Code, e.Body)
}
