# GoMono
GoMono is a Golang API wrapper around Mono's REST API. 

It implements the `Authentication`, `Account`, `User` and `Misc` endpoints as documented on the [Mono's v1 API docs](https://docs.mono.co/reference)

## Install
```
$ go get github.com/jcobhams/gomono
```

## Usage

```go
package main

import (
	"github.com/jcobhams/gomono"
	"net/http"
)

func main() {
	gm, err := gomono.New(gomono.NewDefaultConfig("YOUR_SECRET_KEY")) // If you'd like to use the default config

	//Or you could configure the HTTP Client how you'd like (timeouts, Transports etc
	//gm, err := gomono.New(gomono.Config{
	//	SecretKey:  "YOUR_SECRET_KEY",
	//	HttpClient: &http.Client{
	//		Transport:     nil,
	//		CheckRedirect: nil,
	//		Jar:           nil,
	//		Timeout:       0,
	//	},
	//	ApiUrl:     "https://api.withmono.com",
	//})
	
	if err != nil {
		//do something with error
    }
    
    //Exchange Token
    id, err := gm.ExchangeToken("CODE")
    
    //Get Bank Account Details
    infResponse, err := gm.Information(id)
    
    //Get Bank Statement - Period => in months (1-12) | output => json or pdf
    stmtResponse, err := gm.Statement(id, "period", "output")
    
    // Query PDF Job Status - If the Statement call above had a pdf output, the response will contain a PDF struct with and ID
    pdfStmtResponse, err := gm.PdfStatementJobStatus(id, stmtResponse.PDF.ID)
    
    // Get user transactions - tnxType => debit or credit | paginate => bool
    tnxResponse, err := gm.Transactions(id, "start", "end", "narration", "tnxType", "paginate")

    // Get Credit Transactions
    crdTnxResponse, err := gm.CreditTransactions(id)
    
    // Get Debit Transactions
    dbtTnxResponse, err := gm.DebitTransactions(id)
    
    // Get Income Information
    incResponse, err := gm.Income(id)
    
    // Get Identity Information
    idyResponse, err := gm.Identity(id)
    
    // Get Institutions List
    insResponse, err := gm.Institutions()
    
    // LookupBVN
    bvnResponse, err := gm.LookupBVN("1234567890")

}

```

In all the examples, error handling has been ignore/suppressed. Please handle errors properly to avoid `nil pointer` panics.

## Integration Testing
`Gomono` is an interface that can easily be mocked to ease testing.

You could also use the explicit configuration option shown earlier to create your clients. 

That way you can set a test API Url or intercept HTTP calls using a fake http client - Whatever works best for you :)

## What's Not Covered/Supported
1. Data Sync / Reauth Code

## Run Tests
go test -race -v -coverprofile cover.out

## View Coverage
go tool cover -html=cover.out