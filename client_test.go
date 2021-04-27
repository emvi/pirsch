package pirsch

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("tcp://127.0.0.1:9000", nil)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NoError(t, client.DB.Ping())
}

func TestClient_SaveHit(t *testing.T) {
	cleanupDB()
	assert.NoError(t, dbClient.SaveHits([]Hit{
		{
			ClientID:       1,
			Fingerprint:    "fp",
			Time:           time.Now(),
			Session:        sql.NullTime{Time: time.Now(), Valid: true},
			UserAgent:      "ua",
			Path:           "/path",
			Language:       "en",
			Referrer:       sql.NullString{String: "ref", Valid: true},
			ReferrerName:   sql.NullString{String: "ref_name", Valid: true},
			ReferrerIcon:   sql.NullString{String: "ref_icon", Valid: true},
			OS:             "os",
			OSVersion:      "10",
			Browser:        "browser",
			BrowserVersion: "89",
			CountryCode:    "en",
			Desktop:        true,
			Mobile:         false,
			ScreenWidth:    1920,
			ScreenHeight:   1080,
			ScreenClass:    "XL",
		},
		{
			Fingerprint: "fp",
			Time:        time.Now(),
			UserAgent:   "ua",
			Path:        "/path",
		},
	}))
}

func TestClient_Session(t *testing.T) {
	cleanupDB()
	fp := "session_fp"
	session := Today()
	assert.NoError(t, dbClient.SaveHits([]Hit{
		{
			ClientID:    1,
			Fingerprint: fp,
			Time:        time.Now(),
			Session:     sql.NullTime{Time: session, Valid: true},
			UserAgent:   "ua",
			Path:        "/path",
		},
	}))
	s, err := dbClient.Session(1, fp, time.Now().Add(-time.Second))
	assert.NoError(t, err)
	assert.Equal(t, session, s)
}
