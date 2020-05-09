package datastore

import (
	"encoding/json"
	"time"
)

import "github.com/allegro/bigcache"
type LocalAccountCache struct {
	Cache bigcache.BigCache
	LoadSource DataStore
}

//We can Tune Bigcache parameters like Number of shards and cache size etc. to make it more efficient in handling concurrent requests.
// In this cache default config of 1024 shards should be good enough, otherwise we can create a config struct and pass all these params via constructor
func NewLocalCache(loadSource DataStore) LocalAccountCache {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	return LocalAccountCache{Cache:*cache, LoadSource:loadSource}
}

func (l LocalAccountCache) GetAccountDetailsFromLicenseKey(licenseKey string) (*Account, error) {
	var account Account
	accountBytes, err := l.Cache.Get(licenseKey)
	if err != nil {
		account, err := l.LoadSource.GetAccountDetailsFromLicenseKey(licenseKey)
		if err != nil {
			return nil, err
		}
		accountBytesFromLoader, _ := json.Marshal(account)
		l.Cache.Set(licenseKey, accountBytesFromLoader)
		return account, nil
	}
	err = json.Unmarshal(accountBytes,&account)
	return &account, err
}