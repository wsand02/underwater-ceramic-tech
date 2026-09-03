package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "https://platsbanken-api.arbetsformedlingen.se/jobs/v1/search"
const occupationGroup = "DJh5_yyF_hEM"
const region = "CifL_Rzy_Mku"
const maxRecords = 25

type JobListing struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Occupation string `json:"occupation"`
	Company    string `json:"company"`
	Published  bool   `json:"published"`
	DatePosted string `json:"publishedDate"`
}

type SearchRequest struct {
	Filters []Filter `json:"filters"`
	Size    int      `json:"size"`
}

type Filter struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SearchResponse struct {
	Ads []JobListing `json:"ads"`
}

func main() {
	body, err := json.Marshal(SearchRequest{
		Filters: []Filter{
			{Type: "occupationGroup", Value: occupationGroup},
			{Type: "region", Value: region},
		},
		Size: maxRecords,
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		responseBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		panic(fmt.Sprintf("search failed: %s: %s", resp.Status, responseBody))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		panic(err)
	}
	fmt.Println("Title\tOccupation\tCompany\tDate Posted")
	for _, job := range searchResponse.Ads {
		// req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://platsbanken-api.arbetsformedlingen.se/jobs/v1/job/%s", job.ID), nil)
		// if err != nil {
		// 	panic(err)
		// }
		// req.Header.Set("Content-Type", "application/json")

		// resp, err := http.DefaultClient.Do(req)
		// if err != nil {
		// 	panic(err)
		// }
		// defer resp.Body.Close()

		// if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		// 	responseBody, readErr := io.ReadAll(resp.Body)
		// 	if readErr != nil {
		// 		panic(readErr)
		// 	}
		// 	panic(fmt.Sprintf("job details request failed: %s: %s", resp.Status, responseBody))
		// }

		// var jobDetails JobListing
		// if err := json.NewDecoder(resp.Body).Decode(&jobDetails); err != nil {
		// 	panic(err)
		// }

		fmt.Printf("%s\t%s\t%s\t%s\n", job.Title, job.Occupation, job.Company, job.DatePosted)
		//time.Sleep(1 * time.Second) // Sleep for 1 second to avoid overwhelming the API
	}
}
