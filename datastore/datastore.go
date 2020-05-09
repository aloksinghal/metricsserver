package datastore

type DataStore interface {
	GetAccountDetailsFromLicenseKey(licenseKey string) (*Account, error)
}
