import cRole from '@/assets/labeled/c-role-128.png'
import cm from '@/assets/labeled/cm-128.png'
import crb from '@/assets/labeled/crb-128.png'
import cronjob from '@/assets/labeled/cronjob-128.png'
import deploy from '@/assets/labeled/deploy-128.png'
import ds from '@/assets/labeled/ds-128.png'
import ep from '@/assets/labeled/ep-128.png'
import group from '@/assets/labeled/group-128.png'
import hpa from '@/assets/labeled/hpa-128.png'
import ing from '@/assets/labeled/ing-128.png'
import job from '@/assets/labeled/job-128.png'
import limits from '@/assets/labeled/limits-128.png'
import netpol from '@/assets/labeled/netpol-128.png'
import ns from '@/assets/labeled/ns-128.png'
import pod from '@/assets/labeled/pod-128.png'
import psp from '@/assets/labeled/psp-128.png'
import pv from '@/assets/labeled/pv-128.png'
import pvc from '@/assets/labeled/pvc-128.png'
import quota from '@/assets/labeled/quota-128.png'
import rb from '@/assets/labeled/rb-128.png'
import role from '@/assets/labeled/role-128.png'
import rs from '@/assets/labeled/rs-128.png'
import sa from '@/assets/labeled/sa-128.png'
import sc from '@/assets/labeled/sc-128.png'
import secret from '@/assets/labeled/secret-128.png'
import sts from '@/assets/labeled/sts-128.png'
import svc from '@/assets/labeled/svc-128.png'
import user from '@/assets/labeled/user-128.png'
import volume from '@/assets/labeled/vol-128.png'
import nodeIcon from '@/assets/labeled/node-128.png'
import crd from '@/assets/labeled/crd-128.png'
import kubernetes from '@/assets/kubernetes.png'

export const ICON_MAP = {
  ClusterRole: cRole,
  ConfigMap: cm,
  ClusterRoleBinding: crb,
  CronJob: cronjob,
  Deployment: deploy,
  CafeDeployment: crd,
  DaemonSet: ds,
  Endpoint: ep,
  Group: group,
  HorizontalPodAutoscaler: hpa,
  Ingress: ing,
  Job: job,
  Limits: limits,
  NetworkPolicy: netpol,
  Namespace: ns,
  Pod: pod,
  PodSecurityPolicy: psp,
  PersistentVolume: pv,
  PersistentVolumeClaim: pvc,
  ResourceQuota: quota,
  RoleBinding: rb,
  Role: role,
  ReplicaSet: rs,
  ServiceAccount: sa,
  StorageClass: sc,
  Secret: secret,
  StatefulSet: sts,
  Service: svc,
  User: user,
  Volume: volume,
  Node: nodeIcon,
  InPlaceSet: crd,
  PodDisruptionBudget: crd,
  CRD: crd,
  Kubernetes: kubernetes,
}
