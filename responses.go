// Copyright 2020 Joseph Cobhams. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE.md file.
//
package gomono

type (
	InformationResponse struct {
		Meta struct {
			DataStatus string `json:"data_status"`
		} `json:"meta"`
		Account struct {
			ID            string  `json:"_id"`
			Name          string  `json:"name"`
			Currency      string  `json:"currency"`
			Type          string  `json:"type"`
			AccountNumber string  `json:"accountNumber"`
			Balance       float64 `json:"balance"`
			BVN           string  `json:"bvn"`
			Institution   struct {
				Name     string `json:"name"`
				BankCode string `json:"bankCode"`
				Type     string `json:"type"`
			} `json:"institution"`
		}
	}

	StatementResponse struct {
		JSON *StatementResponseJson
		PDF  *StatementResponsePdf
	}

	StatementResponseJson struct {
		Meta struct{ Count int } `json:"meta"`
		Data []struct {
			ID        string  `json:"_id"`
			Type      string  `json:"type"`
			Date      string  `json:"date"`
			Narration string  `json:"narration"`
			Amount    float64 `json:"amount"`
			Balance   float64 `json:"balance"`
		}
	}

	StatementResponsePdf struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Path   string `json:"path"`
	}

	TransactionsResponse struct {
		Paging struct {
			Total    int    `json:"total"`
			Page     int    `json:"page"`
			Previous string `json:"previous"`
			Next     string `json:"next"`
		}
		Data []struct {
			ID        string  `json:"_id"`
			Amount    float64 `json:"amount"`
			Date      string  `json:"date"`
			Narration string  `json:"narration"`
			Type      string  `json:"type"`
			Category  string  `json:"category"`
		}
	}

	TransactionByTypeResponse struct {
		Total   float64 `json:"total"`
		History []struct {
			Amount float64 `json:"amount"`
			Period string  `json:"period"`
		}
	}

	IncomeResponse struct {
		Type       string  `json:"type"`
		Amount     float64 `json:"amount"`
		Employer   string  `json:"employer"`
		Confidence float64 `json:"confidence"`
	}

	IdentityResponse struct {
		FirstName          string `json:"firstName"`
		MiddleName         string `json:"middleName"`
		LastName           string `json:"lastName"`
		DateOfBirth        string `json:"dateOfBirth"`
		PhoneNumber1       string `json:"phoneNumber1"`
		PhoneNumber2       string `json:"phoneNumber2"`
		RegistrationDate   string `json:"registrationDate"`
		Email              string `json:"email"`
		Gender             string `json:"gender"`
		LevelOfAccount     string `json:"levelOfAccount"`
		LgaOfOrigin        string `json:"lgaOfOrigin"`
		LgaOfResidence     string `json:"lgaOfResidence"`
		MaritalStatus      string `json:"maritalStatus"`
		NIN                string `json:"nin"`
		Nationality        string `json:"nationality"`
		ResidentialAddress string `json:"residentialAddress"`
		StateOfOrigin      string `json:"stateOfOrigin"`
		StateOfResidence   string `json:"stateOfResidence"`
		Title              string `json:"title"`
		WatchListed        string `json:"watchListed"`
		BVN                string `json:"bvn"`
		PhotoID            string `json:"photo_id"`
	}

	InstitutionsResponse struct {
		Institutions []Institution
	}

	Institution struct {
		Name     string `json:"name"`
		Icon     string `json:"icon"`
		Website  string `json:"website"`
		Coverage struct {
			Personal  bool     `json:"personal"`
			Business  bool     `json:"business"`
			Countries []string `json:"countries"`
		} `json:"coverage"`
		Products []string `json:"products"`
	}
)
