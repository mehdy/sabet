package rss

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/mehdy/sabet/pkg/jobs/meta"
	"github.com/mmcdole/gofeed"
)

func (j *Job) SetStore(store meta.Store) {
	j.store = store
}

func urlToKey(url string) string {
	key := sha1.New()
	return hex.EncodeToString(key.Sum([]byte(url)))
}

func (j Job) Init() error {
	return nil
}

func (j *Job) Execute(_ io.Reader) (io.Reader, error) {
	fp := gofeed.NewParser()
	newItems := make([]*gofeed.Item, 0)

	for _, step := range j.Spec.Sources {
		now := time.Now()

		lastUpdateKey := urlToKey(step.URL) + "_last_update"
		lastUpdateRaw, err := j.store.Get(lastUpdateKey)
		if err != nil {
			return nil, err
		}
		if lastUpdateRaw == nil {
			lastUpdateRaw = []byte("0")
		}

		var lastUpdateTime time.Time

		timestamp, err := strconv.Atoi(string(lastUpdateRaw))
		if err != nil {
			return nil, err
		}
		lastUpdateTime = time.Unix(int64(timestamp), 0)

		feed, err := fp.ParseURL(step.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse feed: %v", err)
		}

		for _, item := range feed.Items {
			if item.PublishedParsed == nil {
				item.PublishedParsed = &now
			}

			if item.PublishedParsed.After(lastUpdateTime) {
				newItems = append(newItems, item)
			}
		}

		if err := j.store.Put(lastUpdateKey, []byte(strconv.Itoa(int(now.Unix())))); err != nil {
			return nil, err
		}
	}

	output, err := json.Marshal(Result{Items: newItems})
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(output), nil
}
