// Copyright 2020 Joseph Cobhams. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE.md file.
//
package gomono

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

//Auth Endpoints

//ExchangeToken - https://docs.mono.co/reference#authentication-endpoint
func (g *gomono) ExchangeToken(code string) (string, error) {
	if code == "" {
		return "", errors.New("gomono: Code cannot be blank")
	}

	payload, err := g.preparePayload(map[string]string{"code": code})
	if err != nil {
		return "", err
	}

	respTarget := make(map[string]string)

	err = g.makeRequest("POST", fmt.Sprintf("%v/account/auth", g.apiUrl), payload, nil, &respTarget)
	if err != nil {
		return "", err
	}
	return respTarget["id"], nil
}

//Account Endpoints

//Information - https://docs.mono.co/reference#bank-account-details
func (g *gomono) Information(id string) (*InformationResponse, error) {
	if id == "" {
		return nil, errors.New("gomono: ID is required")
	}

	var respTarget InformationResponse
	err := g.makeRequest("POST", fmt.Sprintf("%v/accounts/%v", g.apiUrl, id), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}

	return &respTarget, nil
}

//Statement - https://docs.mono.co/reference#bank-statement
func (g *gomono) Statement(id, period, output string) (*StatementResponse, error) {
	if id == "" {
		return nil, errors.New("gomono: ID is required")
	}

	output = strings.ToLower(output)
	if output != "" && output != "pdf" && output != "json" {
		return nil, errors.New("gomono: only json or pdf output supported")
	}

	params := url.Values{}
	if output != "" {
		params.Add("output", output)
	}

	if period != "" {
		params.Add("period", period)
	}

	endpoint := fmt.Sprintf("%v/accounts/%v/statement?%v", g.apiUrl, id, params.Encode())

	var result StatementResponse

	switch output {
	case "pdf":
		var pdfRespTarget StatementResponsePdf
		err := g.makeRequest("GET", endpoint, nil, nil, &pdfRespTarget)
		if err != nil {
			return nil, err
		}
		result.PDF = &pdfRespTarget
	case "json":
		var jsonRespTarget StatementResponseJson
		err := g.makeRequest("POST", endpoint, nil, nil, &jsonRespTarget)
		if err != nil {
			return nil, err
		}
		result.JSON = &jsonRespTarget
	}

	return &result, nil
}

//PdfStatementJobStatus - https://docs.mono.co/reference#poll-statement-status
func (g *gomono) PdfStatementJobStatus(id, jobId string) (*StatementResponsePdf, error) {
	if id == "" {
		return nil, errors.New("gomono: ID is required")
	}

	if jobId == "" {
		return nil, errors.New("gomono: JOBID is required")
	}

	var respTarget StatementResponsePdf
	err := g.makeRequest("GET", fmt.Sprintf("%v/accounts/%v/statement/jobs/%v?", g.apiUrl, id, jobId), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}

	return &respTarget, nil
}

// User Endpoints

//Transactions - https://docs.mono.co/reference#poll-statement-status
func (g *gomono) Transactions(id, start, end, narration, tnxType string, paginate bool) (*TransactionsResponse, error) {
	if id == "" {
		return nil, errors.New("gomono: ID is required")
	}

	params := url.Values{}
	if start != "" {
		params.Add("start", start)
	}

	if end != "" {
		params.Add("end", end)
	}

	if narration != "" {
		params.Add("narration", narration)
	}

	if tnxType != "" {
		params.Add("type", tnxType)
	}

	params.Add("paginate", strconv.FormatBool(paginate))

	var respTarget TransactionsResponse
	err := g.makeRequest("GET", fmt.Sprintf("%v/accounts/%v/transactions?%v", g.apiUrl, id, params.Encode()), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}
	return &respTarget, nil
}

//CreditTransactions - https://docs.mono.co/reference#credits
func (g *gomono) CreditTransactions(id string) (*TransactionByTypeResponse, error) {
	return g.transactionByType(id, "credit")
}

//DebitTransactions - https://docs.mono.co/reference#debits
func (g *gomono) DebitTransactions(id string) (*TransactionByTypeResponse, error) {
	return g.transactionByType(id, "debit")
}

func (g *gomono) transactionByType(id, tnxType string) (*TransactionByTypeResponse, error) {
	if id == "" {
		return nil, errors.New("gomono: ID is required")
	}

	var respTarget TransactionByTypeResponse
	err := g.makeRequest("GET", fmt.Sprintf("%v/accounts/%v/%v", g.apiUrl, id, tnxType), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}
	return &respTarget, nil
}

//Income - https://docs.mono.co/reference#income
func (g *gomono) Income(id string) (*IncomeResponse, error) {
	if id == "" {
		return nil, errors.New("gomono: ID is required")
	}

	var respTarget IncomeResponse
	err := g.makeRequest("GET", fmt.Sprintf("%v/accounts/%v/income", g.apiUrl, id), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}
	return &respTarget, nil
}

//Identity - https://docs.mono.co/reference#identity
func (g *gomono) Identity(id string) (*IdentityResponse, error) {
	if id == "" {
		return nil, errors.New("gomono: ID is required")
	}

	var respTarget IdentityResponse
	err := g.makeRequest("GET", fmt.Sprintf("%v/accounts/%v/identity", g.apiUrl, id), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}
	return &respTarget, nil
}

//MISC Endpoints

//Institutions - https://docs.mono.co/reference#list-institutions
func (g *gomono) Institutions() (*InstitutionsResponse, error) {
	var respTarget []Institution
	err := g.makeRequest("GET", fmt.Sprintf("%v/coverage", g.apiUrl), nil, nil, &respTarget)
	if err != nil {
		return nil, err
	}
	return &InstitutionsResponse{
		Institutions: respTarget,
	}, nil
}

func (g *gomono) LookupBVN(bvn string) (*IdentityResponse, error) {
	if bvn == "" {
		return nil, errors.New("gomono: BVN is required")
	}

	payload, err := g.preparePayload(map[string]string{"bvn": bvn})
	if err != nil {
		return nil, err
	}

	var respTarget IdentityResponse
	err = g.makeRequest("POST", fmt.Sprintf("%v/v1/lookup/bvn/identity", g.apiUrl), payload, nil, &respTarget)
	if err != nil {
		return nil, err
	}
	return &respTarget, nil
}
