package models

// ID
// Title
// Occupation
// Company Name
// Date Posted
// Description
// URL
// Reference

const JobSchema = `
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
