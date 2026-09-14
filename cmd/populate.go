/*
Copyright © 2026 William Sandbrink
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wsand02/underwater-ceramic-tech/internal/client"
	"github.com/wsand02/underwater-ceramic-tech/internal/models"
)

// TODO: decide on orm library
// TODO: create models

const occupationGroup = "DJh5_yyF_hEM"
const region = "CifL_Rzy_Mku"
const maxRecords = 25

// populateCmd represents the populate command
var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: populate,
}

func populate(cmd *cobra.Command, args []string) {
	ads, err := client.JobSearch(occupationGroup, region, maxRecords, 0)
	if err != nil {
		fmt.Println("Error searching for jobs:", err)
		return
	}

	for _, ad := range ads {
		details, err := client.JobDetails(ad.ID)
		if err != nil {
			fmt.Printf("Error fetching details for job ID %s: %v\n", ad.ID, err)
			continue
		}
		id, err := strconv.ParseInt(details.ID, 10, 64)
		if err != nil {
			fmt.Printf("invalid job ID %q: %v", details.ID, err)
			continue
		}
		job := &models.Job{
			ID:          id,
			Title:       details.Title,
			Occupation:  details.Occupation,
			Company:     ad.Company,
			DatePosted:  ad.DatePosted,
			Description: details.Description,
			URL:         details.Application.URL,
			Reference:   details.Application.Reference,
		}
		err = job.Upsert()
		if err != nil {
			fmt.Printf("Error upserting job ID: %s: %v\n", ad.ID, err)
		}
		fmt.Printf("Job ID: %v\nTitle: %s\nOccupation: %s\nCompany: %s\nPublished: %t\nDate Posted: %s\nDescription: %s\nApplication Reference: %s\nApplication URL: %s\n\n",
			details.ID, details.Title, details.Occupation, ad.Company, ad.Published, ad.DatePosted, details.Description, details.Application.Reference, details.Application.URL)
	}
}

func init() {
	rootCmd.AddCommand(populateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// populateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// populateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
