package datastore


type Account struct {
	LicenseKey string `json:"licenseKey"`
	AccountId string  `json:"accountId"`
	IsValid bool `json:"isValid"`
}

func (a Account) IsAccountValid() bool {
	return a.IsValid
}