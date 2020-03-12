package main_test

import (
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"testing"
	"time"
	"app"
	"app/mocks"
)

func TestCollectFunc(t *testing.T) {
	t.Log("Start TestCollectFunc")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	httpClient := mocks.NewMockHttpClient(mockCtrl)
	pcCollector := main.PrometheusCollector{
		"testName", "testHelpMsg",
		prometheus.Labels{"url": "testUrl"},
		"testUrl",
		false,
		httpClient,
	}

	// Case 1: Up/Down status check and http status code is 200
	httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusOK,}, nil).Times(1)

	result := pcCollector.CollectFunc()
	if result != 1 {
		t.Errorf("Return value of CollectFunc() should be 1, but is is %f", result)
	}

	// Case 2: Up/Down status check and http status code is NOT 200
	httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusServiceUnavailable,}, nil).Times(1)
	result = pcCollector.CollectFunc()
	if result != 0 {
		t.Errorf("Return value of CollectFunc() should be 0, but is is %f", result)
	}

	// Case 2: Simulate the condition that resp time is 1000 milliseconds
	// So the resp time should between 1000 and 1100 milliseconds(1000 milliseconds sleep time introduced + code runtime)
	pcCollector.IsRespTimeCollector = true
	httpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func (req *http.Request) (*http.Response, error) {
		time.Sleep(time.Millisecond * 1000)
		return &http.Response{StatusCode: http.StatusServiceUnavailable,}, nil
	}).Times(1)
	result = pcCollector.CollectFunc()
	if result < 1000 || result > 1100 {
		t.Errorf("Expected response time should be around 1000 ms, but it is %f ms", result)
	}

	t.Log("End TestCollectFunc")
}
