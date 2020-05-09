package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"metricsserver/datastore"
	"metricsserver/publish"
	"net/http"
	"time"
)

type MetricHandler struct {
	store datastore.DataStore
	publisher publish.Publisher
	maxPermissiblePayload int
	allowableOffset int64
}

type Metric struct {
	AccountId string
	UnixTimeEpochMs int64
	Metric string
	Value string
}

func NewMetricHandler(store datastore.DataStore, publisher publish.Publisher, maxPermissiblePayload int, allowableOffset int64) MetricHandler {
	return MetricHandler{store:store , publisher:publisher, maxPermissiblePayload:maxPermissiblePayload, allowableOffset:allowableOffset}
}

func (m MetricHandler) ProcessMetrics(wr http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(wr,"Method not allowed", 405)
	}

	licenseKey := r.Header.Get("X-License-Key")
	if licenseKey == "" {
		http.Error(wr,errors.New("missing license key").Error() , 401)
		return
	}

	account, err := m.store.GetAccountDetailsFromLicenseKey(licenseKey)
	if err != nil {
		http.Error(wr, err.Error(), 401)
		return
	}
	if !account.IsAccountValid() {
		http.Error(wr, errors.New("account not valid").Error(), 401)
		return
	}

	timestamp := makeTimestamp()
	_ , err = validateTimestamp(timestamp, m.allowableOffset)
	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}
	if len(body) > m.maxPermissiblePayload {
		http.Error(wr, "data size exceeded max permissible limit", 500)
		return
	}

	requestData := map[string]string{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}

	go publishMessages(account.AccountId, timestamp, requestData, m.publisher)
	wr.Write([]byte("successfully processed metric"))
}

func validateTimestamp(timestamp int64, allowableOffsetInMs int64) (bool, error){
	if (timestamp/1000) % 2 == 0 {
		if (timestamp % 1000) > allowableOffsetInMs {
			return false, errors.New("jitter higher than offset")
		}
	} else {
		if (timestamp % 1000) < (1000 - allowableOffsetInMs) {
			return false, errors.New("jitter higher than offset")
		}
	}
	return true, nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func publishMessages(accountId string, timeStamp int64, data map[string]string, publisher publish.Publisher) {
	metrics := make([][]byte, len(data))
	for key, value := range data {
		metric := Metric{
			AccountId: accountId,
			UnixTimeEpochMs:timeStamp,
			Metric:key,
			Value:value,
		}
		byteData, err := json.Marshal(metric)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		metrics = append(metrics, byteData)
	}
	err := publisher.Publish(metrics, accountId)
	if err != nil {
		fmt.Println(err)
	}
}