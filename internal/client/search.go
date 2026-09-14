package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JobListingResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Occupation string `json:"occupation"`
	Company    string `json:"workplaceName"`
	Published  bool   `json:"published"`
	DatePosted string `json:"publishedDate"`
}

type SearchRequest struct {
	Filters []Filter `json:"filters"`
	Size    int      `json:"maxRecords"`
	Start   int      `json:"startIndex"`
}

type Filter struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SearchResponse struct {
	Ads []JobListingResponse `json:"ads"`
}

type JobDetailsResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Occupation  string `json:"occupation"`
	Application struct {
		Reference string `json:"reference"`
		URL       string `json:"webAddress"`
	} `json:"application"`
	Description string `json:"description"`
}

const searchURL = "https://platsbanken-api.arbetsformedlingen.se/jobs/v1/search"

func JobSearch(occupationGroup, region string, maxRecords int, startIndex int) ([]JobListingResponse, error) {
	body, err := json.Marshal(SearchRequest{
		Filters: []Filter{
			{Type: "occupationGroup", Value: occupationGroup},
			{Type: "region", Value: region},
		},
		Size:  maxRecords,
		Start: startIndex,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, searchURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		responseBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		return nil, fmt.Errorf("search failed: %s: %s", resp.Status, responseBody)
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	return searchResponse.Ads, nil
}

const jobDetailsURL = "https://platsbanken-api.arbetsformedlingen.se/jobs/v1/job/"

func JobDetails(jobID string) (*JobDetailsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", jobDetailsURL, jobID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		responseBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		return nil, fmt.Errorf("job details request failed: %s: %s", resp.Status, responseBody)
	}

	var jobDetails JobDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobDetails); err != nil {
		return nil, err
	}
	return &jobDetails, nil
}
