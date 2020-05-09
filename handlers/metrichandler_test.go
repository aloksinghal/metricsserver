package handlers

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"metricsserver/datastore"
	"metricsserver/publish"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_validateTimestamp(t *testing.T) {
	valid, err:=  validateTimestamp(1589037264461, 900)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, true, valid)
}

func TestMetricHandler_ProcessMetricsWOLicenseKey(t *testing.T) {
	publisher := publish.MockPublisher{}
	dataStore := datastore.Mockstore{
		Account:datastore.Account{
			AccountId:"abc",
			LicenseKey:"def",
			IsValid:true,
		},
		Err:nil,
	}
	metricHandler := NewMetricHandler(dataStore,publisher,900, 10000)

	bytes, _ := ioutil.ReadFile("testData.json")
	req := httptest.NewRequest("POST", "http://example.com/foo",strings.NewReader(string(bytes)))
	w := httptest.NewRecorder()
	metricHandler.ProcessMetrics(w, req)
	resp := w.Result()
	if resp.StatusCode != 401 {
		t.Error("request without license key passed")
	}
}

func TestMetricHandler_ProcessMetricsWithInvalidLicenseKey(t *testing.T) {
	publisher := publish.MockPublisher{}
	dataStore := datastore.Mockstore{
		Account:datastore.Account{
			AccountId:"abc",
			LicenseKey:"def",
			IsValid:true,
		},
		Err:nil,
	}
	metricHandler := NewMetricHandler(dataStore,publisher,900, 10000)

	bytes, _ := ioutil.ReadFile("testData.json")
	req := httptest.NewRequest("POST", "http://example.com/foo",strings.NewReader(string(bytes)))
	req.Header.Set("X-License-Key", "fgh")
	w := httptest.NewRecorder()
	metricHandler.ProcessMetrics(w, req)
	resp := w.Result()
	if resp.StatusCode != 401 {
		t.Error("request with invalid license key passed")
	}
}

func TestMetricHandler_ProcessMetricsWithValidLicenseKey(t *testing.T) {
	publisher := publish.MockPublisher{}
	dataStore := datastore.Mockstore{
		Account:datastore.Account{
			AccountId:"abc",
			LicenseKey:"def",
			IsValid:true,
		},
		Err:nil,
	}
	metricHandler := NewMetricHandler(dataStore,publisher,900, 1000)

	bytes, _ := ioutil.ReadFile("testData.json")
	req := httptest.NewRequest("POST", "http://example.com/foo",strings.NewReader(string(bytes)))
	req.Header.Set("X-License-Key", "def")
	w := httptest.NewRecorder()
	metricHandler.ProcessMetrics(w, req)
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Error("request with invalid license key passed")
	}
}


func TestMetricHandler_ProcessMetricsWithDrift(t *testing.T) {
	publisher := publish.MockPublisher{}
	dataStore := datastore.Mockstore{
		Account:datastore.Account{
			AccountId:"abc",
			LicenseKey:"def",
			IsValid:true,
		},
		Err:nil,
	}
	metricHandler := NewMetricHandler(dataStore,publisher,900, -1)

	bytes, _ := ioutil.ReadFile("testData.json")
	req := httptest.NewRequest("POST", "http://example.com/foo",strings.NewReader(string(bytes)))
	req.Header.Set("X-License-Key", "def")
	w := httptest.NewRecorder()
	metricHandler.ProcessMetrics(w, req)
	resp := w.Result()
	if resp.StatusCode != 500 {
		t.Error("request with more than allowed drift passed")
	}
}




