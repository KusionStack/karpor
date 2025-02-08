// Issue Response Type

export interface Issue {
  scanner: string
  severity: string
  title: string
  message: string
}

export interface ResourceGroup {
  cluster: string
  apiVersion: string
  kind: string
  namespace: string
  name: string
}

export interface IssueGroup {
  issue: Issue
  resourceGroups: ResourceGroup[]
}

export type SeverityLevel = 'High' | 'Medium' | 'Low'

export interface SeverityCount {
  High?: number
  Medium?: number
  Low?: number
}

export interface IssueResponse {
  issueTotal: number
  resourceTotal: number
  bySeverity: SeverityCount
  issueGroups: IssueGroup[]
}

// Score Response Type
export interface ResourceScore {
  score: number
  resourceTotal: number
  issuesTotal: number
  severityStatistic: SeverityCount
}

// Cluster Summary Response Type

interface Metrics {
  points: null | any
}

export interface ClusterSummaryInfo {
  nodeCount: number
  serverVersion: string
  readyNodes: number
  notReadyNodes: number
  memoryCapacity: number
  memoryUsage: number
  cpuCapacity: number
  cpuUsage: number
  podsCapacity: number
  podsUsage: number
  latency: number
  metricsEnabled: boolean
  cpuMetrics: Metrics
  memoryMetrics: Metrics
}

// Kind Response Summary Type

export interface KindSummaryInfo {
  cluster: string
  group: string
  version: string
  kind: string
  count: number
}

// Namespace Response Summary Type

interface CountByGVK {
  [key: string]: number
}

export interface NamespaceSummaryInfo {
  cluster: string
  namespace: string
  countByGVK: CountByGVK
}

// Resource Summary Response Type

export interface ResourceSummaryInfo {
  cluster: string
  apiVersion: string
  kind: string
  namespace: string
  name: string
  status: string
}

// Topology Response Data Type

// Cluster / Namespace
export interface TopologyClusterResourceGroup {
  cluster: string
  apiVersion: string
  kind: string
}

export interface TopologyRelationship {
  [key: string]: 'parent' | 'child'
}

export interface TopologyClusterResourceType {
  resourceGroup: TopologyClusterResourceGroup
  count: number
  relationship: TopologyRelationship
}

export interface TopologyClusterResourcesData {
  [key: string]: TopologyClusterResourceType
}

export interface ClusterResourceData {
  [clusterName: string]: TopologyClusterResourcesData
}

// Resource

export interface TopologyResourceGroup {
  cluster: string
  apiVersion: string
  kind: string
  name: string
  namespace?: string
}

export interface TopologyResourceItem {
  resourceGroup: TopologyResourceGroup
  parents: string[]
  children: string[]
}

export interface TopologyClusterResources {
  [resourceKey: string]: TopologyResourceItem
}

export interface TopologyClusterResourceData {
  [resourceClusterName: string]: TopologyClusterResources
}

// Source Table Response Data Type

interface ManagedField {
  apiVersion: string
  fieldsType: string
  fieldsV1: {
    [key: string]: any
  }
  manager: string
  operation: string
  time: string
}

interface Metadata {
  creationTimestamp: string
  labels: {
    [key: string]: string
  }
  managedFields: ManagedField[]
  name: string
  namespace: string
  resourceVersion: string
  uid: string
}

interface ResourceObject {
  apiVersion: string
  kind: string
  metadata: Metadata
}

interface Item {
  cluster: string
  object: ResourceObject
  syncAt: string
  deleted: boolean
}

interface TableData {
  items: Item[]
  sqlQuery: string
  total: number
  currentPage: number
  pageSize: number
}

export interface TableDataResponse {
  data: TableData
}
