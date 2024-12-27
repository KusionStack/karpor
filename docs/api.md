


# Karpor
Karpor is a brand new Kubernetes visualization tool that focuses on search, insights, and AI at its core
  

## Informations

### Version

1.0

### Contact

  

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json
  * multipart/form-data
  * text/plain

### Produces
  * application/json
  * text/event-stream
  * text/plain

## All endpoints

###  authn

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /authn | [get authn](#get-authn) | Get returns an authn result of user's token. |
  


###  cluster

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /rest-api/v1/cluster/{clusterName} | [delete rest API v1 cluster cluster name](#delete-rest-api-v1-cluster-cluster-name) | Delete removes a cluster resource by name. |
| GET | /rest-api/v1/cluster/{clusterName} | [get rest API v1 cluster cluster name](#get-rest-api-v1-cluster-cluster-name) | Get returns a cluster resource by name. |
| GET | /rest-api/v1/clusters | [get rest API v1 clusters](#get-rest-api-v1-clusters) | List lists all cluster resources. |
| POST | /rest-api/v1/cluster/{clusterName} | [post rest API v1 cluster cluster name](#post-rest-api-v1-cluster-cluster-name) | Create creates a cluster resource. |
| POST | /rest-api/v1/cluster/config/file | [post rest API v1 cluster config file](#post-rest-api-v1-cluster-config-file) | Upload kubeConfig file for cluster |
| POST | /rest-api/v1/cluster/config/validate | [post rest API v1 cluster config validate](#post-rest-api-v1-cluster-config-validate) | Validate KubeConfig |
| PUT | /rest-api/v1/cluster/{clusterName} | [put rest API v1 cluster cluster name](#put-rest-api-v1-cluster-cluster-name) | Update updates the cluster metadata by name. |
  


###  debug

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /endpoints | [get endpoints](#get-endpoints) | List all available endpoints |
  


###  insight

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /insight/aggregator/event/{cluster}/{namespace}/{name} | [get insight aggregator event cluster namespace name](#get-insight-aggregator-event-cluster-namespace-name) | Stream resource events using Server-Sent Events |
| GET | /insight/aggregator/log/pod/{cluster}/{namespace}/{name} | [get insight aggregator log pod cluster namespace name](#get-insight-aggregator-log-pod-cluster-namespace-name) | Stream pod logs using Server-Sent Events |
| GET | /rest-api/v1/insight/audit | [get rest API v1 insight audit](#get-rest-api-v1-insight-audit) | Audit based on resource group. |
| GET | /rest-api/v1/insight/detail | [get rest API v1 insight detail](#get-rest-api-v1-insight-detail) | GetDetail returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind. |
| GET | /rest-api/v1/insight/events | [get rest API v1 insight events](#get-rest-api-v1-insight-events) | GetEvents returns events for a Kubernetes resource by name, namespace, cluster, apiVersion and kind. |
| GET | /rest-api/v1/insight/score | [get rest API v1 insight score](#get-rest-api-v1-insight-score) | ScoreHandler calculates a score for the audited manifest. |
| GET | /rest-api/v1/insight/stats | [get rest API v1 insight stats](#get-rest-api-v1-insight-stats) | Get returns a global statistics info. |
| GET | /rest-api/v1/insight/summary | [get rest API v1 insight summary](#get-rest-api-v1-insight-summary) | Get returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind. |
| GET | /rest-api/v1/insight/topology | [get rest API v1 insight topology](#get-rest-api-v1-insight-topology) | GetTopology returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind. |
| POST | /insight/aggregator/event/diagnosis/stream | [post insight aggregator event diagnosis stream](#post-insight-aggregator-event-diagnosis-stream) | Diagnose events using AI |
| POST | /insight/aggregator/log/diagnosis/stream | [post insight aggregator log diagnosis stream](#post-insight-aggregator-log-diagnosis-stream) | Diagnose pod logs using AI |
| POST | /insight/yaml/interpret/stream | [post insight yaml interpret stream](#post-insight-yaml-interpret-stream) | Interpret YAML using AI |
  


###  resourcegroup

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /rest-api/v1/resource-groups/{resourceGroupRuleName} | [get rest API v1 resource groups resource group rule name](#get-rest-api-v1-resource-groups-resource-group-rule-name) | List lists all ResourceGroups by rule name. |
  


###  resourcegrouprule

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /rest-api/v1/resource-group-rule/{resourceGroupRuleName} | [delete rest API v1 resource group rule resource group rule name](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name) | Delete removes a ResourceGroupRule by name. |
| GET | /rest-api/v1/resource-group-rule/{resourceGroupRuleName} | [get rest API v1 resource group rule resource group rule name](#get-rest-api-v1-resource-group-rule-resource-group-rule-name) | Get returns a ResourceGroupRule by name. |
| GET | /rest-api/v1/resource-group-rules | [get rest API v1 resource group rules](#get-rest-api-v1-resource-group-rules) | List lists all ResourceGroupRules. |
| POST | /rest-api/v1/resource-group-rule | [post rest API v1 resource group rule](#post-rest-api-v1-resource-group-rule) | Create creates a ResourceGroupRule. |
| PUT | /rest-api/v1/resource-group-rule | [put rest API v1 resource group rule](#put-rest-api-v1-resource-group-rule) | Update updates the ResourceGroupRule metadata by name. |
  


###  search

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /rest-api/v1/search | [get rest API v1 search](#get-rest-api-v1-search) | SearchForResource returns an array of Kubernetes runtime Object matched using the query from context. |
  


## Paths

### <span id="delete-rest-api-v1-cluster-cluster-name"></span> Delete removes a cluster resource by name. (*DeleteRestAPIV1ClusterClusterName*)

```
DELETE /rest-api/v1/cluster/{clusterName}
```

This endpoint deletes the cluster resource by name.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| clusterName | `path` | string | `string` |  | ✓ |  | The name of the cluster |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-rest-api-v1-cluster-cluster-name-200) | OK | Operation status |  | [schema](#delete-rest-api-v1-cluster-cluster-name-200-schema) |
| [400](#delete-rest-api-v1-cluster-cluster-name-400) | Bad Request | Bad Request |  | [schema](#delete-rest-api-v1-cluster-cluster-name-400-schema) |
| [401](#delete-rest-api-v1-cluster-cluster-name-401) | Unauthorized | Unauthorized |  | [schema](#delete-rest-api-v1-cluster-cluster-name-401-schema) |
| [404](#delete-rest-api-v1-cluster-cluster-name-404) | Not Found | Not Found |  | [schema](#delete-rest-api-v1-cluster-cluster-name-404-schema) |
| [405](#delete-rest-api-v1-cluster-cluster-name-405) | Method Not Allowed | Method Not Allowed |  | [schema](#delete-rest-api-v1-cluster-cluster-name-405-schema) |
| [429](#delete-rest-api-v1-cluster-cluster-name-429) | Too Many Requests | Too Many Requests |  | [schema](#delete-rest-api-v1-cluster-cluster-name-429-schema) |
| [500](#delete-rest-api-v1-cluster-cluster-name-500) | Internal Server Error | Internal Server Error |  | [schema](#delete-rest-api-v1-cluster-cluster-name-500-schema) |

#### Responses


##### <span id="delete-rest-api-v1-cluster-cluster-name-200"></span> 200 - Operation status
Status: OK

###### <span id="delete-rest-api-v1-cluster-cluster-name-200-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-cluster-cluster-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-rest-api-v1-cluster-cluster-name-400-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-cluster-cluster-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="delete-rest-api-v1-cluster-cluster-name-401-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-cluster-cluster-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-rest-api-v1-cluster-cluster-name-404-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-cluster-cluster-name-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="delete-rest-api-v1-cluster-cluster-name-405-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-cluster-cluster-name-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="delete-rest-api-v1-cluster-cluster-name-429-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-cluster-cluster-name-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="delete-rest-api-v1-cluster-cluster-name-500-schema"></span> Schema
   
  



### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name"></span> Delete removes a ResourceGroupRule by name. (*DeleteRestAPIV1ResourceGroupRuleResourceGroupRuleName*)

```
DELETE /rest-api/v1/resource-group-rule/{resourceGroupRuleName}
```

This endpoint deletes the ResourceGroupRule by name.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| resourceGroupRuleName | `path` | string | `string` |  | ✓ |  | The name of the resource group rule |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-200) | OK | Operation status |  | [schema](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-200-schema) |
| [400](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-400) | Bad Request | Bad Request |  | [schema](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-400-schema) |
| [401](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-401) | Unauthorized | Unauthorized |  | [schema](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-401-schema) |
| [404](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-404) | Not Found | Not Found |  | [schema](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-404-schema) |
| [405](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-405) | Method Not Allowed | Method Not Allowed |  | [schema](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-405-schema) |
| [429](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-429) | Too Many Requests | Too Many Requests |  | [schema](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-429-schema) |
| [500](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-500) | Internal Server Error | Internal Server Error |  | [schema](#delete-rest-api-v1-resource-group-rule-resource-group-rule-name-500-schema) |

#### Responses


##### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-200"></span> 200 - Operation status
Status: OK

###### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-200-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-400-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-401-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-404-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-405-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-429-schema"></span> Schema
   
  



##### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="delete-rest-api-v1-resource-group-rule-resource-group-rule-name-500-schema"></span> Schema
   
  



### <span id="get-authn"></span> Get returns an authn result of user's token. (*GetAuthn*)

```
GET /authn
```

This endpoint returns an authn result.

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-authn-200) | OK | OK |  | [schema](#get-authn-200-schema) |
| [400](#get-authn-400) | Bad Request | Bad Request |  | [schema](#get-authn-400-schema) |
| [401](#get-authn-401) | Unauthorized | Unauthorized |  | [schema](#get-authn-401-schema) |
| [404](#get-authn-404) | Not Found | Not Found |  | [schema](#get-authn-404-schema) |
| [405](#get-authn-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-authn-405-schema) |
| [429](#get-authn-429) | Too Many Requests | Too Many Requests |  | [schema](#get-authn-429-schema) |
| [500](#get-authn-500) | Internal Server Error | Internal Server Error |  | [schema](#get-authn-500-schema) |

#### Responses


##### <span id="get-authn-200"></span> 200 - OK
Status: OK

###### <span id="get-authn-200-schema"></span> Schema
   
  



##### <span id="get-authn-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-authn-400-schema"></span> Schema
   
  



##### <span id="get-authn-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-authn-401-schema"></span> Schema
   
  



##### <span id="get-authn-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-authn-404-schema"></span> Schema
   
  



##### <span id="get-authn-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-authn-405-schema"></span> Schema
   
  



##### <span id="get-authn-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-authn-429-schema"></span> Schema
   
  



##### <span id="get-authn-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-authn-500-schema"></span> Schema
   
  



### <span id="get-endpoints"></span> List all available endpoints (*GetEndpoints*)

```
GET /endpoints
```

List all registered endpoints in the router

#### Consumes
  * text/plain

#### Produces
  * text/plain

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-endpoints-200) | OK | Endpoints listed successfully |  | [schema](#get-endpoints-200-schema) |

#### Responses


##### <span id="get-endpoints-200"></span> 200 - Endpoints listed successfully
Status: OK

###### <span id="get-endpoints-200-schema"></span> Schema
   
  



### <span id="get-insight-aggregator-event-cluster-namespace-name"></span> Stream resource events using Server-Sent Events (*GetInsightAggregatorEventClusterNamespaceName*)

```
GET /insight/aggregator/event/{cluster}/{namespace}/{name}
```

This endpoint streams resource events in real-time using SSE. It supports event type filtering and automatic updates.

#### Produces
  * text/event-stream

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cluster | `path` | string | `string` |  | ✓ |  | The cluster name |
| name | `path` | string | `string` |  | ✓ |  | The resource name |
| namespace | `path` | string | `string` |  | ✓ |  | The namespace name |
| apiVersion | `query` | string | `string` |  | ✓ |  | The resource API version |
| kind | `query` | string | `string` |  | ✓ |  | The resource kind |
| type | `query` | string | `string` |  |  |  | Event type filter (Normal or Warning) |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-insight-aggregator-event-cluster-namespace-name-200) | OK | OK |  | [schema](#get-insight-aggregator-event-cluster-namespace-name-200-schema) |
| [400](#get-insight-aggregator-event-cluster-namespace-name-400) | Bad Request | Bad Request |  | [schema](#get-insight-aggregator-event-cluster-namespace-name-400-schema) |
| [401](#get-insight-aggregator-event-cluster-namespace-name-401) | Unauthorized | Unauthorized |  | [schema](#get-insight-aggregator-event-cluster-namespace-name-401-schema) |
| [404](#get-insight-aggregator-event-cluster-namespace-name-404) | Not Found | Not Found |  | [schema](#get-insight-aggregator-event-cluster-namespace-name-404-schema) |

#### Responses


##### <span id="get-insight-aggregator-event-cluster-namespace-name-200"></span> 200 - OK
Status: OK

###### <span id="get-insight-aggregator-event-cluster-namespace-name-200-schema"></span> Schema
   
  

[][AiEvent](#ai-event)

##### <span id="get-insight-aggregator-event-cluster-namespace-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-insight-aggregator-event-cluster-namespace-name-400-schema"></span> Schema
   
  



##### <span id="get-insight-aggregator-event-cluster-namespace-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-insight-aggregator-event-cluster-namespace-name-401-schema"></span> Schema
   
  



##### <span id="get-insight-aggregator-event-cluster-namespace-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-insight-aggregator-event-cluster-namespace-name-404-schema"></span> Schema
   
  



### <span id="get-insight-aggregator-log-pod-cluster-namespace-name"></span> Stream pod logs using Server-Sent Events (*GetInsightAggregatorLogPodClusterNamespaceName*)

```
GET /insight/aggregator/log/pod/{cluster}/{namespace}/{name}
```

This endpoint streams pod logs in real-time using SSE. It supports container selection and automatic reconnection.

#### Produces
  * application/json
  * text/event-stream

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cluster | `path` | string | `string` |  | ✓ |  | The cluster name |
| name | `path` | string | `string` |  | ✓ |  | The pod name |
| namespace | `path` | string | `string` |  | ✓ |  | The namespace name |
| container | `query` | string | `string` |  |  |  | The container name (optional if pod has only one container) |
| download | `query` | boolean | `bool` |  |  |  | Download logs as file instead of streaming |
| since | `query` | string | `string` |  |  |  | Only return logs newer than a relative duration like 5s, 2m, or 3h |
| sinceTime | `query` | string | `string` |  |  |  | Only return logs after a specific date (RFC3339) |
| tailLines | `query` | integer | `int64` |  |  |  | Number of lines from the end of the logs to show |
| timestamps | `query` | boolean | `bool` |  |  |  | Include timestamps in log output |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-insight-aggregator-log-pod-cluster-namespace-name-200) | OK | OK |  | [schema](#get-insight-aggregator-log-pod-cluster-namespace-name-200-schema) |
| [400](#get-insight-aggregator-log-pod-cluster-namespace-name-400) | Bad Request | Bad Request |  | [schema](#get-insight-aggregator-log-pod-cluster-namespace-name-400-schema) |
| [401](#get-insight-aggregator-log-pod-cluster-namespace-name-401) | Unauthorized | Unauthorized |  | [schema](#get-insight-aggregator-log-pod-cluster-namespace-name-401-schema) |
| [404](#get-insight-aggregator-log-pod-cluster-namespace-name-404) | Not Found | Not Found |  | [schema](#get-insight-aggregator-log-pod-cluster-namespace-name-404-schema) |

#### Responses


##### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-200"></span> 200 - OK
Status: OK

###### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-200-schema"></span> Schema
   
  

[AggregatorLogEntry](#aggregator-log-entry)

##### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-400-schema"></span> Schema
   
  



##### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-401-schema"></span> Schema
   
  



##### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-insight-aggregator-log-pod-cluster-namespace-name-404-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-cluster-cluster-name"></span> Get returns a cluster resource by name. (*GetRestAPIV1ClusterClusterName*)

```
GET /rest-api/v1/cluster/{clusterName}
```

This endpoint returns a cluster resource by name.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| clusterName | `path` | string | `string` |  | ✓ |  | The name of the cluster |
| format | `query` | string | `string` |  |  |  | The format of the response. Either in json or yaml |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-cluster-cluster-name-200) | OK | Unstructured object |  | [schema](#get-rest-api-v1-cluster-cluster-name-200-schema) |
| [400](#get-rest-api-v1-cluster-cluster-name-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-cluster-cluster-name-400-schema) |
| [401](#get-rest-api-v1-cluster-cluster-name-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-cluster-cluster-name-401-schema) |
| [404](#get-rest-api-v1-cluster-cluster-name-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-cluster-cluster-name-404-schema) |
| [405](#get-rest-api-v1-cluster-cluster-name-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-cluster-cluster-name-405-schema) |
| [429](#get-rest-api-v1-cluster-cluster-name-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-cluster-cluster-name-429-schema) |
| [500](#get-rest-api-v1-cluster-cluster-name-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-cluster-cluster-name-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-cluster-cluster-name-200"></span> 200 - Unstructured object
Status: OK

###### <span id="get-rest-api-v1-cluster-cluster-name-200-schema"></span> Schema
   
  

[UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="get-rest-api-v1-cluster-cluster-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-cluster-cluster-name-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-cluster-cluster-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-cluster-cluster-name-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-cluster-cluster-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-cluster-cluster-name-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-cluster-cluster-name-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-cluster-cluster-name-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-cluster-cluster-name-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-cluster-cluster-name-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-cluster-cluster-name-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-cluster-cluster-name-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-clusters"></span> List lists all cluster resources. (*GetRestAPIV1Clusters*)

```
GET /rest-api/v1/clusters
```

This endpoint lists all cluster resources.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| descending | `query` | boolean | `bool` |  |  |  | Whether to sort the list in descending order. Default to false |
| orderBy | `query` | string | `string` |  |  |  | The order to list the cluster. Default to order by name |
| summary | `query` | boolean | `bool` |  |  |  | Whether to display summary or not. Default to false |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-clusters-200) | OK | List of cluster objects |  | [schema](#get-rest-api-v1-clusters-200-schema) |
| [400](#get-rest-api-v1-clusters-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-clusters-400-schema) |
| [401](#get-rest-api-v1-clusters-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-clusters-401-schema) |
| [404](#get-rest-api-v1-clusters-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-clusters-404-schema) |
| [405](#get-rest-api-v1-clusters-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-clusters-405-schema) |
| [429](#get-rest-api-v1-clusters-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-clusters-429-schema) |
| [500](#get-rest-api-v1-clusters-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-clusters-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-clusters-200"></span> 200 - List of cluster objects
Status: OK

###### <span id="get-rest-api-v1-clusters-200-schema"></span> Schema
   
  

[][UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="get-rest-api-v1-clusters-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-clusters-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-clusters-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-clusters-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-clusters-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-clusters-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-clusters-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-clusters-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-clusters-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-clusters-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-clusters-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-clusters-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-insight-audit"></span> Audit based on resource group. (*GetRestAPIV1InsightAudit*)

```
GET /rest-api/v1/insight/audit
```

This endpoint audits based on the specified resource group.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiVersion | `query` | string | `string` |  |  |  | The specified apiVersion, such as 'apps/v1' |
| cluster | `query` | string | `string` |  |  |  | The specified cluster name, such as 'example-cluster' |
| forceNew | `query` | boolean | `bool` |  |  |  | Switch for forced scanning, default is 'false' |
| kind | `query` | string | `string` |  |  |  | The specified kind, such as 'Deployment' |
| name | `query` | string | `string` |  |  |  | The specified resource name, such as 'foo' |
| namespace | `query` | string | `string` |  |  |  | The specified namespace, such as 'default' |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-insight-audit-200) | OK | Audit results |  | [schema](#get-rest-api-v1-insight-audit-200-schema) |
| [400](#get-rest-api-v1-insight-audit-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-insight-audit-400-schema) |
| [401](#get-rest-api-v1-insight-audit-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-insight-audit-401-schema) |
| [404](#get-rest-api-v1-insight-audit-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-insight-audit-404-schema) |
| [429](#get-rest-api-v1-insight-audit-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-insight-audit-429-schema) |
| [500](#get-rest-api-v1-insight-audit-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-insight-audit-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-insight-audit-200"></span> 200 - Audit results
Status: OK

###### <span id="get-rest-api-v1-insight-audit-200-schema"></span> Schema
   
  

[ScannerAuditData](#scanner-audit-data)

##### <span id="get-rest-api-v1-insight-audit-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-insight-audit-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-audit-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-insight-audit-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-audit-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-insight-audit-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-audit-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-insight-audit-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-audit-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-insight-audit-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-insight-detail"></span> GetDetail returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind. (*GetRestAPIV1InsightDetail*)

```
GET /rest-api/v1/insight/detail
```

This endpoint returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiVersion | `query` | string | `string` |  |  |  | The specified apiVersion, such as 'apps/v1'. Should be percent-encoded |
| cluster | `query` | string | `string` |  |  |  | The specified cluster name, such as 'example-cluster' |
| format | `query` | string | `string` |  |  |  | The format of the response. Either in json or yaml. Default to json |
| kind | `query` | string | `string` |  |  |  | The specified kind, such as 'Deployment' |
| name | `query` | string | `string` |  |  |  | The specified resource name, such as 'foo' |
| namespace | `query` | string | `string` |  |  |  | The specified namespace, such as 'default' |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-insight-detail-200) | OK | Unstructured object |  | [schema](#get-rest-api-v1-insight-detail-200-schema) |
| [400](#get-rest-api-v1-insight-detail-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-insight-detail-400-schema) |
| [401](#get-rest-api-v1-insight-detail-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-insight-detail-401-schema) |
| [404](#get-rest-api-v1-insight-detail-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-insight-detail-404-schema) |
| [405](#get-rest-api-v1-insight-detail-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-insight-detail-405-schema) |
| [429](#get-rest-api-v1-insight-detail-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-insight-detail-429-schema) |
| [500](#get-rest-api-v1-insight-detail-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-insight-detail-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-insight-detail-200"></span> 200 - Unstructured object
Status: OK

###### <span id="get-rest-api-v1-insight-detail-200-schema"></span> Schema
   
  

[UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="get-rest-api-v1-insight-detail-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-insight-detail-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-detail-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-insight-detail-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-detail-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-insight-detail-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-detail-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-insight-detail-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-detail-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-insight-detail-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-detail-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-insight-detail-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-insight-events"></span> GetEvents returns events for a Kubernetes resource by name, namespace, cluster, apiVersion and kind. (*GetRestAPIV1InsightEvents*)

```
GET /rest-api/v1/insight/events
```

This endpoint returns events for a Kubernetes resource YAML by name, namespace, cluster, apiVersion and kind.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiVersion | `query` | string | `string` |  |  |  | The specified apiVersion, such as 'apps/v1'. Should be percent-encoded |
| cluster | `query` | string | `string` |  |  |  | The specified cluster name, such as 'example-cluster' |
| kind | `query` | string | `string` |  |  |  | The specified kind, such as 'Deployment' |
| name | `query` | string | `string` |  |  |  | The specified resource name, such as 'foo' |
| namespace | `query` | string | `string` |  |  |  | The specified namespace, such as 'default' |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-insight-events-200) | OK | List of events |  | [schema](#get-rest-api-v1-insight-events-200-schema) |
| [400](#get-rest-api-v1-insight-events-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-insight-events-400-schema) |
| [401](#get-rest-api-v1-insight-events-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-insight-events-401-schema) |
| [404](#get-rest-api-v1-insight-events-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-insight-events-404-schema) |
| [405](#get-rest-api-v1-insight-events-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-insight-events-405-schema) |
| [429](#get-rest-api-v1-insight-events-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-insight-events-429-schema) |
| [500](#get-rest-api-v1-insight-events-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-insight-events-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-insight-events-200"></span> 200 - List of events
Status: OK

###### <span id="get-rest-api-v1-insight-events-200-schema"></span> Schema
   
  

[][UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="get-rest-api-v1-insight-events-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-insight-events-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-events-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-insight-events-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-events-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-insight-events-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-events-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-insight-events-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-events-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-insight-events-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-events-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-insight-events-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-insight-score"></span> ScoreHandler calculates a score for the audited manifest. (*GetRestAPIV1InsightScore*)

```
GET /rest-api/v1/insight/score
```

This endpoint calculates a score for the provided manifest based on the number and severity of issues detected during the audit.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiVersion | `query` | string | `string` |  |  |  | The specified apiVersion, such as 'apps/v1' |
| cluster | `query` | string | `string` |  |  |  | The specified cluster name, such as 'example-cluster' |
| forceNew | `query` | boolean | `bool` |  |  |  | Switch for forced compute score, default is 'false' |
| kind | `query` | string | `string` |  |  |  | The specified kind, such as 'Deployment' |
| name | `query` | string | `string` |  |  |  | The specified resource name, such as 'foo' |
| namespace | `query` | string | `string` |  |  |  | The specified namespace, such as 'default' |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-insight-score-200) | OK | Score calculation result |  | [schema](#get-rest-api-v1-insight-score-200-schema) |
| [400](#get-rest-api-v1-insight-score-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-insight-score-400-schema) |
| [401](#get-rest-api-v1-insight-score-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-insight-score-401-schema) |
| [404](#get-rest-api-v1-insight-score-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-insight-score-404-schema) |
| [429](#get-rest-api-v1-insight-score-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-insight-score-429-schema) |
| [500](#get-rest-api-v1-insight-score-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-insight-score-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-insight-score-200"></span> 200 - Score calculation result
Status: OK

###### <span id="get-rest-api-v1-insight-score-200-schema"></span> Schema
   
  

[InsightScoreData](#insight-score-data)

##### <span id="get-rest-api-v1-insight-score-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-insight-score-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-score-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-insight-score-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-score-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-insight-score-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-score-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-insight-score-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-score-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-insight-score-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-insight-stats"></span> Get returns a global statistics info. (*GetRestAPIV1InsightStats*)

```
GET /rest-api/v1/insight/stats
```

This endpoint returns a global statistics info.

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-insight-stats-200) | OK | Global statistics info |  | [schema](#get-rest-api-v1-insight-stats-200-schema) |
| [400](#get-rest-api-v1-insight-stats-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-insight-stats-400-schema) |
| [401](#get-rest-api-v1-insight-stats-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-insight-stats-401-schema) |
| [404](#get-rest-api-v1-insight-stats-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-insight-stats-404-schema) |
| [405](#get-rest-api-v1-insight-stats-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-insight-stats-405-schema) |
| [429](#get-rest-api-v1-insight-stats-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-insight-stats-429-schema) |
| [500](#get-rest-api-v1-insight-stats-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-insight-stats-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-insight-stats-200"></span> 200 - Global statistics info
Status: OK

###### <span id="get-rest-api-v1-insight-stats-200-schema"></span> Schema
   
  

[InsightStatistics](#insight-statistics)

##### <span id="get-rest-api-v1-insight-stats-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-insight-stats-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-stats-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-insight-stats-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-stats-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-insight-stats-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-stats-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-insight-stats-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-stats-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-insight-stats-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-stats-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-insight-stats-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-insight-summary"></span> Get returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind. (*GetRestAPIV1InsightSummary*)

```
GET /rest-api/v1/insight/summary
```

This endpoint returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiVersion | `query` | string | `string` |  |  |  | The specified apiVersion, such as 'apps/v1'. Should be percent-encoded |
| cluster | `query` | string | `string` |  |  |  | The specified cluster name, such as 'example-cluster' |
| kind | `query` | string | `string` |  |  |  | The specified kind, such as 'Deployment' |
| name | `query` | string | `string` |  |  |  | The specified resource name, such as 'foo' |
| namespace | `query` | string | `string` |  |  |  | The specified namespace, such as 'default' |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-insight-summary-200) | OK | Resource Summary |  | [schema](#get-rest-api-v1-insight-summary-200-schema) |
| [400](#get-rest-api-v1-insight-summary-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-insight-summary-400-schema) |
| [401](#get-rest-api-v1-insight-summary-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-insight-summary-401-schema) |
| [404](#get-rest-api-v1-insight-summary-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-insight-summary-404-schema) |
| [405](#get-rest-api-v1-insight-summary-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-insight-summary-405-schema) |
| [429](#get-rest-api-v1-insight-summary-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-insight-summary-429-schema) |
| [500](#get-rest-api-v1-insight-summary-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-insight-summary-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-insight-summary-200"></span> 200 - Resource Summary
Status: OK

###### <span id="get-rest-api-v1-insight-summary-200-schema"></span> Schema
   
  

[InsightResourceSummary](#insight-resource-summary)

##### <span id="get-rest-api-v1-insight-summary-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-insight-summary-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-summary-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-insight-summary-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-summary-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-insight-summary-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-summary-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-insight-summary-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-summary-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-insight-summary-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-summary-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-insight-summary-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-insight-topology"></span> GetTopology returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind. (*GetRestAPIV1InsightTopology*)

```
GET /rest-api/v1/insight/topology
```

This endpoint returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| apiVersion | `query` | string | `string` |  |  |  | The specified apiVersion, such as 'apps/v1'. Should be percent-encoded |
| cluster | `query` | string | `string` |  |  |  | The specified cluster name, such as 'example-cluster' |
| forceNew | `query` | boolean | `bool` |  |  |  | Force re-generating the topology, default is 'false' |
| kind | `query` | string | `string` |  |  |  | The specified kind, such as 'Deployment' |
| name | `query` | string | `string` |  |  |  | The specified resource name, such as 'foo' |
| namespace | `query` | string | `string` |  |  |  | The specified namespace, such as 'default' |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-insight-topology-200) | OK | map from string to resource.ResourceTopology |  | [schema](#get-rest-api-v1-insight-topology-200-schema) |
| [400](#get-rest-api-v1-insight-topology-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-insight-topology-400-schema) |
| [401](#get-rest-api-v1-insight-topology-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-insight-topology-401-schema) |
| [404](#get-rest-api-v1-insight-topology-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-insight-topology-404-schema) |
| [405](#get-rest-api-v1-insight-topology-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-insight-topology-405-schema) |
| [429](#get-rest-api-v1-insight-topology-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-insight-topology-429-schema) |
| [500](#get-rest-api-v1-insight-topology-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-insight-topology-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-insight-topology-200"></span> 200 - map from string to resource.ResourceTopology
Status: OK

###### <span id="get-rest-api-v1-insight-topology-200-schema"></span> Schema
   
  

map of [InsightResourceTopology](#insight-resource-topology)

##### <span id="get-rest-api-v1-insight-topology-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-insight-topology-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-topology-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-insight-topology-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-topology-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-insight-topology-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-topology-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-insight-topology-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-topology-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-insight-topology-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-insight-topology-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-insight-topology-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name"></span> Get returns a ResourceGroupRule by name. (*GetRestAPIV1ResourceGroupRuleResourceGroupRuleName*)

```
GET /rest-api/v1/resource-group-rule/{resourceGroupRuleName}
```

This endpoint returns a ResourceGroupRule by name.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| resourceGroupRuleName | `path` | string | `string` |  | ✓ |  | The name of the resource group rule |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-200) | OK | Unstructured object |  | [schema](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-200-schema) |
| [400](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-400-schema) |
| [401](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-401-schema) |
| [404](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-404-schema) |
| [405](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-405-schema) |
| [429](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-429-schema) |
| [500](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-resource-group-rule-resource-group-rule-name-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-200"></span> 200 - Unstructured object
Status: OK

###### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-200-schema"></span> Schema
   
  

[UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-resource-group-rule-resource-group-rule-name-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-resource-group-rules"></span> List lists all ResourceGroupRules. (*GetRestAPIV1ResourceGroupRules*)

```
GET /rest-api/v1/resource-group-rules
```

This endpoint lists all ResourceGroupRules.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| descending | `query` | boolean | `bool` |  |  |  | Whether to sort the list in descending order. Default to false |
| orderBy | `query` | string | `string` |  |  |  | The order to list the resourceGroupRule. Default to order by name |
| summary | `query` | boolean | `bool` |  |  |  | Whether to display summary or not. Default to false |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-resource-group-rules-200) | OK | List of resourceGroupRule objects |  | [schema](#get-rest-api-v1-resource-group-rules-200-schema) |
| [400](#get-rest-api-v1-resource-group-rules-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-resource-group-rules-400-schema) |
| [401](#get-rest-api-v1-resource-group-rules-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-resource-group-rules-401-schema) |
| [404](#get-rest-api-v1-resource-group-rules-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-resource-group-rules-404-schema) |
| [405](#get-rest-api-v1-resource-group-rules-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-resource-group-rules-405-schema) |
| [429](#get-rest-api-v1-resource-group-rules-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-resource-group-rules-429-schema) |
| [500](#get-rest-api-v1-resource-group-rules-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-resource-group-rules-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-resource-group-rules-200"></span> 200 - List of resourceGroupRule objects
Status: OK

###### <span id="get-rest-api-v1-resource-group-rules-200-schema"></span> Schema
   
  

[][UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="get-rest-api-v1-resource-group-rules-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-resource-group-rules-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rules-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-resource-group-rules-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rules-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-resource-group-rules-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rules-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-resource-group-rules-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rules-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-resource-group-rules-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-group-rules-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-resource-group-rules-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name"></span> List lists all ResourceGroups by rule name. (*GetRestAPIV1ResourceGroupsResourceGroupRuleName*)

```
GET /rest-api/v1/resource-groups/{resourceGroupRuleName}
```

This endpoint lists all ResourceGroups.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| resourceGroupRuleName | `path` | string | `string` |  | ✓ |  | The name of the resource group rule |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-resource-groups-resource-group-rule-name-200) | OK | List of resourceGroup objects |  | [schema](#get-rest-api-v1-resource-groups-resource-group-rule-name-200-schema) |
| [400](#get-rest-api-v1-resource-groups-resource-group-rule-name-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-resource-groups-resource-group-rule-name-400-schema) |
| [401](#get-rest-api-v1-resource-groups-resource-group-rule-name-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-resource-groups-resource-group-rule-name-401-schema) |
| [404](#get-rest-api-v1-resource-groups-resource-group-rule-name-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-resource-groups-resource-group-rule-name-404-schema) |
| [405](#get-rest-api-v1-resource-groups-resource-group-rule-name-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-resource-groups-resource-group-rule-name-405-schema) |
| [429](#get-rest-api-v1-resource-groups-resource-group-rule-name-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-resource-groups-resource-group-rule-name-429-schema) |
| [500](#get-rest-api-v1-resource-groups-resource-group-rule-name-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-resource-groups-resource-group-rule-name-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-200"></span> 200 - List of resourceGroup objects
Status: OK

###### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-200-schema"></span> Schema
   
  

[][UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-resource-groups-resource-group-rule-name-500-schema"></span> Schema
   
  



### <span id="get-rest-api-v1-search"></span> SearchForResource returns an array of Kubernetes runtime Object matched using the query from context. (*GetRestAPIV1Search*)

```
GET /rest-api/v1/search
```

This endpoint returns an array of Kubernetes runtime Object matched using the query from context.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| page | `query` | string | `string` |  |  |  | The current page to fetch. Default to 1 |
| pageSize | `query` | string | `string` |  |  |  | The size of the page. Default to 10 |
| pattern | `query` | string | `string` |  | ✓ |  | The search pattern. Can be either sql, dsl or nl. Required |
| query | `query` | string | `string` |  | ✓ |  | The query to use for search. Required |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-rest-api-v1-search-200) | OK | Array of runtime.Object |  | [schema](#get-rest-api-v1-search-200-schema) |
| [400](#get-rest-api-v1-search-400) | Bad Request | Bad Request |  | [schema](#get-rest-api-v1-search-400-schema) |
| [401](#get-rest-api-v1-search-401) | Unauthorized | Unauthorized |  | [schema](#get-rest-api-v1-search-401-schema) |
| [404](#get-rest-api-v1-search-404) | Not Found | Not Found |  | [schema](#get-rest-api-v1-search-404-schema) |
| [405](#get-rest-api-v1-search-405) | Method Not Allowed | Method Not Allowed |  | [schema](#get-rest-api-v1-search-405-schema) |
| [429](#get-rest-api-v1-search-429) | Too Many Requests | Too Many Requests |  | [schema](#get-rest-api-v1-search-429-schema) |
| [500](#get-rest-api-v1-search-500) | Internal Server Error | Internal Server Error |  | [schema](#get-rest-api-v1-search-500-schema) |

#### Responses


##### <span id="get-rest-api-v1-search-200"></span> 200 - Array of runtime.Object
Status: OK

###### <span id="get-rest-api-v1-search-200-schema"></span> Schema
   
  

[][interface{}](#interface)

##### <span id="get-rest-api-v1-search-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-rest-api-v1-search-400-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-search-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="get-rest-api-v1-search-401-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-search-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-rest-api-v1-search-404-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-search-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="get-rest-api-v1-search-405-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-search-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="get-rest-api-v1-search-429-schema"></span> Schema
   
  



##### <span id="get-rest-api-v1-search-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-rest-api-v1-search-500-schema"></span> Schema
   
  



### <span id="post-insight-aggregator-event-diagnosis-stream"></span> Diagnose events using AI (*PostInsightAggregatorEventDiagnosisStream*)

```
POST /insight/aggregator/event/diagnosis/stream
```

This endpoint analyzes events using AI to identify issues and provide solutions

#### Consumes
  * application/json

#### Produces
  * text/event-stream

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [AggregatorEventDiagnoseRequest](#aggregator-event-diagnose-request) | `models.AggregatorEventDiagnoseRequest` | | ✓ | | The events to analyze |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-insight-aggregator-event-diagnosis-stream-200) | OK | OK |  | [schema](#post-insight-aggregator-event-diagnosis-stream-200-schema) |
| [400](#post-insight-aggregator-event-diagnosis-stream-400) | Bad Request | Bad Request |  | [schema](#post-insight-aggregator-event-diagnosis-stream-400-schema) |
| [500](#post-insight-aggregator-event-diagnosis-stream-500) | Internal Server Error | Internal Server Error |  | [schema](#post-insight-aggregator-event-diagnosis-stream-500-schema) |

#### Responses


##### <span id="post-insight-aggregator-event-diagnosis-stream-200"></span> 200 - OK
Status: OK

###### <span id="post-insight-aggregator-event-diagnosis-stream-200-schema"></span> Schema
   
  

[AiDiagnosisEvent](#ai-diagnosis-event)

##### <span id="post-insight-aggregator-event-diagnosis-stream-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-insight-aggregator-event-diagnosis-stream-400-schema"></span> Schema
   
  



##### <span id="post-insight-aggregator-event-diagnosis-stream-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-insight-aggregator-event-diagnosis-stream-500-schema"></span> Schema
   
  



### <span id="post-insight-aggregator-log-diagnosis-stream"></span> Diagnose pod logs using AI (*PostInsightAggregatorLogDiagnosisStream*)

```
POST /insight/aggregator/log/diagnosis/stream
```

This endpoint analyzes pod logs using AI to identify issues and provide solutions

#### Consumes
  * application/json

#### Produces
  * text/event-stream

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [AggregatorDiagnoseRequest](#aggregator-diagnose-request) | `models.AggregatorDiagnoseRequest` | | ✓ | | The logs to analyze |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-insight-aggregator-log-diagnosis-stream-200) | OK | OK |  | [schema](#post-insight-aggregator-log-diagnosis-stream-200-schema) |
| [400](#post-insight-aggregator-log-diagnosis-stream-400) | Bad Request | Bad Request |  | [schema](#post-insight-aggregator-log-diagnosis-stream-400-schema) |
| [500](#post-insight-aggregator-log-diagnosis-stream-500) | Internal Server Error | Internal Server Error |  | [schema](#post-insight-aggregator-log-diagnosis-stream-500-schema) |

#### Responses


##### <span id="post-insight-aggregator-log-diagnosis-stream-200"></span> 200 - OK
Status: OK

###### <span id="post-insight-aggregator-log-diagnosis-stream-200-schema"></span> Schema
   
  

[AiDiagnosisEvent](#ai-diagnosis-event)

##### <span id="post-insight-aggregator-log-diagnosis-stream-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-insight-aggregator-log-diagnosis-stream-400-schema"></span> Schema
   
  



##### <span id="post-insight-aggregator-log-diagnosis-stream-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-insight-aggregator-log-diagnosis-stream-500-schema"></span> Schema
   
  



### <span id="post-insight-yaml-interpret-stream"></span> Interpret YAML using AI (*PostInsightYamlInterpretStream*)

```
POST /insight/yaml/interpret/stream
```

This endpoint analyzes YAML content using AI to provide detailed interpretation and insights

#### Consumes
  * application/json

#### Produces
  * text/event-stream

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [DetailInterpretRequest](#detail-interpret-request) | `models.DetailInterpretRequest` | | ✓ | | The YAML content to interpret |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-insight-yaml-interpret-stream-200) | OK | OK |  | [schema](#post-insight-yaml-interpret-stream-200-schema) |
| [400](#post-insight-yaml-interpret-stream-400) | Bad Request | Bad Request |  | [schema](#post-insight-yaml-interpret-stream-400-schema) |
| [500](#post-insight-yaml-interpret-stream-500) | Internal Server Error | Internal Server Error |  | [schema](#post-insight-yaml-interpret-stream-500-schema) |

#### Responses


##### <span id="post-insight-yaml-interpret-stream-200"></span> 200 - OK
Status: OK

###### <span id="post-insight-yaml-interpret-stream-200-schema"></span> Schema
   
  

[AiInterpretEvent](#ai-interpret-event)

##### <span id="post-insight-yaml-interpret-stream-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-insight-yaml-interpret-stream-400-schema"></span> Schema
   
  



##### <span id="post-insight-yaml-interpret-stream-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-insight-yaml-interpret-stream-500-schema"></span> Schema
   
  



### <span id="post-rest-api-v1-cluster-cluster-name"></span> Create creates a cluster resource. (*PostRestAPIV1ClusterClusterName*)

```
POST /rest-api/v1/cluster/{clusterName}
```

This endpoint creates a new cluster resource using the payload.

#### Consumes
  * application/json
  * text/plain

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| clusterName | `path` | string | `string` |  | ✓ |  | The name of the cluster |
| request | `body` | [ClusterClusterPayload](#cluster-cluster-payload) | `models.ClusterClusterPayload` | | ✓ | | cluster to create (either plain text or JSON format) |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-rest-api-v1-cluster-cluster-name-200) | OK | Unstructured object |  | [schema](#post-rest-api-v1-cluster-cluster-name-200-schema) |
| [400](#post-rest-api-v1-cluster-cluster-name-400) | Bad Request | Bad Request |  | [schema](#post-rest-api-v1-cluster-cluster-name-400-schema) |
| [401](#post-rest-api-v1-cluster-cluster-name-401) | Unauthorized | Unauthorized |  | [schema](#post-rest-api-v1-cluster-cluster-name-401-schema) |
| [404](#post-rest-api-v1-cluster-cluster-name-404) | Not Found | Not Found |  | [schema](#post-rest-api-v1-cluster-cluster-name-404-schema) |
| [405](#post-rest-api-v1-cluster-cluster-name-405) | Method Not Allowed | Method Not Allowed |  | [schema](#post-rest-api-v1-cluster-cluster-name-405-schema) |
| [429](#post-rest-api-v1-cluster-cluster-name-429) | Too Many Requests | Too Many Requests |  | [schema](#post-rest-api-v1-cluster-cluster-name-429-schema) |
| [500](#post-rest-api-v1-cluster-cluster-name-500) | Internal Server Error | Internal Server Error |  | [schema](#post-rest-api-v1-cluster-cluster-name-500-schema) |

#### Responses


##### <span id="post-rest-api-v1-cluster-cluster-name-200"></span> 200 - Unstructured object
Status: OK

###### <span id="post-rest-api-v1-cluster-cluster-name-200-schema"></span> Schema
   
  

[UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="post-rest-api-v1-cluster-cluster-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-rest-api-v1-cluster-cluster-name-400-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-cluster-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="post-rest-api-v1-cluster-cluster-name-401-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-cluster-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="post-rest-api-v1-cluster-cluster-name-404-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-cluster-name-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="post-rest-api-v1-cluster-cluster-name-405-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-cluster-name-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="post-rest-api-v1-cluster-cluster-name-429-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-cluster-name-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-rest-api-v1-cluster-cluster-name-500-schema"></span> Schema
   
  



### <span id="post-rest-api-v1-cluster-config-file"></span> Upload kubeConfig file for cluster (*PostRestAPIV1ClusterConfigFile*)

```
POST /rest-api/v1/cluster/config/file
```

Uploads a KubeConfig file for cluster, with a maximum size of 2MB.

#### Consumes
  * multipart/form-data

#### Produces
  * text/plain

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| description | `formData` | string | `string` |  | ✓ |  | cluster description |
| displayName | `formData` | string | `string` |  | ✓ |  | cluster display name |
| file | `formData` | file | `io.ReadCloser` |  | ✓ |  | Upload file with field name 'file' |
| name | `formData` | string | `string` |  | ✓ |  | cluster name |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-rest-api-v1-cluster-config-file-200) | OK | Returns the content of the uploaded KubeConfig file. |  | [schema](#post-rest-api-v1-cluster-config-file-200-schema) |
| [400](#post-rest-api-v1-cluster-config-file-400) | Bad Request | The uploaded file is too large or the request is invalid. |  | [schema](#post-rest-api-v1-cluster-config-file-400-schema) |
| [500](#post-rest-api-v1-cluster-config-file-500) | Internal Server Error | Internal server error. |  | [schema](#post-rest-api-v1-cluster-config-file-500-schema) |

#### Responses


##### <span id="post-rest-api-v1-cluster-config-file-200"></span> 200 - Returns the content of the uploaded KubeConfig file.
Status: OK

###### <span id="post-rest-api-v1-cluster-config-file-200-schema"></span> Schema
   
  

[ClusterUploadData](#cluster-upload-data)

##### <span id="post-rest-api-v1-cluster-config-file-400"></span> 400 - The uploaded file is too large or the request is invalid.
Status: Bad Request

###### <span id="post-rest-api-v1-cluster-config-file-400-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-config-file-500"></span> 500 - Internal server error.
Status: Internal Server Error

###### <span id="post-rest-api-v1-cluster-config-file-500-schema"></span> Schema
   
  



### <span id="post-rest-api-v1-cluster-config-validate"></span> Validate KubeConfig (*PostRestAPIV1ClusterConfigValidate*)

```
POST /rest-api/v1/cluster/config/validate
```

Validates the provided KubeConfig using cluster manager methods.

#### Consumes
  * application/json
  * text/plain

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [ClusterValidatePayload](#cluster-validate-payload) | `models.ClusterValidatePayload` | | ✓ | | KubeConfig payload to validate |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-rest-api-v1-cluster-config-validate-200) | OK | Verification passed server version |  | [schema](#post-rest-api-v1-cluster-config-validate-200-schema) |
| [400](#post-rest-api-v1-cluster-config-validate-400) | Bad Request | Bad Request |  | [schema](#post-rest-api-v1-cluster-config-validate-400-schema) |
| [401](#post-rest-api-v1-cluster-config-validate-401) | Unauthorized | Unauthorized |  | [schema](#post-rest-api-v1-cluster-config-validate-401-schema) |
| [404](#post-rest-api-v1-cluster-config-validate-404) | Not Found | Not Found |  | [schema](#post-rest-api-v1-cluster-config-validate-404-schema) |
| [429](#post-rest-api-v1-cluster-config-validate-429) | Too Many Requests | Too Many Requests |  | [schema](#post-rest-api-v1-cluster-config-validate-429-schema) |
| [500](#post-rest-api-v1-cluster-config-validate-500) | Internal Server Error | Internal Server Error |  | [schema](#post-rest-api-v1-cluster-config-validate-500-schema) |

#### Responses


##### <span id="post-rest-api-v1-cluster-config-validate-200"></span> 200 - Verification passed server version
Status: OK

###### <span id="post-rest-api-v1-cluster-config-validate-200-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-config-validate-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-rest-api-v1-cluster-config-validate-400-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-config-validate-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="post-rest-api-v1-cluster-config-validate-401-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-config-validate-404"></span> 404 - Not Found
Status: Not Found

###### <span id="post-rest-api-v1-cluster-config-validate-404-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-config-validate-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="post-rest-api-v1-cluster-config-validate-429-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-cluster-config-validate-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-rest-api-v1-cluster-config-validate-500-schema"></span> Schema
   
  



### <span id="post-rest-api-v1-resource-group-rule"></span> Create creates a ResourceGroupRule. (*PostRestAPIV1ResourceGroupRule*)

```
POST /rest-api/v1/resource-group-rule
```

This endpoint creates a new ResourceGroupRule using the payload.

#### Consumes
  * application/json
  * text/plain

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [ResourcegroupruleResourceGroupRulePayload](#resourcegrouprule-resource-group-rule-payload) | `models.ResourcegroupruleResourceGroupRulePayload` | | ✓ | | resourceGroupRule to create (either plain text or JSON format) |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-rest-api-v1-resource-group-rule-200) | OK | Unstructured object |  | [schema](#post-rest-api-v1-resource-group-rule-200-schema) |
| [400](#post-rest-api-v1-resource-group-rule-400) | Bad Request | Bad Request |  | [schema](#post-rest-api-v1-resource-group-rule-400-schema) |
| [401](#post-rest-api-v1-resource-group-rule-401) | Unauthorized | Unauthorized |  | [schema](#post-rest-api-v1-resource-group-rule-401-schema) |
| [404](#post-rest-api-v1-resource-group-rule-404) | Not Found | Not Found |  | [schema](#post-rest-api-v1-resource-group-rule-404-schema) |
| [405](#post-rest-api-v1-resource-group-rule-405) | Method Not Allowed | Method Not Allowed |  | [schema](#post-rest-api-v1-resource-group-rule-405-schema) |
| [429](#post-rest-api-v1-resource-group-rule-429) | Too Many Requests | Too Many Requests |  | [schema](#post-rest-api-v1-resource-group-rule-429-schema) |
| [500](#post-rest-api-v1-resource-group-rule-500) | Internal Server Error | Internal Server Error |  | [schema](#post-rest-api-v1-resource-group-rule-500-schema) |

#### Responses


##### <span id="post-rest-api-v1-resource-group-rule-200"></span> 200 - Unstructured object
Status: OK

###### <span id="post-rest-api-v1-resource-group-rule-200-schema"></span> Schema
   
  

[UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="post-rest-api-v1-resource-group-rule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-rest-api-v1-resource-group-rule-400-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-resource-group-rule-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="post-rest-api-v1-resource-group-rule-401-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-resource-group-rule-404"></span> 404 - Not Found
Status: Not Found

###### <span id="post-rest-api-v1-resource-group-rule-404-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-resource-group-rule-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="post-rest-api-v1-resource-group-rule-405-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-resource-group-rule-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="post-rest-api-v1-resource-group-rule-429-schema"></span> Schema
   
  



##### <span id="post-rest-api-v1-resource-group-rule-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-rest-api-v1-resource-group-rule-500-schema"></span> Schema
   
  



### <span id="put-rest-api-v1-cluster-cluster-name"></span> Update updates the cluster metadata by name. (*PutRestAPIV1ClusterClusterName*)

```
PUT /rest-api/v1/cluster/{clusterName}
```

This endpoint updates the display name and description of an existing cluster resource.

#### Consumes
  * application/json
  * text/plain

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| clusterName | `path` | string | `string` |  | ✓ |  | The name of the cluster |
| request | `body` | [ClusterClusterPayload](#cluster-cluster-payload) | `models.ClusterClusterPayload` | | ✓ | | cluster to update (either plain text or JSON format) |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#put-rest-api-v1-cluster-cluster-name-200) | OK | Unstructured object |  | [schema](#put-rest-api-v1-cluster-cluster-name-200-schema) |
| [400](#put-rest-api-v1-cluster-cluster-name-400) | Bad Request | Bad Request |  | [schema](#put-rest-api-v1-cluster-cluster-name-400-schema) |
| [401](#put-rest-api-v1-cluster-cluster-name-401) | Unauthorized | Unauthorized |  | [schema](#put-rest-api-v1-cluster-cluster-name-401-schema) |
| [404](#put-rest-api-v1-cluster-cluster-name-404) | Not Found | Not Found |  | [schema](#put-rest-api-v1-cluster-cluster-name-404-schema) |
| [405](#put-rest-api-v1-cluster-cluster-name-405) | Method Not Allowed | Method Not Allowed |  | [schema](#put-rest-api-v1-cluster-cluster-name-405-schema) |
| [429](#put-rest-api-v1-cluster-cluster-name-429) | Too Many Requests | Too Many Requests |  | [schema](#put-rest-api-v1-cluster-cluster-name-429-schema) |
| [500](#put-rest-api-v1-cluster-cluster-name-500) | Internal Server Error | Internal Server Error |  | [schema](#put-rest-api-v1-cluster-cluster-name-500-schema) |

#### Responses


##### <span id="put-rest-api-v1-cluster-cluster-name-200"></span> 200 - Unstructured object
Status: OK

###### <span id="put-rest-api-v1-cluster-cluster-name-200-schema"></span> Schema
   
  

[UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="put-rest-api-v1-cluster-cluster-name-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="put-rest-api-v1-cluster-cluster-name-400-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-cluster-cluster-name-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="put-rest-api-v1-cluster-cluster-name-401-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-cluster-cluster-name-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-rest-api-v1-cluster-cluster-name-404-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-cluster-cluster-name-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="put-rest-api-v1-cluster-cluster-name-405-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-cluster-cluster-name-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="put-rest-api-v1-cluster-cluster-name-429-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-cluster-cluster-name-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="put-rest-api-v1-cluster-cluster-name-500-schema"></span> Schema
   
  



### <span id="put-rest-api-v1-resource-group-rule"></span> Update updates the ResourceGroupRule metadata by name. (*PutRestAPIV1ResourceGroupRule*)

```
PUT /rest-api/v1/resource-group-rule
```

This endpoint updates the display name and description of an existing ResourceGroupRule.

#### Consumes
  * application/json
  * text/plain

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [ResourcegroupruleResourceGroupRulePayload](#resourcegrouprule-resource-group-rule-payload) | `models.ResourcegroupruleResourceGroupRulePayload` | | ✓ | | resourceGroupRule to update (either plain text or JSON format) |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#put-rest-api-v1-resource-group-rule-200) | OK | Unstructured object |  | [schema](#put-rest-api-v1-resource-group-rule-200-schema) |
| [400](#put-rest-api-v1-resource-group-rule-400) | Bad Request | Bad Request |  | [schema](#put-rest-api-v1-resource-group-rule-400-schema) |
| [401](#put-rest-api-v1-resource-group-rule-401) | Unauthorized | Unauthorized |  | [schema](#put-rest-api-v1-resource-group-rule-401-schema) |
| [404](#put-rest-api-v1-resource-group-rule-404) | Not Found | Not Found |  | [schema](#put-rest-api-v1-resource-group-rule-404-schema) |
| [405](#put-rest-api-v1-resource-group-rule-405) | Method Not Allowed | Method Not Allowed |  | [schema](#put-rest-api-v1-resource-group-rule-405-schema) |
| [429](#put-rest-api-v1-resource-group-rule-429) | Too Many Requests | Too Many Requests |  | [schema](#put-rest-api-v1-resource-group-rule-429-schema) |
| [500](#put-rest-api-v1-resource-group-rule-500) | Internal Server Error | Internal Server Error |  | [schema](#put-rest-api-v1-resource-group-rule-500-schema) |

#### Responses


##### <span id="put-rest-api-v1-resource-group-rule-200"></span> 200 - Unstructured object
Status: OK

###### <span id="put-rest-api-v1-resource-group-rule-200-schema"></span> Schema
   
  

[UnstructuredUnstructured](#unstructured-unstructured)

##### <span id="put-rest-api-v1-resource-group-rule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="put-rest-api-v1-resource-group-rule-400-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-resource-group-rule-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="put-rest-api-v1-resource-group-rule-401-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-resource-group-rule-404"></span> 404 - Not Found
Status: Not Found

###### <span id="put-rest-api-v1-resource-group-rule-404-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-resource-group-rule-405"></span> 405 - Method Not Allowed
Status: Method Not Allowed

###### <span id="put-rest-api-v1-resource-group-rule-405-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-resource-group-rule-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="put-rest-api-v1-resource-group-rule-429-schema"></span> Schema
   
  



##### <span id="put-rest-api-v1-resource-group-rule-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="put-rest-api-v1-resource-group-rule-500-schema"></span> Schema
   
  



## Models

### <span id="aggregator-diagnose-request"></span> aggregator.DiagnoseRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| language | string| `string` |  | | Language code for AI response |  |
| logs | []string| `[]string` |  | |  |  |



### <span id="aggregator-event-diagnose-request"></span> aggregator.EventDiagnoseRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| events | [][AiEvent](#ai-event)| `[]*AiEvent` |  | |  |  |
| language | string| `string` |  | |  |  |



### <span id="aggregator-log-entry"></span> aggregator.LogEntry


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| content | string| `string` |  | |  |  |
| error | string| `string` |  | |  |  |
| timestamp | string| `string` |  | |  |  |



### <span id="ai-diagnosis-event"></span> ai.DiagnosisEvent


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| content | string| `string` |  | | Event content |  |
| type | string| `string` |  | | Event type: start/chunk/error/complete |  |



### <span id="ai-event"></span> ai.Event


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| count | integer| `int64` |  | |  |  |
| firstTimestamp | string| `string` |  | |  |  |
| lastTimestamp | string| `string` |  | |  |  |
| message | string| `string` |  | |  |  |
| reason | string| `string` |  | |  |  |
| type | string| `string` |  | |  |  |



### <span id="ai-interpret-event"></span> ai.InterpretEvent


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| content | string| `string` |  | | Event content or error message |  |
| type | string| `string` |  | | Event type: start, chunk, error, complete |  |



### <span id="cluster-cluster-payload"></span> cluster.ClusterPayload


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| description | string| `string` |  | | ClusterDescription is the description of cluster to be created |  |
| displayName | string| `string` |  | | ClusterDisplayName is the display name of cluster to be created |  |
| kubeConfig | string| `string` |  | | ClusterKubeConfig is the kubeconfig of cluster to be created |  |



### <span id="cluster-upload-data"></span> cluster.UploadData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| content | string| `string` |  | |  |  |
| fileName | string| `string` |  | |  |  |
| fileSize | integer| `int64` |  | |  |  |
| sanitizedClusterContent | string| `string` |  | |  |  |



### <span id="cluster-validate-payload"></span> cluster.ValidatePayload


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| kubeConfig | string| `string` |  | |  |  |



### <span id="detail-interpret-request"></span> detail.InterpretRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| language | string| `string` |  | |  |  |
| yaml | string| `string` |  | |  |  |



### <span id="entity-resource-group"></span> entity.ResourceGroup


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| annotations | map of string| `map[string]string` |  | |  |  |
| apiVersion | string| `string` |  | |  |  |
| cluster | string| `string` |  | |  |  |
| kind | string| `string` |  | |  |  |
| labels | map of string| `map[string]string` |  | |  |  |
| name | string| `string` |  | |  |  |
| namespace | string| `string` |  | |  |  |
| status | string| `string` |  | |  |  |



### <span id="insight-resource-summary"></span> insight.ResourceSummary


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| creationTimestamp | string| `string` |  | |  |  |
| resource | [EntityResourceGroup](#entity-resource-group)| `EntityResourceGroup` |  | |  |  |
| resourceVersion | string| `string` |  | |  |  |
| uid | string| `string` |  | |  |  |



### <span id="insight-resource-topology"></span> insight.ResourceTopology


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| children | []string| `[]string` |  | |  |  |
| parents | []string| `[]string` |  | |  |  |
| resourceGroup | [EntityResourceGroup](#entity-resource-group)| `EntityResourceGroup` |  | |  |  |



### <span id="insight-score-data"></span> insight.ScoreData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| issuesTotal | integer| `int64` |  | | IssuesTotal is the total count of all issues found during the audit.
This count can be used to understand the overall number of problems
that need to be addressed. |  |
| resourceTotal | integer| `int64` |  | | ResourceTotal is the count of unique resources audited during the scan. |  |
| score | number| `float64` |  | | Score represents the calculated score of the audited manifest based on
the number and severity of issues. It provides a quantitative measure
of the security posture of the resources in the manifest. |  |
| severityStatistic | map of integer| `map[string]int64` |  | | SeverityStatistic is a mapping of severity levels to their respective
number of occurrences. It allows for a quick overview of the distribution
of issues across different severity categories. |  |



### <span id="insight-statistics"></span> insight.Statistics


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| clusterCount | integer| `int64` |  | |  |  |
| resourceCount | integer| `int64` |  | |  |  |
| resourceGroupRuleCount | integer| `int64` |  | |  |  |



### <span id="resourcegrouprule-resource-group-rule-payload"></span> resourcegrouprule.ResourceGroupRulePayload


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| description | string| `string` |  | |  |  |
| fields | []string| `[]string` |  | |  |  |
| name | string| `string` |  | |  |  |



### <span id="scanner-audit-data"></span> scanner.AuditData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| bySeverity | map of integer| `map[string]int64` |  | |  |  |
| issueGroups | [][ScannerIssueGroup](#scanner-issue-group)| `[]*ScannerIssueGroup` |  | |  |  |
| issueTotal | integer| `int64` |  | |  |  |
| resourceTotal | integer| `int64` |  | |  |  |



### <span id="scanner-issue"></span> scanner.Issue


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| message | string| `string` |  | | Message provides a detailed human-readable description of the issue. |  |
| scanner | string| `string` |  | | Scanner is the name of the scanner that discovered the issue. |  |
| severity | integer| `int64` |  | | Severity indicates how critical the issue is, using the IssueSeverityLevel constants. |  |
| title | string| `string` |  | | Title is a brief summary of the issue. |  |



### <span id="scanner-issue-group"></span> scanner.IssueGroup


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| issue | [ScannerIssue](#scanner-issue)| `ScannerIssue` |  | |  |  |
| resourceGroups | [][EntityResourceGroup](#entity-resource-group)| `[]*EntityResourceGroup` |  | |  |  |



### <span id="unstructured-unstructured"></span> unstructured.Unstructured


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| object | [interface{}](#interface)| `interface{}` |  | | Object is a JSON compatible map with string, float, int, bool, []interface{}, or
map[string]interface{}
children. |  |


