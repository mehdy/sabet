package rss

import (
	"github.com/mehdy/sabet/pkg/jobs/meta"
	"github.com/mmcdole/gofeed"
)

type Job struct {
	meta.TypeMeta  `json:",inline"`
	meta.JobMeta   `json:"metadata,omitempty"`
	meta.StoreMeta `json:"store,omitempty"`

	Spec struct {
		Sources []Source `json:"sources,omitempty"`
	} `json:"spec,omitempty"`

	store meta.Store
}

type Source struct {
	URL        string   `json:"url,omitempty"`
	Title      string   `json:"title,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Languages  []string `json:"languages,omitempty"`
}

type Result struct {
	Items []*gofeed.Item `json:"items,omitempty"`
}
