package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const searchURL = "https://platsbanken-api.arbetsformedlingen.se/jobs/v1/search"
const occupationGroup = "DJh5_yyF_hEM"
const region = "CifL_Rzy_Mku"
const maxRecords = 25

const jobDetailsURL = "https://platsbanken-api.arbetsformedlingen.se/jobs/v1/job/"

type JobListing struct {
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
	Ads []JobListing `json:"ads"`
}

type JobDetails struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Occupation  string `json:"occupation"`
	Application struct {
		Reference string `json:"reference"`
		URL       string `json:"url"`
	} `json:"application"`
	Description string `json:"description"`
}

func jobSearch(occupationGroup, region string, maxRecords int, startIndex int) ([]JobListing, error) {
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

func jobDetails(jobID string) (*JobDetails, error) {
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

	var jobDetails JobDetails
	if err := json.NewDecoder(resp.Body).Decode(&jobDetails); err != nil {
		return nil, err
	}
	return &jobDetails, nil
}

func main() {
	ads, err := jobSearch(occupationGroup, region, maxRecords, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("Title\tOccupation\tCompany\tDate Posted")
	for _, job := range ads {
		fmt.Printf("%s\t%s\t%s\t%s\n", job.Title, job.Occupation, job.Company, job.DatePosted)
		details, err := jobDetails(job.ID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Details for %s:\n", job.Title)
		fmt.Printf("  ID: %s\n", details.ID)
		fmt.Printf("  Description: %s\n", details.Description)
		time.Sleep(1 * time.Second)
	}
}
