package models

type Schema interface {
	Migrate() error
}
