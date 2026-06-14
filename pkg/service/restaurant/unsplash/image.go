package unsplash

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func GetUnSplashImageURL(client HttpClient, menuItem string) (string, error) {
	imageUrl := "https://api.unsplash.com/search/photos/?page=1&query=" + url.QueryEscape(menuItem) + "&w=400&h=400"
	req, err := http.NewRequest("GET", imageUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Client-ID", os.Getenv("UNSPLASH_API_KEY")))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse UnSplash
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(apiResponse.Results) == 0 {
		return "", fmt.Errorf("no image results found for %q", menuItem)
	}

	return apiResponse.Results[0].Urls.Small, nil
}
