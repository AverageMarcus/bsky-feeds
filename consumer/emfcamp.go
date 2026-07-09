package consumer

import (
	"database/sql"
	dbpkg "jetstream-feed-generator/db/sqlc"
	"log/slog"
	"regexp"
	"slices"
	"strings"
	"time"

	apibsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/jetstream/pkg/models"
)

type EMFCampFeed struct {
	name   string
	logger *slog.Logger
	q      *dbpkg.Queries
}

func NewEMFCampFeed(name string, logger *slog.Logger, db *sql.DB) *EMFCampFeed {
	feedLogger := logger.With("feed", name)
	return &EMFCampFeed{name, feedLogger, dbpkg.New(db)}
}

func (f *EMFCampFeed) Name() string {
	return f.name
}

func (f *EMFCampFeed) DB() *dbpkg.Queries {
	return f.q
}

func (f *EMFCampFeed) Match(event *models.Event, post *apibsky.FeedPost) bool {
	if len(post.Text) == 0 {
		return false
	}

	var (
		// Official hashtags we want to include
		hashtags = []string{
			"EMFCamp", "EMF2026", "EMF26",
		}
		// Posts from accounts (DIDs) to always include
		accounts = []string{
			"did:plc:r5tbkz2suj4hz6kyadj73y6n", // emfcamp.bsky.social
		}
		// Some common phrases that might be used without hashtags that we'd want to match on
		re = regexp.MustCompile(`(?mi)(^|\s|#)(EMF ?Camp)(\d{2,4})?(\W|$)`)
	)

	// Used to include year-specific hashtags for each
	yearFull := time.Now().Format("2006")
	yearShort := time.Now().Format("06")

	// Check if a relevant account to always include posts from
	if event != nil && event.Account != nil && slices.Contains(accounts, event.Account.Did) {
		return true
	}

	// First check for relevant hashtags matching exactly (case-insensitive)
	for _, facet := range post.Facets {
		for _, feat := range facet.Features {
			if feat.RichtextFacet_Tag != nil {
				hashtag := feat.RichtextFacet_Tag.Tag
				if slices.ContainsFunc(hashtags, func(h string) bool {
					return strings.EqualFold(h, hashtag) || strings.EqualFold(h, hashtag+yearFull) || strings.EqualFold(h, hashtag+yearShort)
				}) {
					return true
				}
			}
		}
	}

	// Finally attempt some more generic non-hashtag matches
	return re.MatchString(post.Text)
}
