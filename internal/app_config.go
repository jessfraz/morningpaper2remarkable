package internal

import (
	"time"

	"github.com/pseudo-su/morningpaper2remarkable/internal/remarkable"
)

type AppConfig struct {
	Debug   bool
	DataDir string

	Interval time.Duration
	Once     bool
	MaxPages int

	RemarkableAPI remarkable.Remarkable
	MorningPaperRSSFeedURL string
	DefaultDir string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		MorningPaperRSSFeedURL: "https://blog.acolyer.org/feed/?paged=%d",
		DefaultDir: "morningpaper",
	}
}
