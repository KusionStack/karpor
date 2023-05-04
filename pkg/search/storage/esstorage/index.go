// Copyright The Karbour Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package esstorage

import (
	"context"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	defaultIndexName = "elastic-default-index"
	defaultMapping   = `{
  "settings":{
    "index":{
      "number_of_shards":1,
      "auto_expand_replicas":"0-1",
      "number_of_replicas":0
    },
    "analysis":{
      "normalizer":{
        "lowercase":{
          "type":"custom",
          "filter":[
            "lowercase"
          ]
        }
      }
    }
  },
  "mappings":{
    "_source":{
      "excludes":[
        "custom"
      ]
    },
    "properties":{
      "apiVersion":{
        "type":"keyword"
      },
      "kind":{
        "type":"keyword",
        "normalizer":"lowercase"
      },
      "name":{
        "type":"keyword"
      },
      "namespace":{
        "type":"keyword"
      },
      "cluster":{
        "type":"keyword"
      },
      "object":{
        "properties":{
          "metadata":{
            "properties":{
              "annotations":{
                "type":"flattened"
              },
              "managedFields":{
                "type":"object",
                "enabled":false
              },
              "creationTimestamp":{
                "type":"date",
                "format":"yyyy-MM-dd'T'HH:mm:ss'Z'"
              },
              "deletionTimestamp":{
                "type":"date",
                "format":"yyyy-MM-dd'T'HH:mm:ss'Z'"
              },
              "labels":{
                "type":"flattened"
              },
              "name":{
                "type":"keyword"
              },
              "namespace":{
                "type":"keyword"
              },
              "ownerReferences":{
                "type":"flattened"
              },
              "resourceVersion":{
                "type":"keyword",
                "ignore_above":256
              }
            }
          },
          "spec":{
            "type":"flattened",
            "ignore_above":1024,
            "depth_limit":200
          }
        }
      }
    }
  }
}`
)

func createIndex(client *elasticsearch.Client, mapping string, indexName string) error {
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}
	resp, err := req.Do(context.TODO(), client)
	if err != nil {
		return err
	}
	if resp.IsError() {
		msg := resp.String()
		if strings.Contains(resp.String(), "resource_already_exists_exception") {
			return nil
		}
		return fmt.Errorf(msg)
	}
	return nil
}
