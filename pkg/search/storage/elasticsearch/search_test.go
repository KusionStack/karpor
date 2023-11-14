package elasticsearch

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/require"
)

func TestESClient_Search(t *testing.T) {
	cl, err := NewESClient(elasticsearch.Config{
		Addresses: []string{
			"http://100.88.101.58:9200",
		},
	})
	require.NoError(t, err)
	res, err := cl.searchByDSL(context.TODO(), "apiVersion=v1")
	require.NoError(t, err)
	t.Log(res)
}

func TestTranslate(t *testing.T) {
	str := "select * from `elastic-default-index` where apiVersion='v1'"
	cl, err := NewESClient(elasticsearch.Config{
		Addresses: []string{
			"http://100.88.101.58:9200",
		},
	})
	require.NoError(t, err)
	cl.searchBySQL(context.TODO(), str)
	res, err := cl.searchBySQL(context.TODO(), str)
	require.NoError(t, err)
	fmt.Printf("%v", res)
}

func TestEscape(t *testing.T) {
	str := "select * from `elastic-default-index` where apiVersion='v1'"
	estr := url.QueryEscape(str)
	fmt.Println(estr)
}
