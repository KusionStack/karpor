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

package elasticsearch

const (
	defaultIndexName = "resources"
	defaultMapping   = `{
  "settings":{
    "index":{
      "max_result_window": "1000000",
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
      "cluster":{
        "type":"keyword"
      },
      "apiVersion":{
        "type":"keyword"
      },
      "kind":{
        "type":"keyword",
        "normalizer":"lowercase"
      },
      "namespace":{
        "type":"keyword"
      },
      "name":{
        "type":"keyword"
      },
      "labels":{
        "type":"flattened"
      },
      "annotations":{
        "type":"flattened"
      },
      "creationTimestamp":{
        "type":"date",
        "format":"yyyy-MM-dd'T'HH:mm:ss'Z'"
      },
      "deletionTimestamp":{
        "type":"date",
        "format":"yyyy-MM-dd'T'HH:mm:ss'Z'"
      },
      "ownerReferences":{
        "type":"flattened"
      },
      "resourceVersion":{
        "type":"keyword",
        "ignore_above":256
      },
      "content":{
        "type":"text"
      }
    }
  }
}`
)
