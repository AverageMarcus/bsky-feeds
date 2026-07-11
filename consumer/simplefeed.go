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

type SimpleFeedConfig struct {
	Name        string   `mapstructure:"name"`
	Hashtags    []string `mapstructure:"hashtags"`
	AccountDIDs []string `mapstructure:"account_dids"`
	Regex       string   `mapstructure:"regex"`
	IncludeYear bool     `mapstructure:"include_year"`
}

type SimpleFeed struct {
	logger *slog.Logger
	q      *dbpkg.Queries

	config SimpleFeedConfig
}

func NewSimpleFeed(cfg SimpleFeedConfig, logger *slog.Logger, db *sql.DB) *SimpleFeed {
	feedLogger := logger.With("feed", cfg.Name)
	return &SimpleFeed{feedLogger, dbpkg.New(db), cfg}
}

func (f *SimpleFeed) Name() string {
	return f.config.Name
}

func (f *SimpleFeed) DB() *dbpkg.Queries {
	return f.q
}

func (f *SimpleFeed) Match(event *models.Event, post *apibsky.FeedPost) bool {
	if len(post.Text) == 0 {
		return false
	}

	// Used to include year-specific hashtags for each
	yearFull := time.Now().Format("2006")
	yearShort := time.Now().Format("06")

	// Check if a relevant account to always include posts from
	if event != nil && event.Account != nil && slices.Contains(f.config.AccountDIDs, event.Account.Did) {
		return true
	}

	// First check for relevant hashtags matching exactly (case-insensitive)
	for _, facet := range post.Facets {
		for _, feat := range facet.Features {
			if feat.RichtextFacet_Tag != nil {
				hashtag := feat.RichtextFacet_Tag.Tag
				if slices.ContainsFunc(f.config.Hashtags, func(h string) bool {
					return strings.EqualFold(h, hashtag) || (f.config.IncludeYear && (strings.EqualFold(h, hashtag+yearFull) || strings.EqualFold(h, hashtag+yearShort)))
				}) {
					return true
				}
			}
		}
	}

	if f.config.Regex != "" {
		re := regexp.MustCompile(f.config.Regex)
		return re.MatchString(post.Text)
	}

	return false
}
