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

// måste stödja att annonser blir borttagna från af
// vissa annonser saknar ju url för ansökan pga ansök via mail etc etc måste stödja det, troligtvis via en AFUrl metod där jag konstruerar url till af bara.

const jobSchema = `
CREATE TABLE IF NOT EXISTS jobs (
	id bigint PRIMARY KEY,
	title text NOT NULL,
	occupation text,
	company text,
	date_posted timestamptz,
	description text,
	url text,
	reference text,
	created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
)
`

type Job struct {
	ID          int64
	Title       string
	Occupation  string
	Company     string
	DatePosted  time.Time
	Description string
	URL         string
	Reference   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

func (j *Job) Upsert() error {
	db, err := database.GetInstance()
	if err != nil {
		return fmt.Errorf("Job.Upsert: %w", err)
	}

	_, err = db.Exec(`
		INSERT INTO jobs (
			id,
			title,
			occupation,
			company,
			date_posted,
			description,
			url,
			reference
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			occupation = EXCLUDED.occupation,
			company = EXCLUDED.company,
			date_posted = EXCLUDED.date_posted,
			description = EXCLUDED.description,
			url = EXCLUDED.url,
			reference = EXCLUDED.reference,
			updated_at = NOW()
	`,
		j.ID,
		j.Title,
		j.Occupation,
		j.Company,
		j.DatePosted,
		j.Description,
		j.URL,
		j.Reference,
	)

	if err != nil {
		return fmt.Errorf("Job.Upsert: %w", err)
	}

	return nil
}
