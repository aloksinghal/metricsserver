package datastore

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type SqlStore struct {
	Connection sql.DB
}

//We should pass the host and other configs via constructor. We can use max open connections and Max Idle connections param to ensure that we
//don't overwhelm the database due to cache misses or extra load
func NewSqlStore(maxIdleConnections int, maxOpenConnections int) DataStore {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/metrics")
	db.SetConnMaxLifetime(1000 * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
	if err != nil {
		panic(err)
	}
	return SqlStore{Connection: *db}
}

func (s SqlStore) GetAccountDetailsFromLicenseKey(licenseKey string) (*Account, error) {
	var account Account
	err := s.Connection.QueryRow("SELECT licenseKey,accountId,isValid FROM accounts where licenseKey= ?", licenseKey).
		Scan(&account.LicenseKey, &account.AccountId, &account.IsValid)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
