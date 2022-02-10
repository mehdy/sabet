package rss

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/mehdy/sabet/pkg/jobs/meta"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

func (j *Job) SetStore(store meta.Store) {
	j.store = store
}

func urlToKey(url string) string {
	key := sha1.New()
	return hex.EncodeToString(key.Sum([]byte(url)))[:8]
}

func (j Job) Init() error {
	return nil
}

func (j *Job) log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"type": j.GetType(),
		"name": j.GetName(),
	})
}

func (j *Job) Execute(_ io.Reader) (io.Reader, error) {
	fp := gofeed.NewParser()
	newItems := make([]*gofeed.Item, 0)

	for _, step := range j.Spec.Sources {

		now := time.Now()

		lastUpdateKey := urlToKey(step.URL) + "_last_update"
		lastUpdateRaw, err := j.store.Get(lastUpdateKey)
		if err != nil {
			j.log().WithField("url", step.URL).WithError(err).Error("Failed to get last update time")
			continue
		}
		if lastUpdateRaw == nil {
			lastUpdateRaw = []byte("0")
		}

		var lastUpdateTime time.Time

		timestamp, err := strconv.Atoi(string(lastUpdateRaw))
		if err != nil {
			j.log().WithField("url", step.URL).WithError(err).Error("Failed to parse last update time")
			continue
		}
		lastUpdateTime = time.Unix(int64(timestamp), 0)

		j.log().WithField("url", step.URL).Debugf("Last updated at %s", lastUpdateTime)

		j.log().WithField("url", step.URL).Info("Fetching feed")
		feed, err := fp.ParseURL(step.URL)
		if err != nil {
			j.log().WithField("url", step.URL).WithError(err).Error("Failed to fetch feed")
			continue
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
			j.log().WithField("url", step.URL).WithError(err).Error("Failed to save last update time")
			continue
		}
	}

	output, err := json.Marshal(Result{Items: newItems})
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(output), nil
}
