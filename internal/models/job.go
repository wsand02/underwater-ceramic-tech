package models

import (
	"fmt"
	"time"

	"github.com/wsand02/underwater-ceramic-tech/internal/database"
)

// ID
// Title
// Occupation
// Company Name
// Date Posted
// Description
// URL
// Reference

const jobSchema = `
CREATE TABLE IF NOT EXISTS jobs (
	id bigint PRIMARY KEY,
	title text NOT NULL,
	occupation text,
	company text,
	date_posted timestamptz,
	description text,
	url text,
	reference text
)
`

type Job struct {
	ID          int64
	Title       string
	Occupation  string
	Company     string
	DatePosted  *time.Time
	Description string
	URL         string
	Reference   string
}

func (j *Job) Migrate() error {
	db, err := database.GetInstance()
	if err != nil {
		return fmt.Errorf("Job.Migrate: %v", err)
	}
	_, err = db.Exec(jobSchema)
	if err != nil {
		return fmt.Errorf("Job.Migrate: %v", err)
	}
	return nil
}
