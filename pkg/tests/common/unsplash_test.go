package common

import (
	"bytes"
	"github.com/Ayocodes24/GO-Eats/pkg/service/restaurant/unsplash"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (mc *MockClient) Do(req *http.Request) (*http.Response, error) {
	return mc.DoFunc(req)
}

func TestGetUnSplashImageURL(t *testing.T) {
	mockResponse := `{
		"results": [{
			"urls": {
				"small": "https://example.com/image.jpg"
			}
		}]
	}`
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
				Header:     make(http.Header),
			}, nil
		},
	}

	imageURL := unsplash.GetUnSplashImageURL(mockClient, "test")
	expectedURL := "https://example.com/image.jpg"
	if imageURL != expectedURL {
		t.Fatalf("expected %s, got %s", expectedURL, imageURL)
	}
}
