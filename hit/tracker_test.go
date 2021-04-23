package hit

/*
func TestTrackerConfigValidate(t *testing.T) {
	cfg := &TrackerConfig{}
	cfg.validate()
	assert.Equal(t, runtime.NumCPU(), cfg.Worker)
	assert.Equal(t, defaultWorkerBufferSize, cfg.WorkerBufferSize)
	assert.Equal(t, defaultWorkerTimeout, cfg.WorkerTimeout)
	assert.Len(t, cfg.ReferrerDomainBlacklist, 0)
	assert.False(t, cfg.ReferrerDomainBlacklistIncludesSubdomains)
	cfg = &TrackerConfig{
		Worker:                  123,
		WorkerBufferSize:        42,
		WorkerTimeout:           time.Second * 57,
		ReferrerDomainBlacklist: []string{"localhost"},
		ReferrerDomainBlacklistIncludesSubdomains: true,
	}
	cfg.validate()
	assert.Equal(t, 123, cfg.Worker)
	assert.Equal(t, 42, cfg.WorkerBufferSize)
	assert.Equal(t, time.Second*57, cfg.WorkerTimeout)
	assert.Len(t, cfg.ReferrerDomainBlacklist, 1)
	assert.True(t, cfg.ReferrerDomainBlacklistIncludesSubdomains)
	cfg = &TrackerConfig{WorkerTimeout: time.Second * 142}
	cfg.validate()
	assert.Equal(t, maxWorkerTimeout, cfg.WorkerTimeout)
}

func TestTrackerHitTimeout(t *testing.T) {
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.Header.Add("User-Agent", "valid")
	req2 := httptest.NewRequest(http.MethodGet, "/hello-world", nil)
	req2.Header.Add("User-Agent", "valid")
	store := newTestStore()
	tracker := NewTracker(store, "salt", &TrackerConfig{WorkerTimeout: time.Millisecond * 200})
	tracker.Hit(req1, nil)
	tracker.Hit(req2, nil)
	time.Sleep(time.Millisecond * 210)
	assert.Len(t, store.hits, 2)

	// ignore order...
	if store.hits[0].Path != "/" && store.hits[0].Path != "/hello-world" ||
		store.hits[1].Path != "/" && store.hits[1].Path != "/hello-world" {
		t.Fatalf("Hits not as expected: %v %v", store.hits[0], store.hits[1])
	}
}

func TestTrackerHitLimit(t *testing.T) {
	store := newTestStore()
	tracker := NewTracker(store, "salt", &TrackerConfig{
		Worker:           1,
		WorkerBufferSize: 10,
	})

	for i := 0; i < 7; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("User-Agent", "valid")
		tracker.Hit(req, nil)
	}

	tracker.Stop()
	assert.Len(t, store.hits, 7)
}

func TestTrackerHitDiscard(t *testing.T) {
	store := newTestStore()
	tracker := NewTracker(store, "salt", &TrackerConfig{
		Worker:           1,
		WorkerBufferSize: 5,
	})

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("User-Agent", "valid")
		tracker.Hit(req, nil)

		if i > 3 {
			tracker.Stop()
		}
	}

	assert.Len(t, store.hits, 5)
}

func TestTrackerCountryCode(t *testing.T) {
	geoDB, err := NewGeoDB(GeoDBConfig{
		File: filepath.Join("geodb/GeoIP2-Country-Test.mmdb"),
	})
	assert.NoError(t, err)
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.Header.Add("User-Agent", "valid")
	req1.RemoteAddr = "81.2.69.142"
	req2 := httptest.NewRequest(http.MethodGet, "/hello-world", nil)
	req2.Header.Add("User-Agent", "valid")
	req2.RemoteAddr = "127.0.0.1"
	store := newTestStore()
	tracker := NewTracker(store, "salt", &TrackerConfig{
		WorkerTimeout: time.Second,
	})
	tracker.SetGeoDB(geoDB)
	tracker.Hit(req1, nil)
	tracker.Hit(req2, nil)
	tracker.Stop()
	assert.Len(t, store.hits, 2)
	foundGB := false
	foundEmpty := false

	for _, hit := range store.hits {
		if hit.CountryCode.String == "gb" {
			foundGB = true
		} else if hit.CountryCode.String == "" {
			foundEmpty = true
		}
	}

	assert.True(t, foundGB)
	assert.True(t, foundEmpty)
}

func TestTrackerHitSession(t *testing.T) {
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.Header.Add("User-Agent", "valid")
	req2 := httptest.NewRequest(http.MethodGet, "/hello-world", nil)
	req2.Header.Add("User-Agent", "valid")
	store := newTestStore()
	tracker := NewTracker(store, "salt", &TrackerConfig{
		WorkerTimeout: time.Second,
		Sessions:      true,
	})
	tracker.Hit(req1, nil)
	tracker.Hit(req2, nil)
	tracker.Stop()
	assert.Len(t, store.hits, 2)

	// ignore order...
	if !store.hits[0].Session.Valid || !store.hits[1].Session.Valid ||
		store.hits[0].Session.Time.IsZero() || store.hits[1].Session.Time.IsZero() {
		t.Fatalf("Hits not as expected: %v %v", store.hits[0], store.hits[1])
	}
}

func TestTrackerIgnoreSubdomain(t *testing.T) {
	store := newTestStore()
	tracker := NewTracker(store, "salt", &TrackerConfig{
		WorkerTimeout: time.Second,
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("User-Agent", "valid")
	req.RemoteAddr = "81.2.69.142"
	tracker.Hit(req, &Options{
		ReferrerDomainBlacklist: []string{"pirsch.io"},
		Referrer:                "https://pirsch.io/",
	})
	tracker.Hit(req, &Options{
		ReferrerDomainBlacklist:                   []string{"pirsch.io"},
		ReferrerDomainBlacklistIncludesSubdomains: true,
		Referrer: "https://www.pirsch.io/",
	})
	tracker.Hit(req, &Options{
		ReferrerDomainBlacklist: []string{"pirsch.io", "www.pirsch.io"},
		Referrer:                "https://www.pirsch.io/",
	})
	tracker.Hit(req, &Options{
		ReferrerDomainBlacklist: []string{"pirsch.io"},
		Referrer:                "pirsch.io",
	})
	tracker.Stop()
	assert.Len(t, store.hits, 4)

	for _, hit := range store.hits {
		assert.False(t, hit.Referrer.Valid)
	}
}

func BenchmarkTracker(b *testing.B) {
	geoDB, err := NewGeoDB(GeoDBConfig{
		File: filepath.Join("geodb/GeoIP2-Country-Test.mmdb"),
	})
	assert.NoError(b, err)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("User-Agent", "valid")
	req.RemoteAddr = "81.2.69.142"
	store := NewPostgresStore(postgresDB, nil)
	tracker := NewTracker(store, "salt", &TrackerConfig{
		WorkerTimeout: time.Second,
		Sessions:      true,
		GeoDB:         geoDB,
	})

	for i := 0; i < 10000; i++ {
		tracker.Hit(req, nil)
	}

	tracker.Stop()
}
*/