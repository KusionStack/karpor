export const yamlStr = `apiVersion: v1
kind: Pod
metadata:
  annotations:
    alibabacloud.com/actual-pod-cgroup-path: /sigma/pod883a112f-386b-4b34-a727-9b786426cab3
    cafe.sofastack.io/rulesets: '{"ocmpaas-cz30a-test":{"passPreCheck":false,"passPostCheck":false,"ruleName":"","state":"","timestamp":"0001-01-01T00:00:00Z"}}'
    custom.k8s.alipay.com/original-resource: '{"containers":[{"name":"apiserver","Resources":{"limits":{"cpu":"6","ephemeral-storage":"100Gi","memory":"8Gi"},"requests":{"cpu":"6","ephemeral-storage":"100Gi","memory":"8Gi"}}}]}'
    meta.k8s.alipay.com/last-spec-hash: a12c8cd2e2b218307e2091c1846a37d8
    meta.k8s.alipay.com/pod-zappinfo: '{"spec":{"appName":"ocmpaas","zone":"CZ30A","serverType":"DOCKER","fqdn":"ocmpaas-cz30a-011166232074.eu95.alipay.net","expectStatus":""},"status":{"registered":true,"message":"","status":"online"}}'
    meta.k8s.alipay.com/trace-context: '[{"trace_id":"a812a15cdde89ace0000000000000000","parent_id":"","root_span_id":"28cc32b01ddaf697","delivery_type":"PodCreate","status":"closed","services":[{"component":"cloud-scheduler","span_id":"3ca83c434e40cc23"},{"component":"default-scheduler","span_id":"340f2185df97e4d6"},{"component":"cni-service","span_id":"7ebf535062961ff1"},{"component":"kubelet","span_id":"cf7869c88a5ed2fe"},{"component":"zappinfo-controller","span_id":"c340b3fb48e0ab4c"},{"component":"naming-controller","span_id":"84f24b8adf45d2dc"}],"start_at":"2023-05-09T11:52:48+08:00","finish_at":"2023-05-09T11:53:09+08:00","extra_info":null}]'
    orca.identity.alipay.com/serviceaccount: "true"
    paascore.alipay.com/upgrade-diff: 37e3d1263c41823d25978acbec3a75ec
    pod.beta1.sigma.ali/alloc-spec: '{"containers":[{"name":"apiserver","resource":{"cpu":{},"gpu":{"shareMode":"exclusive"}},"hostConfig":{"cgroupParent":"/sigma","diskQuotaMode":"","memorySwap":8589934592,"pidsLimit":32767,"cpuBvtWarpNs":2,"memoryWmarkRatio":95,"cpuShares":6144,"oomScoreAdj":-1}}]}'
    pod.beta1.sigma.ali/hostname-template: ocmpaas-cz30a-{{.IpAddress}}
    pod.beta1.sigma.ali/net-priority: "5"
    pod.beta1.sigma.ali/network-status: '{"ipam":"ais-ipam","vlan":"701","networkPrefixLen":24,"gateway":"11.166.232.247","netType":"vlan","sandboxId":"","ip":"11.166.232.74","securityDomain":"ALIPAY_TEST"}'
    pod.beta1.sigma.ali/pod-spec-hash: cluster-extension-apiserver-jqntf-549d7f46d8
    pod.beta1.sigma.ali/scheduler-update-time: "2023-05-09T11:52:49.383784666+08:00"
    pod.beta1.sigma.ali/trace-id: 96055ad3-c1cb-4b76-a1db-746697bb452b
    pod.beta1.sigma.ali/trace-naming: '{"id":"96055ad3-c1cb-4b76-a1db-746697bb452b","service":"naming","creationTimestamp":"2023-05-09T11:53:06.619814492+08:00","completionTimestamp":"2023-05-09T11:53:06.891066682+08:00","executionCount":1,"span":{"operation":"ReconcilerHelper.Reconcile","success":true,"startTimestamp":"2023-05-09T11:53:06.619814492+08:00","endTimestamp":"2023-05-09T11:53:06.891066682+08:00","logs":[{"time":"2023-05-09T11:53:06.619817226+08:00","fields":[{"key":"action","value":"update"}]}],"children":[{"operation":"reconcile","success":true,"startTimestamp":"2023-05-09T11:53:06.619821911+08:00","endTimestamp":"2023-05-09T11:53:06.891066396+08:00"}]}}'
    pod.beta1.sigma.ali/trace-podfqdn: '{"TraceID":"96055ad3-c1cb-4b76-a1db-746697bb452b","Service":"podfqdn","Operation":"AddResourceRecord","Error":false,"Message":"","StartTimestamp":"2023-05-09T11:53:06.620143737+08:00","FinishTimestamp":"2023-05-09T11:53:06.737906027+08:00","Logs":{"error":""}}'
    pod.beta1.sigma.ali/trace-zappinfo: '{"id":"96055ad3-c1cb-4b76-a1db-746697bb452b","service":"zappinfo","creationTimestamp":"2023-05-09T11:53:09.585985326+08:00","completionTimestamp":"2023-05-09T11:53:09.638850489+08:00","executionCount":1,"span":{"operation":"ReconcilerHelper.Reconcile","success":true,"startTimestamp":"2023-05-09T11:53:09.585985326+08:00","endTimestamp":"2023-05-09T11:53:09.638850489+08:00","logs":[{"time":"2023-05-09T11:53:09.585988543+08:00","fields":[{"key":"action","value":"update"}]}],"children":[{"operation":"reconcile","success":true,"startTimestamp":"2023-05-09T11:53:09.585991693+08:00","endTimestamp":"2023-05-09T11:53:09.638850084+08:00","children":[{"operation":"zappinfo.update","success":true,"startTimestamp":"2023-05-09T11:53:09.586046382+08:00","endTimestamp":"2023-05-09T11:53:09.638837609+08:00","children":[{"operation":"getPodZappinfoMetaSpec","success":true,"startTimestamp":"2023-05-09T11:53:09.586047934+08:00","endTimestamp":"2023-05-09T11:53:09.586212157+08:00"},{"operation":"updateServerStatus","success":true,"startTimestamp":"2023-05-09T11:53:09.586216493+08:00","endTimestamp":"2023-05-09T11:53:09.638837148+08:00"}]}]}]}}'
    pod.beta1.sigma.ali/update-status: '{"statuses":{"apiserver":{"creationTimestamp":"2023-05-09T11:52:51.333512804+08:00","finishTimestamp":"2023-05-09T11:53:05.240263629+08:00","retryCount":0,"currentState":"running","lastState":"unknown","action":"start","success":true,"message":"create
      start and post start success","specHash":"cluster-extension-apiserver-jqntf-549d7f46d8"}}}'
    pod.k8s.alipay.com/auto-eviction: "true"
    pod.k8s.alipay.com/fqdn-registered-timestamp: 2023-05-09 11:53:06.737936235 +0800
      CST m=+1304212.764360580
    sigma.ali/container-diskQuotaID: '{"apiserver":"16913742"}'
    trace.cafe.sofastack.io/distribution-info: '{"Stage":"2","Id":"RELEASE_202305060509499566181"}'
    ulogfs.k8s.alipay.com/biz-disk-quota-repaired: "true"
    ulogfs.k8s.alipay.com/enable-zclean: "true"
    ulogfs.k8s.alipay.com/inject: enabled
  creationTimestamp: "2023-05-09T03:52:48Z"
  finalizers:
  - finalizer.k8s.alipay.com/zappinfo
  - protection-delete.pod.sigma.ali/naming-registered
  - pod.beta1.sigma.ali/cni-allocated
  - finalizers.k8s.alipay.com/pod-fqdn
  generateName: cluster-extension-apiserver-jqntf-
  labels:
    alibabacloud.com/quota-name: ocmpaas-test-sigmaguaranteed-daily
    cafe.sofastack.io/app-instance-group: ""
    cafe.sofastack.io/app-instance-group-name: ""
    cafe.sofastack.io/cell: CZ30A
    cafe.sofastack.io/control: "true"
    cafe.sofastack.io/creator: huanyu
    cafe.sofastack.io/deploy-type: workload
    cafe.sofastack.io/global-tenant: MAIN_SITE
    cafe.sofastack.io/pod-ip: 11.166.232.74
    cafe.sofastack.io/pod-number: "3"
    cafe.sofastack.io/pre-check: "false"
    cafe.sofastack.io/service-available: "1683604389587189074"
    cafe.sofastack.io/version: cluster-extension-apiserver-jqntf-549d7f46d8
    cluster.x-k8s.io/cluster-name: eu95
    component: cluster-extension-apiserver
    controller-revision-hash: cluster-extension-apiserver-jqntf-549d7f46d8
    meta.k8s.alipay.com/app-env: TEST
    meta.k8s.alipay.com/biz-group: ocmpaas
    meta.k8s.alipay.com/biz-group-id: cluster-extension-apiserver-jqntf-30bd4911-a01c-4bd8-be1d-48dc67520374
    meta.k8s.alipay.com/biz-name: cloudprovision
    meta.k8s.alipay.com/delivery-workload: paascore-cafeext
    meta.k8s.alipay.com/fqdn: ocmpaas-cz30a-011166232074.eu95.alipay.net
    meta.k8s.alipay.com/hostname: ocmpaas-cz30a-011166232074
    meta.k8s.alipay.com/migration-level: L2
    meta.k8s.alipay.com/min-replicas: "1"
    meta.k8s.alipay.com/original-pod-namespace: ocmpaas
    meta.k8s.alipay.com/priority: production
    meta.k8s.alipay.com/qoc-class: ProdGeneral
    meta.k8s.alipay.com/qos-class: Prod
    meta.k8s.alipay.com/replicas: "1"
    meta.k8s.alipay.com/schedule-time-limit: 10m0s
    meta.k8s.alipay.com/situation: normal
    meta.k8s.alipay.com/slo-resource: 8C16G
    meta.k8s.alipay.com/slo-scale: "10"
    meta.k8s.alipay.com/zone: CZ30A
    paascore.alipay.com/adopted: "1683604368844361705"
    sigma.ali/app-name: ocmpaas
    sigma.ali/deploy-unit: ocmpaas-test
    sigma.ali/force-update-quota-name: 20230516-135758
    sigma.ali/instance-group: ocmpaassqa
    sigma.ali/ip: 11.166.232.74
    sigma.ali/qos: SigmaBurstable
    sigma.ali/site: eu95
    sigma.ali/sn: 1dff5a86-80b7-4803-b9ad-dcd7d5a98adb
    strategy.cafe.sofastack.io/batch-index: RELEASE_202305060509499566181-2
  name: cluster-extension-apiserver-jqntf-6x676
  namespace: ocmpaas
  ownerReferences:
  - apiVersion: apps.cafe.cloud.alipay.com/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: InPlaceSet
    name: cluster-extension-apiserver-jqntf
    uid: 30bd4911-a01c-4bd8-be1d-48dc67520374
  resourceVersion: "32675456614"
  selfLink: /api/v1/namespaces/ocmpaas/pods/cluster-extension-apiserver-jqntf-6x676
  uid: 883a112f-386b-4b34-a727-9b786426cab3`
