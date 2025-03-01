package images

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// PexelsAPIKey is the API key for Pexels
// In production, this should be loaded from environment variables
var PexelsAPIKey = os.Getenv("PEXELS_API_KEY")

// PexelsResponse represents the response from Pexels API
type PexelsResponse struct {
	Page         int     `json:"page"`
	PerPage      int     `json:"per_page"`
	Photos       []Photo `json:"photos"`
	TotalResults int     `json:"total_results"`
	NextPage     string  `json:"next_page"`
}

// Photo represents a photo from Pexels API
type Photo struct {
	ID              int         `json:"id"`
	Width           int         `json:"width"`
	Height          int         `json:"height"`
	URL             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerURL string      `json:"photographer_url"`
	PhotographerID  int         `json:"photographer_id"`
	AvgColor        string      `json:"avg_color"`
	Src             PhotoSource `json:"src"`
	Liked           bool        `json:"liked"`
	Alt             string      `json:"alt"`
}

// PhotoSource represents the different sizes of a photo
type PhotoSource struct {
	Original  string `json:"original"`
	Large2x   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

// ImageResult represents the result of an image search
type ImageResult struct {
	URL             string `json:"url"`
	Photographer    string `json:"photographer"`
	PhotographerURL string `json:"photographer_url"`
	Alt             string `json:"alt"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
}

// GetTravelImage fetches a travel image for a destination from Pexels API
func GetTravelImage() string {
	defaultImage := "https://images.pexels.com/photos/2245436/pexels-photo-2245436.png?auto=compress&cs=tinysrgb&h=650&w=940"

	// Ensure API key is set
	if PexelsAPIKey == "" {
		return defaultImage
	}

	// Create the request
	url := "https://api.pexels.com/v1/search?query=travel&per_page=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return defaultImage
	}

	// Add authorization header
	req.Header.Add("Authorization", PexelsAPIKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return defaultImage
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return defaultImage
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultImage
	}

	// Parse the response
	var pexelsResp PexelsResponse
	if err := json.Unmarshal(body, &pexelsResp); err != nil {
		return defaultImage
	}

	// Check if we got any photos
	if len(pexelsResp.Photos) == 0 {
		return defaultImage
	}

	// Get the first photo and return its large URL
	return pexelsResp.Photos[0].Src.Large
}
