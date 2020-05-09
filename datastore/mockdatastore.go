package datastore

import "errors"

type Mockstore struct {
	Account Account
	Err error
}


func (m Mockstore) GetAccountDetailsFromLicenseKey(licenseKey string) (*Account, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Account.LicenseKey != licenseKey {
		return nil, errors.New("license key does not exist")
	}
	return &m.Account, nil
}