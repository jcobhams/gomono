// Copyright 2020 Joseph Cobhams. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE.md file.
//
package gomono

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	testAccountId       = "5fc68b964bdcbe4eb164e852"
	testJobId           = "MvRh2vWwv5CGafudTivY"
	testSecretKey       = "TEST_SECRET_KEY"
	testMonoConnectCode = "TEST_MONO_CONNECT_CODE"
	mockServer          *httptest.Server
	client              Gomono
)

func TestMain(m *testing.M) {

	mockServer = testServer()
	client, _ = New(Config{
		SecretKey:  testSecretKey,
		HttpClient: &http.Client{Timeout: 1 * time.Second},
		ApiUrl:     mockServer.URL,
	})
	mockServer.Close()

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	g, err := New(NewDefaultConfig(""))
	assert.Nil(t, g)
	assert.NotNil(t, err)

	g, err = New(NewDefaultConfig(testSecretKey))
	assert.NotNil(t, g)
	assert.Nil(t, err)
}

func TestGomono_ExchangeToken(t *testing.T) {
	id, err := client.ExchangeToken("")
	assert.Empty(t, id)
	assert.NotNil(t, err)

	id, err = client.ExchangeToken(testMonoConnectCode)
	assert.NotEmpty(t, id)
	assert.Equal(t, testAccountId, id)
	assert.Nil(t, err)
}

func TestGomono_Information(t *testing.T) {
	r, err := client.Information("")
	assert.Nil(t, r)
	assert.NotNil(t, err)

	r, err = client.Information(testAccountId)
	assert.NotNil(t, r)
	assert.Equal(t, testAccountId, r.Account.ID)
	assert.Nil(t, err)
}

func TestGomono_Statement(t *testing.T) {
	r, err := client.Statement(testAccountId, "", "x")
	assert.Nil(t, r)
	assert.NotNil(t, err)

	r, err = client.Statement(testAccountId, "", "json")
	assert.NotNil(t, r)
	assert.Equal(t, 2, r.JSON.Meta.Count)
	assert.Equal(t, 2, len(r.JSON.Data))
	assert.Nil(t, r.PDF)
	assert.Nil(t, err)

	r, err = client.Statement(testAccountId, "", "pdf")
	assert.NotNil(t, r)
	assert.Equal(t, "BUILDING", r.PDF.Status)
	assert.NotEmpty(t, r.PDF.ID)
	assert.NotEmpty(t, r.PDF.Path)
	assert.Nil(t, r.JSON)
	assert.Nil(t, err)
}

func TestGomono_PdfStatementJobStatus(t *testing.T) {
	r, err := client.PdfStatementJobStatus(testAccountId, testJobId)
	assert.NotNil(t, r)
	assert.Equal(t, testJobId, r.ID)
	assert.Equal(t, "COMPLETE", r.Status)
	assert.Nil(t, err)
}

func TestGomono_Transactions(t *testing.T) {
	r, err := client.Transactions(testAccountId, "01-10-2020", "07-10-2020", "test", "debit", true)
	assert.NotNil(t, r)
	assert.Equal(t, 190, r.Paging.Total)
	assert.Equal(t, 2, r.Paging.Page)
	assert.Equal(t, 2, len(r.Data))
	assert.Nil(t, err)
}

func TestGomono_TransactionByType(t *testing.T) {
	r, err := client.CreditTransactions(testAccountId)
	assert.NotNil(t, r)
	assert.Equal(t, float64(2000000), r.Total)
	assert.Equal(t, 2, len(r.History))
	assert.Nil(t, err)

	r, err = client.DebitTransactions(testAccountId)
	assert.NotNil(t, r)
	assert.Equal(t, float64(1000000), r.Total)
	assert.Equal(t, 2, len(r.History))
	assert.Nil(t, err)
}

func TestGomono_Income(t *testing.T) {
	r, err := client.Income(testAccountId)
	assert.NotNil(t, r)
	assert.Equal(t, "INCOME", r.Type)
	assert.Equal(t, float64(59700000), r.Amount)
	assert.Nil(t, err)
}

func TestGomono_Identity(t *testing.T) {
	r, err := client.Identity(testAccountId)
	assert.NotNil(t, r)
	assert.Equal(t, "ABDULHAMID", r.FirstName)
	assert.Equal(t, "HASSAN", r.LastName)
	assert.Equal(t, "NO", r.WatchListed)
	assert.Equal(t, "06-May-1996", r.DateOfBirth)
	assert.Nil(t, err)
}

func TestGomono_Institutions(t *testing.T) {
	r, err := client.Institutions()
	assert.NotNil(t, r)
	assert.Equal(t, "GTBank", r.Institutions[0].Name)
	assert.Equal(t, 4, len(r.Institutions))
	assert.Nil(t, err)
}

func TestGomono_LookupBVN(t *testing.T) {
	r, err := client.LookupBVN("1234567897418")
	assert.NotNil(t, r)
	assert.Equal(t, "ABDULHAMID", r.FirstName)
	assert.Equal(t, "HASSAN", r.LastName)
	assert.Equal(t, "NO", r.WatchListed)
	assert.Equal(t, "06-May-1996", r.DateOfBirth)
	assert.Nil(t, err)
}

//StartServer initializes a test HTTP server useful for request mocking, Integration tests and Client configuration
func testServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		switch r.URL.Path {
		case "/account/auth":
			successBody := fmt.Sprintf(`{"id": "%v"}`, testAccountId)
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case fmt.Sprintf("/accounts/%v", testAccountId):
			body := fmt.Sprintf(`{
    "meta": {"data_status": "AVAILABLE"},
    "account": {
        "_id": "%v",
        "institution": {
            "name": "Access Bank",
            "bankCode": "044",
            "type": "PERSONAL_BANKING"
        },
        "name": "IDORENYIN OBONG OBONG",
        "currency": "NGN",
        "type": "Current",
        "accountNumber": "0788164862",
        "balance": 37836709,
        "bvn": "6800"
    }
}`, testAccountId)
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case fmt.Sprintf("/accounts/%v/statement", testAccountId):
			body := ""
			switch r.URL.Query().Get("output") {
			case "json":
				body = `{"meta": {"count": 2},
    "data": [
        {
            "_id": "5fc686f98b97632dbef0f8db",
            "type": "debit",
            "date": "2020-12-01T00:00:00.000Z",
            "narration": "VALUE ADDED TAX VAT ON NIP TRANSFER FOR Yusuf Money TO FCMB/OGUNGBEFUN OLADUNNI KHADIJAH ReF:",
            "amount": 375,
            "balance": 10517116
        },
		{
            "_id": "5fc686f98b97632dbef0f8dc",
            "type": "debit",
            "date": "2020-12-01T00:00:00.000Z",
            "narration": "COMMISSION NIP TRANSFER COMMISSION FOR Yusuf Money TO FCMB/OGUNGBEFUN OLADUNNI KHADIJAH ReF:",
            "amount": 5000,
            "balance": 10517491
        }
]
}`
			case "pdf":
				body = fmt.Sprintf(`{"id": "%v", "status": "BUILDING", "path": "https://api.withmono.com/statements/pvLhFR89Id2zrnPGJZcM.pdf"}`, testJobId)
			}
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case fmt.Sprintf("/accounts/%v/statement/jobs/%v", testAccountId, testJobId):
			body := fmt.Sprintf(`{"id": "%v", "status": "COMPLETE", "path": "https://api.withmono.com/statements/pvLhFR89Id2zrnPGJZcM.pdf"}`, testJobId)
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case fmt.Sprintf("/accounts/%v/transactions", testAccountId):
			body := `{
  "paging": {
    "total": 190,
    "page": 2,
    "previous": "https://api.withmono.com/accounts/:id/transactions?page=2",
    "next": "https://api.withmono.com/accounts/:id/transactions?page=3"
  },
  "data": [
    {
      	"_id": "5f171a540295e231abca1154",
      	"amount": 10000,
      	"date": "2020-07-21T00:00:00.000Z",
      	"narration": "TRANSFER from HASSAN ABDULHAMID TOMIWA to OGUNGBEFUN OLADUNNI KHADIJAH",
      	"type": "debit",
      	"category": "E-CHANNELS"
    },
    {
		"_id": "5d171a540295e231abca6654",
      	"amount": 20000,
      	"date": "2020-07-21T00:00:00.000Z",
      	"narration": "TRANSFER from HASSAN ABDULHAMID TOMIWA to UMAR ABDULLAHI",
      	"type": "debit",
      	"category": "E-CHANNELS"
    }
  ]
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case fmt.Sprintf("/accounts/%v/credit", testAccountId):
			body := `{"total": 2000000,
  "history": [
    {
      "amount": 1000000,
      "period": "01-20"
    },
    {
      "amount": 1000000,
      "period": "07-20"
    }
  ]
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case fmt.Sprintf("/accounts/%v/debit", testAccountId):
			body := `{"total": 1000000,
  "history": [
    {
      "amount": 1000000,
      "period": "01-20"
    },
    {
      "amount": 1000000,
      "period": "07-20"
    }
  ]
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case fmt.Sprintf("/accounts/%v/income", testAccountId):
			body := `{"type": "INCOME", "amount": 59700000, "employer": "Relentless Labs Inc", "confidence": 0.95}`
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case fmt.Sprintf("/accounts/%v/identity", testAccountId):
			body := `{
    "firstName": "ABDULHAMID",
    "middleName": "TOMIWA",
    "lastName": "HASSAN",
    "dateOfBirth": "06-May-1996",
    "phoneNumber1": "0000000",
    "phoneNumber2": "",
    "registrationDate": "26-Mar-2018",
    "email": "tomiwa.jr@gmail.com",
    "gender": "Male",
    "levelOfAccount": "Level 1 - Low Level Accounts",
    "lgaOfOrigin": "Abeokuta South",
    "lgaOfResidence": "Ikeja",
    "maritalStatus": "Single",
    "nin": "000000",
    "nationality": "Nigeria",
    "residentialAddress": "23, SHITTU ANIMASHA STR, PHASE2, GBAGADA",
    "stateOfOrigin": "Ogun State",
    "stateOfResidence": "Lagos State",
    "title": "Mr",
    "watchListed": "NO",
    "bvn": "00000000",
    "photoId": "/9j/4AAQSkZJRgABAgAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAGQASwDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDA1JO45APTiooyA4Bzk+lMRpWiELk1Bqsm2IjP4Vbtj8mcfnWPrE3O3PNIOpkR8z1sRBhFjNY8ESSyfOM4ORWwmNg5pIpkZYMWUhht9RwfpVWUBckVoMMrVKVccdfrTYBbykNtBAPvVq6TfF0BxVKA4frzWkDuXHYigRglSjn61Mk20jd0ovEKOSaroQ/Q0ii8xjnUZXPoR1FVJLeQfcbd7Hqb0ak86MyGMMC4AJXuB/kUhjwOKORxTDncDn5R1FPznmgBe1RlVkKvlhtPr1p2TnHWjoDxTAjIwx+tSIMLz1qMkk07NAH//2Q=="
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		case "/coverage":
			successBody := `[
    {
        "name": "GTBank",
        "icon": "https://connect.withmono.com/build/img/guaranty-trust-bank.png",
        "website": "https://www.gtbank.com",
        "coverage": {
            "personal": true,
            "business": true,
            "countries": [
                "NG"
            ]
        },
        "products": [
            "Auth",
            "Accounts",
            "Transactions",
            "Balance",
            "Income",
            "Identity",
            "Direct Debit"
        ]
    },
    {
        "name": "Access Bank",
        "icon": "https://connect.withmono.com/build/img/access-bank.png",
        "website": "https://www.accessbankplc.com",
        "coverage": {
            "personal": true,
            "business": true,
            "countries": [
                "NG"
            ]
        },
        "products": [
            "Auth",
            "Accounts",
            "Transactions",
            "Balance",
            "Income",
            "Identity",
            "Direct Debit"
        ]
    },
    {
        "name": "First Bank",
        "icon": "https://connect.withmono.com/build/img/first-bank-of-nigeria.png",
        "website": "https://www.firstbanknigeria.com",
        "coverage": {
            "personal": true,
            "business": true,
            "countries": [
                "NG"
            ]
        },
        "products": [
            "Auth",
            "Accounts",
            "Transactions",
            "Balance",
            "Income",
            "Identity",
            "Direct Debit"
        ]
    },
    {
        "name": "Fidelity Bank",
        "icon": "https://connect.withmono.com/build/img/fidelity-bank.png",
        "website": "https://www.fidelitybank.ng",
        "coverage": {
            "personal": true,
            "business": true,
            "countries": [
                "NG"
            ]
        },
        "products": [
            "Auth",
            "Accounts",
            "Transactions",
            "Balance",
            "Income",
            "Identity",
            "Direct Debit"
        ]
    }]`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case "/v1/lookup/bvn/identity":
			body := `{
    "firstName": "ABDULHAMID",
    "middleName": "TOMIWA",
    "lastName": "HASSAN",
    "dateOfBirth": "06-May-1996",
    "phoneNumber1": "0000000",
    "phoneNumber2": "",
    "registrationDate": "26-Mar-2018",
    "email": "tomiwa.jr@gmail.com",
    "gender": "Male",
    "levelOfAccount": "Level 1 - Low Level Accounts",
    "lgaOfOrigin": "Abeokuta South",
    "lgaOfResidence": "Ikeja",
    "maritalStatus": "Single",
    "nin": "000000",
    "nationality": "Nigeria",
    "residentialAddress": "23, SHITTU ANIMASHA STR, PHASE2, GBAGADA",
    "stateOfOrigin": "Ogun State",
    "stateOfResidence": "Lagos State",
    "title": "Mr",
    "watchListed": "NO",
    "bvn": "00000000",
    "photoId": "/9j/4AAQSkZJRgABAgAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAGQASwDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDA1JO45APTiooyA4Bzk+lMRpWiELk1Bqsm2IjP4Vbtj8mcfnWPrE3O3PNIOpkR8z1sRBhFjNY8ESSyfOM4ORWwmNg5pIpkZYMWUhht9RwfpVWUBckVoMMrVKVccdfrTYBbykNtBAPvVq6TfF0BxVKA4frzWkDuXHYigRglSjn61Mk20jd0ovEKOSaroQ/Q0ii8xjnUZXPoR1FVJLeQfcbd7Hqb0ak86MyGMMC4AJXuB/kUhjwOKORxTDncDn5R1FPznmgBe1RlVkKvlhtPr1p2TnHWjoDxTAjIwx+tSIMLz1qMkk07NAH//2Q=="
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, body)

		default:
			w.WriteHeader(500)
		}
	}))
	return server
}
