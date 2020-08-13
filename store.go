package pirsch

import (
	"database/sql"
	"time"
)

// NullTenant can be used to pass no (null) tenant to filters and functions.
// This is an invalid sql.NullInt64 with a value of 0.
var NullTenant = NewTenantID(0)

// Store defines an interface to persists hits and other data.
// The first parameter (if required) is always the tenant ID and can be left out (pirsch.NullTenant), if you don't want to split your data.
// This is usually the case if you integrate Pirsch into your application.
type Store interface {
	// Save persists a list of hits.
	Save([]Hit) error

	// DeleteHitsByDay deletes all hits on given day.
	DeleteHitsByDay(sql.NullInt64, time.Time) error

	// SaveVisitorsPerDay persists unique visitors per day.
	SaveVisitorsPerDay(*VisitorsPerDay) error

	// SaveVisitorsPerHour persists unique visitors per day and hour.
	SaveVisitorsPerHour(*VisitorsPerHour) error

	// SaveVisitorsPerLanguage persists unique visitors per day and language.
	SaveVisitorsPerLanguage(*VisitorsPerLanguage) error

	// SaveVisitorsPerPage persists unique visitors per day and page.
	SaveVisitorsPerPage(*VisitorsPerPage) error

	// SaveVisitorsPerReferrer persists unique visitors per day and referrer.
	SaveVisitorsPerReferrer(*VisitorsPerReferrer) error

	// Days returns the days at least one hit exists for.
	Days(sql.NullInt64) ([]time.Time, error)

	// CountVisitorsPerDay returns the unique visitor count for per day.
	CountVisitorsPerDay(sql.NullInt64, time.Time) (int, error)

	// CountVisitorsPerDayAndHour returns the unique visitor count per day and hour.
	CountVisitorsPerDayAndHour(sql.NullInt64, time.Time) ([]VisitorsPerHour, error)

	// CountVisitorsPerLanguage returns the unique visitor count per language and day.
	CountVisitorsPerLanguage(sql.NullInt64, time.Time) ([]VisitorsPerLanguage, error)

	// CountVisitorsPerPage returns the unique visitor count per page and day.
	CountVisitorsPerPage(sql.NullInt64, time.Time) ([]VisitorsPerPage, error)

	// CountVisitorsPerReferrer returns the unique visitor count per referrer and day.
	CountVisitorsPerReferrer(sql.NullInt64, time.Time) ([]VisitorsPerReferrer, error)

	// Paths returns distinct paths for page visits.
	// This does not include today.
	Paths(sql.NullInt64, time.Time, time.Time) ([]string, error)

	// Referrer returns distinct referrer for page visits.
	// This does not include today.
	Referrer(sql.NullInt64, time.Time, time.Time) ([]string, error)

	// Visitors returns the visitors within given time frame.
	// This does not include today.
	Visitors(sql.NullInt64, time.Time, time.Time) ([]VisitorsPerDay, error)

	// PageVisits returns the page visits within given time frame for given path.
	// This does not include today.
	PageVisits(sql.NullInt64, string, time.Time, time.Time) ([]VisitorsPerDay, error)

	// ReferrerVisits returns the referrer visits within given time frame for given referrer.
	// This does not include today.
	ReferrerVisits(sql.NullInt64, string, time.Time, time.Time) ([]VisitorsPerReferrer, error)

	// VisitorPages returns the pages within given time frame for unique visitors.
	// It does include today.
	VisitorPages(sql.NullInt64, time.Time, time.Time) ([]VisitorPage, error)

	// VisitorLanguages returns the languages within given time frame for unique visitors.
	// It does include today.
	VisitorLanguages(sql.NullInt64, time.Time, time.Time) ([]VisitorLanguage, error)

	// VisitorReferrer returns the languages within given time frame for unique visitors.
	// It does include today.
	VisitorReferrer(sql.NullInt64, time.Time, time.Time) ([]VisitorReferrer, error)

	// HourlyVisitors returns unique visitors per hour for given time frame.
	// It does include today.
	HourlyVisitors(sql.NullInt64, time.Time, time.Time) ([]HourlyVisitors, error)

	// ActiveVisitors returns unique visitors starting at given time.
	ActiveVisitors(sql.NullInt64, time.Time) (int, error)

	// ActiveVisitorsPerPage returns unique visitors per page starting at given time.
	ActiveVisitorsPerPage(sql.NullInt64, time.Time) ([]PageVisitors, error)

	// CountHits returns the number of hits for given tenant ID.
	CountHits(sql.NullInt64) int

	// VisitorsPerDay returns all visitors per day for given tenant ID in order.
	VisitorsPerDay(sql.NullInt64) []VisitorsPerDay

	// VisitorsPerHour returns all visitors per hour for given tenant ID in order.
	VisitorsPerHour(sql.NullInt64) []VisitorsPerHour

	// VisitorsPerLanguage returns all visitors per language for given tenant ID in alphabetical order.
	VisitorsPerLanguage(sql.NullInt64) []VisitorsPerLanguage

	// VisitorsPerPage returns all visitors per page for given tenant ID in alphabetical order.
	VisitorsPerPage(sql.NullInt64) []VisitorsPerPage

	// VisitorsPerReferrer returns all visitors per referrer for given tenant ID in alphabetical order.
	VisitorsPerReferrer(sql.NullInt64) []VisitorsPerReferrer
}
