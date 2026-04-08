package mercure

import (
	"regexp"
	"strings"

	"github.com/maypok86/otter/v2"
	"github.com/yosida95/uritemplate/v3"
)

// DefaultTopicSelectorStoreCacheSize is the default maximum number of entries in the cache.
const DefaultTopicSelectorStoreCacheSize = 2_560_000

type matchCacheKey struct {
	topicSelector string
	topic         string
}

// TopicSelectorStore caches compiled templates to improve memory and CPU usage.
type TopicSelectorStore struct {
	matchCache    *otter.Cache[matchCacheKey, bool]
	templateCache *otter.Cache[string, *regexp.Regexp]
}

// NewTopicSelectorStore creates a TopicSelectorStore.
// If cacheSize > 0, match results and compiled templates are cached.
func NewTopicSelectorStore(cacheSize int) (*TopicSelectorStore, error) {
	if cacheSize <= 0 {
		return &TopicSelectorStore{}, nil
	}

	matchCache, err := otter.New[matchCacheKey, bool](&otter.Options[matchCacheKey, bool]{
		MaximumSize: cacheSize,
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	templateCache, err := otter.New[string, *regexp.Regexp](&otter.Options[string, *regexp.Regexp]{
		MaximumSize: cacheSize / 10, // Templates are fewer but larger
	})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &TopicSelectorStore{matchCache: matchCache, templateCache: templateCache}, nil
}

func (tss *TopicSelectorStore) match(topic, topicSelector string) bool {
	// Always do an exact matching comparison first
	// Also check if the topic selector is the reserved keyword *
	if topicSelector == "*" || topic == topicSelector {
		return true
	}

	k := matchCacheKey{topicSelector: topicSelector, topic: topic}

	if tss.matchCache != nil {
		if value, found := tss.matchCache.GetIfPresent(k); found {
			return value
		}
	}

	r := tss.getRegexp(topicSelector)
	if r == nil {
		return false
	}

	// Use template.Regexp() instead of template.Match() for performance
	// See https://github.com/yosida95/uritemplate/pull/7
	match := r.MatchString(topic)
	if tss.matchCache != nil {
		tss.matchCache.Set(k, match)
	}

	return match
}

// getRegexp retrieves regexp for this template selector.
func (tss *TopicSelectorStore) getRegexp(topicSelector string) *regexp.Regexp {
	// If it's definitely not a URI template, skip to save some resources
	if !strings.Contains(topicSelector, "{") {
		return nil
	}

	if tss.templateCache != nil {
		if r, found := tss.templateCache.GetIfPresent(topicSelector); found {
			return r
		}
	}

	// If an error occurs, it's a raw string
	if tpl, err := uritemplate.New(topicSelector); err == nil {
		r := tpl.Regexp()
		if tss.templateCache != nil {
			tss.templateCache.Set(topicSelector, r)
		}

		return r
	}

	return nil
}
