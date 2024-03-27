


# Karbour
  

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
  * text/plain

## All endpoints

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
| GET | /rest-api/v1/insight/audit | [get rest API v1 insight audit](#get-rest-api-v1-insight-audit) | Audit based on locator. |
| GET | /rest-api/v1/insight/detail | [get rest API v1 insight detail](#get-rest-api-v1-insight-detail) | GetDetail returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind. |
| GET | /rest-api/v1/insight/events | [get rest API v1 insight events](#get-rest-api-v1-insight-events) | GetEvents returns events for a Kubernetes resource by name, namespace, cluster, apiVersion and kind. |
| GET | /rest-api/v1/insight/score | [get rest API v1 insight score](#get-rest-api-v1-insight-score) | ScoreHandler calculates a score for the audited manifest. |
| GET | /rest-api/v1/insight/summary | [get rest API v1 insight summary](#get-rest-api-v1-insight-summary) | Get returns a Kubernetes resource summary by name, namespace, cluster, apiVersion and kind. |
| GET | /rest-api/v1/insight/topology | [get rest API v1 insight topology](#get-rest-api-v1-insight-topology) | GetTopology returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind. |
  


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
   
  



### <span id="get-rest-api-v1-insight-audit"></span> Audit based on locator. (*GetRestAPIV1InsightAudit*)

```
GET /rest-api/v1/insight/audit
```

This endpoint audits based on the specified locator.

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
| pattern | `query` | string | `string` |  | ✓ |  | The search pattern. Can be either sql or dsl. Required |
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
   
  



## Models

### <span id="cluster-cluster-payload"></span> cluster.ClusterPayload


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| description | string| `string` |  | | ClusterDescription is the description of cluster to be created |  |
| displayName | string| `string` |  | | ClusterDisplayName is the display name of cluster to be created |  |
| kubeconfig | string| `string` |  | | ClusterKubeConfig is the kubeconfig of cluster to be created |  |



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



### <span id="core-locator"></span> core.Locator


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| apiVersion | string| `string` |  | |  |  |
| cluster | string| `string` |  | |  |  |
| kind | string| `string` |  | |  |  |
| name | string| `string` |  | |  |  |
| namespace | string| `string` |  | |  |  |



### <span id="insight-resource-summary"></span> insight.ResourceSummary


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| creationTimestamp | string| `string` |  | |  |  |
| resource | [CoreLocator](#core-locator)| `CoreLocator` |  | |  |  |
| resourceVersion | string| `string` |  | |  |  |
| uid | string| `string` |  | |  |  |



### <span id="insight-resource-topology"></span> insight.ResourceTopology


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| children | []string| `[]string` |  | |  |  |
| locator | [CoreLocator](#core-locator)| `CoreLocator` |  | |  |  |
| parents | []string| `[]string` |  | |  |  |



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
| locators | [][CoreLocator](#core-locator)| `[]*CoreLocator` |  | |  |  |



### <span id="unstructured-unstructured"></span> unstructured.Unstructured


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| object | [interface{}](#interface)| `interface{}` |  | | Object is a JSON compatible map with string, float, int, bool, []interface{}, or
map[string]interface{}
children. |  |


