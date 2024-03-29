export const SEVERITY_MAP = {
  Critical: {
    color: 'purple',
    text: 'Critical',
  },
  High: {
    color: 'error',
    text: 'High',
  },
  Medium: {
    color: 'warning',
    text: 'Medium',
  },
  Low: {
    color: 'success',
    text: 'Low',
  },
  Safe: {
    color: 'success',
    text: 'Safe',
  },
}

export const searchSqlPrefix = 'select * from resources'

export const whereKeywords = [
  'cluster',
  'kind',
  'namespace',
  'name',
  'content',
  '`labels.${key}`',
  'labels_key',
  'labels_value',
  '`annotations.${key}`',
  'annotations_key',
  'annotations_value',
]
export const operatorKeywords = ['=', 'like', 'contains']
export const kindCompletions = [
  `'Pod'`,
  `'Service'`,
  `'ReplicaSet'`,
  `'Deployment'`,
  `'StatefulSet'`,
  `'DaemonSet'`,
  `'Job'`,
  `'CronJob'`,
  `'Namespace'`,
  `'ServiceAccount'`,
  `'ConfigMap'`,
  `'Secret'`,
  `'PersistentVolume'`,
  `'PersistentVolumeClaim'`,
  `'Ingress'`,
  `'StorageClass'`,
  `'NetworkPolicy'`,
  `'ResourceQuota'`,
  `'LimitRange'`,
  `'Role'`,
  `'ClusterRole'`,
  `'RoleBinding'`,
  `'ClusterRoleBinding'`,
  `'CustomResourceDefinition'`,
  `'PodDisruptionBudget'`,
  `'HorizontalPodAutoscaler'`,
  `'PodSecurityPolicy'`,
  `'EndpointSlices'`,
]
export const defaultKeywords = [
  'select',
  'from',
  'where',
  'values',
  'as',
  'join',
  'on',
  'group by',
  'order by',
  'having',
  'or',
  'and',
]
