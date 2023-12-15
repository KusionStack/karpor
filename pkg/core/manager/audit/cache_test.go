package audit

import (
	"testing"
	"time"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(time.Minute) // Set expiration time to 1 minute
	locator := core.Locator{
		Cluster:    "test",
		Group:      "group",
		APIVersion: "v1",
		Kind:       "testKind",
		Namespace:  "default",
		Name:       "testResource",
	}

	// Create an AuditData instance to be stored in the cache
	data := &AuditData{
		Issues: []*scanner.Issue{
			{Title: "issue1"},
			{Title: "issue2"},
		},
		Aggregated: map[string]int{"issue1": 2, "issue2": 1},
	}

	// Set data in cache with the locator
	cache.Set(locator, data)

	// Retrieve data from the cache
	retrievedData, exists := cache.Get(locator)
	assert.True(t, exists)
	assert.NotNil(t, retrievedData)
	assert.Equal(t, data.Issues, retrievedData.Issues)
	assert.Equal(t, data.Aggregated, retrievedData.Aggregated)
}

func TestCache_GetExpiredData(t *testing.T) {
	cache := NewCache(time.Millisecond) // Set expiration time to 1 millisecond
	locator := core.Locator{
		Cluster:    "test",
		Group:      "group",
		APIVersion: "v1",
		Kind:       "Kind",
		Namespace:  "default",
		Name:       "testResource",
	}

	data := &AuditData{
		Issues: []*scanner.Issue{
			{Title: "issue1"},
		},
		Aggregated: map[string]int{"issue1": 1},
	}

	// Set data in cache with the locator
	cache.Set(locator, data)

	// Wait for the cache item to expire
	time.Sleep(2 * time.Millisecond)

	// Retrieve data from the cache after expiration
	retrievedData, exists := cache.Get(locator)
	assert.False(t, exists)
	assert.Nil(t, retrievedData)
}
