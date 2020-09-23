package pirsch

import (
	"sort"
	"time"
)

// PathVisitors represents visitor statistics per day for a path.
type PathVisitors struct {
	Path  string
	Stats []Stats
}

// TimeOfDayVisitors represents the visitor count per day and hour for a path.
type TimeOfDayVisitors struct {
	Day   time.Time
	Stats []VisitorTimeStats
}

// Analyzer provides an interface to analyze processed data and hits.
type Analyzer struct {
	store Store
}

// NewAnalyzer returns a new Analyzer for given Store.
func NewAnalyzer(store Store) *Analyzer {
	return &Analyzer{store}
}

// ActiveVisitors returns the active visitors per path and the total number of active visitors for given duration.
// Use time.Minute*5 for example to see the active visitors for the past 5 minutes.
func (analyzer *Analyzer) ActiveVisitors(filter *Filter, duration time.Duration) ([]Stats, int, error) {
	filter = analyzer.getFilter(filter)
	from := time.Now().UTC().Add(-duration)
	stats, err := analyzer.store.ActivePageVisitors(QueryParams{TenantID: filter.TenantID}, from)

	if err != nil {
		return nil, 0, err
	}

	return stats, analyzer.store.ActiveVisitors(QueryParams{TenantID: filter.TenantID}, from), nil
}

// Visitors returns the visitor count, session count, and bounce rate per day.
func (analyzer *Analyzer) Visitors(filter *Filter) ([]Stats, error) {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats, err := analyzer.store.Visitors(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	if addToday {
		visitorsToday := analyzer.store.CountVisitors(nil, QueryParams{TenantID: filter.TenantID}, today)
		bouncesToday := analyzer.store.CountVisitorsByPathAndMaxOneHit(nil, QueryParams{TenantID: filter.TenantID}, today, "")

		if len(stats) > 0 {
			if visitorsToday != nil {
				stats[len(stats)-1].Visitors += visitorsToday.Visitors
				stats[len(stats)-1].Sessions += visitorsToday.Sessions
				stats[len(stats)-1].Bounces += bouncesToday
			}
		} else {
			stats = append(stats, Stats{
				Visitors: visitorsToday.Visitors,
				Sessions: visitorsToday.Sessions,
				Bounces:  visitorsToday.Bounces,
			})
		}
	}

	for i := range stats {
		if stats[i].Visitors > 0 {
			stats[i].BounceRate = float64(stats[i].Bounces) / float64(stats[i].Visitors)
		}
	}

	return stats, nil
}

// VisitorHours returns the visitor and session count grouped by hour of day for given time frame.
func (analyzer *Analyzer) VisitorHours(filter *Filter) ([]VisitorTimeStats, error) {
	filter = analyzer.getFilter(filter)
	stats, err := analyzer.store.VisitorHours(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// Languages returns the visitor count per language.
func (analyzer *Analyzer) Languages(filter *Filter) ([]LanguageStats, error) {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats, err := analyzer.store.VisitorLanguages(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	if addToday {
		visitorsToday, err := analyzer.store.CountVisitorsByLanguage(nil, QueryParams{TenantID: filter.TenantID}, today)

		if err != nil {
			return nil, err
		}

		for _, v := range visitorsToday {
			found := false

			for i, s := range stats {
				if s.Language.String == v.Language.String {
					stats[i].Visitors += v.Visitors
					found = true
					break
				}
			}

			if !found {
				stats = append(stats, v)
			}
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Visitors > stats[j].Visitors
	})

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// Referrer returns the visitor count per referrer.
func (analyzer *Analyzer) Referrer(filter *Filter) ([]ReferrerStats, error) {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats, err := analyzer.store.VisitorReferrer(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	if addToday {
		visitorsToday, err := analyzer.store.CountVisitorsByReferrer(nil, QueryParams{TenantID: filter.TenantID}, today)

		if err != nil {
			return nil, err
		}

		for _, v := range visitorsToday {
			found := false

			for i, s := range stats {
				if s.Referrer.String == v.Referrer.String {
					stats[i].Visitors += v.Visitors
					found = true
					break
				}
			}

			if !found {
				stats = append(stats, v)
			}
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Visitors > stats[j].Visitors
	})

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// OS returns the visitor count per operating system.
func (analyzer *Analyzer) OS(filter *Filter) ([]OSStats, error) {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats, err := analyzer.store.VisitorOS(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	if addToday {
		visitorsToday, err := analyzer.store.CountVisitorsByOS(nil, QueryParams{TenantID: filter.TenantID}, today)

		if err != nil {
			return nil, err
		}

		for _, v := range visitorsToday {
			found := false

			for i, s := range stats {
				if s.OS.String == v.OS.String {
					stats[i].Visitors += v.Visitors
					found = true
					break
				}
			}

			if !found {
				stats = append(stats, v)
			}
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Visitors > stats[j].Visitors
	})

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// Browser returns the visitor count per browser.
func (analyzer *Analyzer) Browser(filter *Filter) ([]BrowserStats, error) {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats, err := analyzer.store.VisitorBrowser(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	if addToday {
		visitorsToday, err := analyzer.store.CountVisitorsByBrowser(nil, QueryParams{TenantID: filter.TenantID}, today)

		if err != nil {
			return nil, err
		}

		for _, v := range visitorsToday {
			found := false

			for i, s := range stats {
				if s.Browser.String == v.Browser.String {
					stats[i].Visitors += v.Visitors
					found = true
					break
				}
			}

			if !found {
				stats = append(stats, v)
			}
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Visitors > stats[j].Visitors
	})

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// Platform returns the visitor count per browser.
func (analyzer *Analyzer) Platform(filter *Filter) *VisitorStats {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats := analyzer.store.VisitorPlatform(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if stats == nil {
		stats = &VisitorStats{}
	}

	if addToday {
		visitorsToday := analyzer.store.CountVisitorsByPlatform(nil, QueryParams{TenantID: filter.TenantID}, today)

		if visitorsToday != nil {
			stats.PlatformDesktop += visitorsToday.PlatformDesktop
			stats.PlatformMobile += visitorsToday.PlatformMobile
			stats.PlatformUnknown += visitorsToday.PlatformUnknown
		}
	}

	sum := float64(stats.PlatformDesktop + stats.PlatformMobile + stats.PlatformUnknown)
	stats.RelativePlatformDesktop = float64(stats.PlatformDesktop) / sum
	stats.RelativePlatformMobile = float64(stats.PlatformMobile) / sum
	stats.RelativePlatformUnknown = float64(stats.PlatformUnknown) / sum
	return stats
}

// Screen returns the visitor count per screen size (width and height).
func (analyzer *Analyzer) Screen(filter *Filter) ([]ScreenStats, error) {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats, err := analyzer.store.VisitorScreenSize(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	if addToday {
		visitorsToday, err := analyzer.store.CountVisitorsByScreenSize(nil, QueryParams{TenantID: filter.TenantID}, today)

		if err != nil {
			return nil, err
		}

		for _, v := range visitorsToday {
			found := false

			for i, s := range stats {
				if s.Width == v.Width && s.Height == v.Height {
					stats[i].Visitors += v.Visitors
					found = true
					break
				}
			}

			if !found {
				stats = append(stats, v)
			}
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Visitors > stats[j].Visitors
	})

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// Country returns the visitor count per country.
func (analyzer *Analyzer) Country(filter *Filter) ([]CountryStats, error) {
	filter = analyzer.getFilter(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats, err := analyzer.store.VisitorCountry(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	if addToday {
		visitorsToday, err := analyzer.store.CountVisitorsByCountryCode(nil, QueryParams{TenantID: filter.TenantID}, today)

		if err != nil {
			return nil, err
		}

		for _, v := range visitorsToday {
			found := false

			for i, s := range stats {
				if s.CountryCode == v.CountryCode {
					stats[i].Visitors += v.Visitors
					found = true
					break
				}
			}

			if !found {
				stats = append(stats, v)
			}
		}
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Visitors > stats[j].Visitors
	})

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// TimeOfDay returns the visitor count per day and hour for given time frame.
func (analyzer *Analyzer) TimeOfDay(filter *Filter) ([]TimeOfDayVisitors, error) {
	filter = analyzer.getFilter(filter)
	from := filter.From
	stats := make([]TimeOfDayVisitors, 0)

	for !from.After(filter.To) {
		s, err := analyzer.VisitorHours(&Filter{From: from, To: from})

		if err != nil {
			return nil, err
		}

		stats = append(stats, TimeOfDayVisitors{
			Day:   from,
			Stats: s,
		})
		from = from.Add(time.Hour * 24)
	}

	return stats, nil
}

// PageVisitors returns the visitor count, session count, and bounce rate per day for the given time frame grouped by path.
func (analyzer *Analyzer) PageVisitors(filter *Filter) ([]PathVisitors, error) {
	filter = analyzer.getFilter(filter)
	paths := analyzer.getPaths(filter)
	today := today()
	addToday := today.Equal(filter.To)
	stats := make([]PathVisitors, 0, len(paths))

	for _, path := range paths {
		visitors, err := analyzer.store.PageVisitors(QueryParams{TenantID: filter.TenantID}, path, filter.From, filter.To)

		if err != nil {
			return nil, err
		}

		if addToday {
			visitorsToday, err := analyzer.store.CountVisitorsByPath(nil, QueryParams{TenantID: filter.TenantID}, today, path, false)

			if err != nil {
				return nil, err
			}

			bouncesToday := analyzer.store.CountVisitorsByPathAndMaxOneHit(nil, QueryParams{TenantID: filter.TenantID}, today, path)

			if len(visitorsToday) > 0 {
				if len(visitors) > 0 {
					visitors[len(visitors)-1].Visitors += visitorsToday[0].Visitors
					visitors[len(visitors)-1].Sessions += visitorsToday[0].Sessions
					visitors[len(visitors)-1].Bounces += bouncesToday
				} else {
					visitors = append(visitors, Stats{
						Visitors: visitorsToday[0].Visitors,
						Sessions: visitorsToday[0].Sessions,
						Bounces:  bouncesToday,
					})
				}
			}
		}

		for i := range visitors {
			if visitors[i].Visitors > 0 {
				visitors[i].BounceRate = float64(visitors[i].Bounces) / float64(visitors[i].Visitors)
			}
		}

		stats = append(stats, PathVisitors{
			Path:  path,
			Stats: visitors,
		})
	}

	return stats, nil
}

// PageLanguages returns the visitor count per language, day, path, and for the given time frame.
// The path is mandatory.
func (analyzer *Analyzer) PageLanguages(filter *Filter) ([]LanguageStats, error) {
	filter = analyzer.getFilter(filter)

	if filter.Path == "" {
		return []LanguageStats{}, nil
	}

	stats, err := analyzer.store.PageLanguages(QueryParams{TenantID: filter.TenantID}, filter.Path, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// PageReferrer returns the visitor count per referrer, day, path, and for the given time frame.
// The path is mandatory.
func (analyzer *Analyzer) PageReferrer(filter *Filter) ([]ReferrerStats, error) {
	filter = analyzer.getFilter(filter)

	if filter.Path == "" {
		return []ReferrerStats{}, nil
	}

	stats, err := analyzer.store.PageReferrer(QueryParams{TenantID: filter.TenantID}, filter.Path, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// PageOS returns the visitor count per operating system, day, path, and for the given time frame.
// The path is mandatory.
func (analyzer *Analyzer) PageOS(filter *Filter) ([]OSStats, error) {
	filter = analyzer.getFilter(filter)

	if filter.Path == "" {
		return []OSStats{}, nil
	}

	stats, err := analyzer.store.PageOS(QueryParams{TenantID: filter.TenantID}, filter.Path, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// PageBrowser returns the visitor count per brower, day, path, and for the given time frame.
// The path is mandatory.
func (analyzer *Analyzer) PageBrowser(filter *Filter) ([]BrowserStats, error) {
	filter = analyzer.getFilter(filter)

	if filter.Path == "" {
		return []BrowserStats{}, nil
	}

	stats, err := analyzer.store.PageBrowser(QueryParams{TenantID: filter.TenantID}, filter.Path, filter.From, filter.To)

	if err != nil {
		return nil, err
	}

	var sum float64

	for i := range stats {
		sum += float64(stats[i].Visitors)
	}

	for i := range stats {
		stats[i].RelativeVisitors = float64(stats[i].Visitors) / sum
	}

	return stats, nil
}

// PagePlatform returns the visitor count per platform, day, path, and for the given time frame.
// The path is mandatory.
func (analyzer *Analyzer) PagePlatform(filter *Filter) *VisitorStats {
	filter = analyzer.getFilter(filter)

	if filter.Path == "" {
		return &VisitorStats{}
	}

	stats := analyzer.store.PagePlatform(QueryParams{TenantID: filter.TenantID}, filter.Path, filter.From, filter.To)

	if stats == nil {
		return &VisitorStats{}
	}

	sum := float64(stats.PlatformDesktop + stats.PlatformMobile + stats.PlatformUnknown)
	stats.RelativePlatformDesktop = float64(stats.PlatformDesktop) / sum
	stats.RelativePlatformMobile = float64(stats.PlatformMobile) / sum
	stats.RelativePlatformUnknown = float64(stats.PlatformUnknown) / sum
	return stats
}

// getFilter validates and returns the given filter or a default filter if it is nil.
func (analyzer *Analyzer) getFilter(filter *Filter) *Filter {
	if filter == nil {
		return NewFilter(NullTenant)
	}

	filter.validate()
	return filter
}

// getPaths returns the paths to filter for. This can either be the one passed in,
// or all relevant paths for the given time frame otherwise.
func (analyzer *Analyzer) getPaths(filter *Filter) []string {
	if filter.Path != "" {
		return []string{filter.Path}
	}

	paths, err := analyzer.store.Paths(QueryParams{TenantID: filter.TenantID}, filter.From, filter.To)

	if err != nil {
		return []string{}
	}

	return paths
}
