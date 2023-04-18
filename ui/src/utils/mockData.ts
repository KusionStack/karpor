export const result = [
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-apiserver-5b7nn-q4qsh",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod78d947e2-fc57-4782-989b-4c20068db1ff",
          "cafe.sofastack.io/available-conditions-ex": "{\"expectedFinalizers\":{\"Service/ocmpaas/ocmpaastestcz30axvip0\":\"prot.cafe.sofastack.io/lbb_ocmpaastestcz30axvip0\"}}",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "0e0170259b1c998495407fd256356623",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100088231095.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"b56192b0cccc068a0000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"3c7c99c978bb551d\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"f5c703c8a0161f19\"},{\"component\":\"default-scheduler\",\"span_id\":\"e42c4fe0a4fe90d8\"},{\"component\":\"cni-service\",\"span_id\":\"5126606e3b7275e4\"},{\"component\":\"kubelet\",\"span_id\":\"5bdcbdf260f69a9e\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"d3e51ca68497f29f\"},{\"component\":\"naming-controller\",\"span_id\":\"a13df309fd1c6222\"}],\"start_at\":\"2023-03-06T20:46:16+08:00\",\"finish_at\":\"2023-03-06T20:50:46+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "paascore.alipay.com/upgrade-diff": "77771c94bd5e87ddeb6289afad6ee70e",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":5120,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.88.231.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.88.231.95\",\"securityDomain\":\"ALIPAY_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-apiserver-5b7nn-74777b5468",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-03-06T20:46:17.001856063+08:00",
          "pod.beta1.sigma.ali/trace-id": "ad180600-7b23-49d2-ba58-6d6ce188e444",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"ad180600-7b23-49d2-ba58-6d6ce188e444\",\"service\":\"naming\",\"creationTimestamp\":\"2023-03-06T20:50:44.775020173+08:00\",\"completionTimestamp\":\"2023-03-06T20:50:44.84883619+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:50:44.775020173+08:00\",\"endTimestamp\":\"2023-03-06T20:50:44.84883619+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:50:44.775023591+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:50:44.77502708+08:00\",\"endTimestamp\":\"2023-03-06T20:50:44.848835907+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"ad180600-7b23-49d2-ba58-6d6ce188e444\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-03-06T20:50:44.70626229+08:00\",\"FinishTimestamp\":\"2023-03-06T20:50:44.771369905+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"ad180600-7b23-49d2-ba58-6d6ce188e444\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-06T20:50:47.538056703+08:00\",\"completionTimestamp\":\"2023-03-06T20:50:47.53826277+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:50:47.538056703+08:00\",\"endTimestamp\":\"2023-03-06T20:50:47.53826277+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:50:47.538059981+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:50:47.538063404+08:00\",\"endTimestamp\":\"2023-03-06T20:50:47.538262561+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:50:47.538108841+08:00\",\"endTimestamp\":\"2023-03-06T20:50:47.538248901+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:50:47.538115341+08:00\",\"endTimestamp\":\"2023-03-06T20:50:47.538248579+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-06T20:46:19.251221454+08:00\",\"finishTimestamp\":\"2023-03-06T20:50:44.078879467+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"unknown\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-apiserver-5b7nn-74777b5468\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-03-06 20:50:44.771418204 +0800 CST m=+2939979.592685391",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16785496\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-03-06T12:46:16Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn",
          "xvip.cafe.sofastack.io/xvip_rs_prot_ocmpaastestcz30axvip0",
          "prot.cafe.sofastack.io/lbb_ocmpaastestcz30axvip0"
        ],
        "generateName": "hub-apiserver-5b7nn-",
        "labels": {
          "ali.EnableDefaultRoute": "true",
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "100.88.231.95",
          "cafe.sofastack.io/pod-number": "3",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1678107142742527132",
          "cafe.sofastack.io/version": "hub-apiserver-5b7nn-74777b5468",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-apiserver",
          "controller-revision-hash": "hub-apiserver-5b7nn-74777b5468",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-apiserver-5b7nn-3100a2d0-41b4-472c-9bb8-7dcedbe9af6d",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100088231095.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100088231095",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "10m0s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "8C16G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1678106776666300515",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-183539",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.88.231.95",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "4db1119a-aa72-4956-8a89-6f179ee9eba3",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "hub-apiserver-5b7nn-q4qsh",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-apiserver-5b7nn",
            "uid": "3100a2d0-41b4-472c-9bb8-7dcedbe9af6d"
          }
        ],
        "resourceVersion": "32427670246",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-apiserver-5b7nn-q4qsh",
        "uid": "78d947e2-fc57-4782-989b-4c20068db1ff"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-apiserver",
              "--allow-privileged=true",
              "--authorization-mode=Node,RBAC",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--endpoint-reconciler-type=none",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--enable-admission-plugins=NodeRestriction",
              "--enable-bootstrap-token-auth=true",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--insecure-port=0",
              "--secure-port=8443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--service-account-key-file=/etc/kubernetes/pki/apiserver/sa.pub",
              "--service-account-signing-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--service-account-issuer=api",
              "--api-audiences=api",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--proxy-client-cert-file=/etc/kubernetes/pki/apiserver/proxy.crt",
              "--proxy-client-key-file=/etc/kubernetes/pki/apiserver/proxy.key",
              "--log-file=/home/admin/logs/apiserver.log",
              "--log-file-max-size=100",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "SN",
                "value": "4db1119a-aa72-4956-8a89-6f179ee9eba3"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "5"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            },
            {
              "name": "single-request-reopen"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "215213809-a",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.public.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                },
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "53687091200",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "77389d7d-431b-4ef1-bff9-40bde3922601"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:46:16Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:46:17Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:50:44Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-03-06T12:50:46Z",
            "lastTransitionTime": "2023-03-06T12:50:46Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:46:17Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:50:48Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:50:44Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:46:17Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:46:17Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "docker://49d3d0e06d769d4c0f8ba7bd0a2fdc01cc1d17fa26bd254e14d97abf78186f16",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "docker-pullable://reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 0,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-06T12:50:43Z"
              }
            }
          }
        ],
        "hostIP": "100.88.107.46",
        "phase": "Running",
        "podIP": "100.88.231.95",
        "podIPs": [
          {
            "ip": "100.88.231.95"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-03-06T12:46:17Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-apiserver-5b7nn-mp7cd",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod28ea47f7-142a-482e-b672-cda12b01e55e",
          "cafe.sofastack.io/available-conditions-ex": "{\"expectedFinalizers\":{\"Service/ocmpaas/ocmpaastestcz30axvip0\":\"prot.cafe.sofastack.io/lbb_ocmpaastestcz30axvip0\"}}",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "c8f34a0ee98cb4f5a1aae048c43baa1c",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100083248042.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"ac489b33401f1eaa0000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"e5d7ab494ae33f5d\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"994acc74ff30deb9\"},{\"component\":\"default-scheduler\",\"span_id\":\"6385fdc9bc042dee\"},{\"component\":\"cni-service\",\"span_id\":\"145b03e8a6c88f29\"},{\"component\":\"kubelet\",\"span_id\":\"988eb5b3d02e53a5\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"8e5c4d94baa686a5\"},{\"component\":\"naming-controller\",\"span_id\":\"f497ec08db3431c9\"}],\"start_at\":\"2023-03-07T11:48:42+08:00\",\"finish_at\":\"2023-03-07T11:48:48+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":5120,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.83.248.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.83.248.42\",\"securityDomain\":\"ALI_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-apiserver-5b7nn-74777b5468",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-03-07T11:48:43.015762068+08:00",
          "pod.beta1.sigma.ali/trace-id": "92956cb6-f127-48f3-8c4d-fd558fd49808",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"92956cb6-f127-48f3-8c4d-fd558fd49808\",\"service\":\"naming\",\"creationTimestamp\":\"2023-03-07T11:48:47.372552933+08:00\",\"completionTimestamp\":\"2023-03-07T11:48:47.441549138+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:48:47.372552933+08:00\",\"endTimestamp\":\"2023-03-07T11:48:47.441549138+08:00\",\"logs\":[{\"time\":\"2023-03-07T11:48:47.372555621+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:48:47.372560192+08:00\",\"endTimestamp\":\"2023-03-07T11:48:47.441548815+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"92956cb6-f127-48f3-8c4d-fd558fd49808\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-03-07T11:48:47.124216762+08:00\",\"FinishTimestamp\":\"2023-03-07T11:48:47.19044454+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"92956cb6-f127-48f3-8c4d-fd558fd49808\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-07T11:48:48.948354895+08:00\",\"completionTimestamp\":\"2023-03-07T11:48:49.360285366+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:48:48.948354895+08:00\",\"endTimestamp\":\"2023-03-07T11:48:49.360285366+08:00\",\"logs\":[{\"time\":\"2023-03-07T11:48:48.948357277+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:48:48.948360581+08:00\",\"endTimestamp\":\"2023-03-07T11:48:49.360285175+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:48:48.948391446+08:00\",\"endTimestamp\":\"2023-03-07T11:48:49.360271136+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:48:48.94839305+08:00\",\"endTimestamp\":\"2023-03-07T11:48:48.948523262+08:00\"},{\"operation\":\"updateServerStatus\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:48:48.948525037+08:00\",\"endTimestamp\":\"2023-03-07T11:48:49.360270709+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-07T11:48:45.055766757+08:00\",\"finishTimestamp\":\"2023-03-07T11:48:45.879659523+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"unknown\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-apiserver-5b7nn-74777b5468\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-03-07 11:48:47.190474006 +0800 CST m=+2993862.011741189",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16815436\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-03-07T03:48:42Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn",
          "xvip.cafe.sofastack.io/xvip_rs_prot_ocmpaastestcz30axvip0"
        ],
        "generateName": "hub-apiserver-5b7nn-",
        "labels": {
          "ali.EnableDefaultRoute": "true",
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/just-created": "1678160922303324561",
          "cafe.sofastack.io/pod-ip": "100.83.248.42",
          "cafe.sofastack.io/pod-number": "4",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/version": "hub-apiserver-5b7nn-74777b5468",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-apiserver",
          "controller-revision-hash": "hub-apiserver-5b7nn-74777b5468",
          "lifecycle.cafe.sofastack.io/finish-upgrade": "1678160928946502620",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-apiserver-5b7nn-3100a2d0-41b4-472c-9bb8-7dcedbe9af6d",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100083248042.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100083248042",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "10m0s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "8C16G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1678160922333236681",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-183955",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.83.248.42",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "0ae37cfe-05ce-41b0-b9a1-2f8ca0e584f5",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "hub-apiserver-5b7nn-mp7cd",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-apiserver-5b7nn",
            "uid": "3100a2d0-41b4-472c-9bb8-7dcedbe9af6d"
          }
        ],
        "resourceVersion": "32427680793",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-apiserver-5b7nn-mp7cd",
        "uid": "28ea47f7-142a-482e-b672-cda12b01e55e"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-apiserver",
              "--allow-privileged=true",
              "--authorization-mode=Node,RBAC",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--endpoint-reconciler-type=none",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--enable-admission-plugins=NodeRestriction",
              "--enable-bootstrap-token-auth=true",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--insecure-port=0",
              "--secure-port=8443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--service-account-key-file=/etc/kubernetes/pki/apiserver/sa.pub",
              "--service-account-signing-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--service-account-issuer=api",
              "--api-audiences=api",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--proxy-client-cert-file=/etc/kubernetes/pki/apiserver/proxy.crt",
              "--proxy-client-key-file=/etc/kubernetes/pki/apiserver/proxy.key",
              "--log-file=/home/admin/logs/apiserver.log",
              "--log-file-max-size=100",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "SN",
                "value": "0ae37cfe-05ce-41b0-b9a1-2f8ca0e584f5"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "5"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "single-request-reopen"
            },
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "817366299",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.public.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "53687091200",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "ba332452-e6c5-4f0f-b3c3-b683cd2d55d9"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:42Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:43Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:47Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-03-07T03:48:48Z",
            "lastTransitionTime": "2023-03-07T03:48:48Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:43Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:50Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:46Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:43Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:43Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://645eb42a0d530f7b1a931a27c860255b3ae6c538d4d54015c2ba27d23585466d",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 0,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-07T03:48:45Z"
              }
            }
          }
        ],
        "hostIP": "100.69.120.224",
        "phase": "Running",
        "podIP": "100.83.248.42",
        "podIPs": [
          {
            "ip": "100.83.248.42"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-03-07T03:48:43Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-controller-manager-fmvm2-swxts",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pode7cfb730-7db0-4d18-9b71-f96b204539ae",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"controller-manager\",\"Resources\":{\"limits\":{\"cpu\":\"1\",\"ephemeral-storage\":\"20Gi\",\"memory\":\"2Gi\"},\"requests\":{\"cpu\":\"1\",\"ephemeral-storage\":\"20Gi\",\"memory\":\"2Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "7389963ce1357078e9cf99db0c35b16c",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100088116148.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"d71045289aa832a90000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"2c420f015fb44df3\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"8d5f410a17ccb0ef\"},{\"component\":\"default-scheduler\",\"span_id\":\"feaa50f02e111dd5\"},{\"component\":\"cni-service\",\"span_id\":\"56ac073ca4dc594c\"},{\"component\":\"kubelet\",\"span_id\":\"8d4c40e674783c17\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"7abe8d0af58b0e81\"},{\"component\":\"naming-controller\",\"span_id\":\"d64a725bc89e2b0d\"}],\"start_at\":\"2023-02-15T20:52:01+08:00\",\"finish_at\":\"2023-02-15T20:52:56+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"controller-manager\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":2147483648,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":1024,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.88.116.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.88.116.148\",\"securityDomain\":\"ALIPAY_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-controller-manager-fmvm2-fd47d46b4",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-02-15T20:52:01.952047036+08:00",
          "pod.beta1.sigma.ali/trace-id": "8095ba02-c5b9-411a-a82d-2db106ef5909",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"8095ba02-c5b9-411a-a82d-2db106ef5909\",\"service\":\"naming\",\"creationTimestamp\":\"2023-02-15T20:52:55.43580956+08:00\",\"completionTimestamp\":\"2023-02-15T20:52:55.501214882+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:52:55.43580956+08:00\",\"endTimestamp\":\"2023-02-15T20:52:55.501214882+08:00\",\"logs\":[{\"time\":\"2023-02-15T20:52:55.435813201+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:52:55.435816607+08:00\",\"endTimestamp\":\"2023-02-15T20:52:55.501214453+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"8095ba02-c5b9-411a-a82d-2db106ef5909\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-02-15T20:52:55.124158519+08:00\",\"FinishTimestamp\":\"2023-02-15T20:52:55.209944461+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"8095ba02-c5b9-411a-a82d-2db106ef5909\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-02-15T20:52:58.836774165+08:00\",\"completionTimestamp\":\"2023-02-15T20:52:58.83700311+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:52:58.836774165+08:00\",\"endTimestamp\":\"2023-02-15T20:52:58.83700311+08:00\",\"logs\":[{\"time\":\"2023-02-15T20:52:58.836776787+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:52:58.83678026+08:00\",\"endTimestamp\":\"2023-02-15T20:52:58.83700292+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:52:58.8368428+08:00\",\"endTimestamp\":\"2023-02-15T20:52:58.836987171+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:52:58.836844738+08:00\",\"endTimestamp\":\"2023-02-15T20:52:58.836987001+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"controller-manager\":{\"creationTimestamp\":\"2023-03-06T20:38:17.986192872+08:00\",\"finishTimestamp\":\"2023-03-06T20:38:18.931110115+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"exited\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-controller-manager-fmvm2-fd47d46b4\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-02-15 20:52:55.210025761 +0800 CST m=+1298510.031292953",
          "sigma.ali/container-diskQuotaID": "{\"controller-manager\":\"17607899\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202302150851427168715\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-02-15T12:52:01Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn"
        ],
        "generateName": "hub-controller-manager-fmvm2-",
        "labels": {
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "100.88.116.148",
          "cafe.sofastack.io/pod-number": "0",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1678106299117288625",
          "cafe.sofastack.io/version": "hub-controller-manager-fmvm2-fd47d46b4",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-controller-manager",
          "controller-revision-hash": "hub-controller-manager-fmvm2-fd47d46b4",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-controller-manager-fmvm2-94cb7efb-f3b1-4715-b86a-808c1d052bd5",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100088116148.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100088116148",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "30s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "1C2G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1676465521608936384",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-160722",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.88.116.148",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "9962e23f-40e4-43a6-8293-ec797574f66b",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202302150851427168715-1"
        },
        "name": "hub-controller-manager-fmvm2-swxts",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-controller-manager-fmvm2",
            "uid": "94cb7efb-f3b1-4715-b86a-808c1d052bd5"
          }
        ],
        "resourceVersion": "32427228128",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-controller-manager-fmvm2-swxts",
        "uid": "e7cfb730-7db0-4d18-9b71-f96b204539ae"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-controller-manager",
              "--kubeconfig=/etc/kubernetes/config/hub-controller-manager.kubeconfig",
              "--authentication-kubeconfig=/etc/kubernetes/config/hub-controller-manager.kubeconfig",
              "--authorization-kubeconfig=/etc/kubernetes/config/hub-controller-manager.kubeconfig",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--bind-address=127.0.0.1",
              "--cluster-name=hub-cluster",
              "--controllers=*",
              "--leader-elect=true",
              "--port=0",
              "--service-account-private-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--use-service-account-credentials=true",
              "--cluster-signing-cert-file=/etc/kubernetes/config/ca.crt",
              "--cluster-signing-key-file=/etc/kubernetes/config/ca.key"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "SN",
                "value": "9962e23f-40e4-43a6-8293-ec797574f66b"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "1"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "1"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "1"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "controller-manager",
            "resources": {
              "limits": {
                "cpu": "1",
                "ephemeral-storage": "20Gi",
                "memory": "2Gi"
              },
              "requests": {
                "cpu": "1",
                "ephemeral-storage": "20Gi",
                "memory": "2Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/config/",
                "name": "hub-controller-manager-kubeconfig",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            },
            {
              "name": "single-request-reopen"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "217337505",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "hub-controller-manager-kubeconfig",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-controller-manager-kubeconfig"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                },
                {
                  "key": "namespace",
                  "path": "namespace"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "21474836480",
                "ulogfs.k8s.alipay.com/enable-zclean": "",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "411f6f7f-3ca3-404f-897f-80eb842562f8"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "1Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:52:01Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:52:02Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:52:55Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-02-15T12:52:56Z",
            "lastTransitionTime": "2023-02-15T12:52:56Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:52:02Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:19Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:19Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:52:02Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:52:02Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://7ce3952fe34c2fc2e7c370402382e8d20c803c69c846043e49f6b3e7706df363",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "controller-manager",
            "ready": true,
            "restartCount": 1,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-06T12:38:18Z"
              }
            }
          }
        ],
        "hostIP": "100.83.13.147",
        "phase": "Running",
        "podIP": "100.88.116.148",
        "podIPs": [
          {
            "ip": "100.88.116.148"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-02-15T12:52:02Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-apiserver-5b7nn-5qqfk",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod679a5f5c-a37b-47cd-aac4-039477afb3c6",
          "cafe.sofastack.io/available-conditions-ex": "{\"expectedFinalizers\":{\"Service/ocmpaas/ocmpaastestcz30axvip0\":\"prot.cafe.sofastack.io/lbb_ocmpaastestcz30axvip0\"}}",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "1550600043c1f91603d581eb7a2a1000",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100081074009.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"d60765743383adf00000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"b87a5b0edaf3d00d\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"35233572715fe6dd\"},{\"component\":\"default-scheduler\",\"span_id\":\"4ed0f1b49cba9ea2\"},{\"component\":\"cni-service\",\"span_id\":\"89e31830f8d83036\"},{\"component\":\"kubelet\",\"span_id\":\"17de19eb35be5da6\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"3878bfb63b63825a\"},{\"component\":\"naming-controller\",\"span_id\":\"12ff92b4cb8313d1\"}],\"start_at\":\"2023-03-07T11:48:42+08:00\",\"finish_at\":\"2023-03-07T11:49:34+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "paascore.alipay.com/upgrade-diff": "77771c94bd5e87ddeb6289afad6ee70e",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":5120,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.81.74.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.81.74.9\",\"securityDomain\":\"ALI_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-apiserver-5b7nn-74777b5468",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-03-07T11:48:42.50480558+08:00",
          "pod.beta1.sigma.ali/trace-id": "2ff49ad7-3fe9-4902-9daa-153ccd811176",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"2ff49ad7-3fe9-4902-9daa-153ccd811176\",\"service\":\"naming\",\"creationTimestamp\":\"2023-03-07T11:49:32.686548376+08:00\",\"completionTimestamp\":\"2023-03-07T11:49:32.757961503+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:49:32.686548376+08:00\",\"endTimestamp\":\"2023-03-07T11:49:32.757961503+08:00\",\"logs\":[{\"time\":\"2023-03-07T11:49:32.686552668+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:49:32.68655754+08:00\",\"endTimestamp\":\"2023-03-07T11:49:32.757961192+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"2ff49ad7-3fe9-4902-9daa-153ccd811176\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-03-07T11:49:32.375517349+08:00\",\"FinishTimestamp\":\"2023-03-07T11:49:32.434505692+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"2ff49ad7-3fe9-4902-9daa-153ccd811176\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-07T11:49:34.431663551+08:00\",\"completionTimestamp\":\"2023-03-07T11:49:34.488243644+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:49:34.431663551+08:00\",\"endTimestamp\":\"2023-03-07T11:49:34.488243644+08:00\",\"logs\":[{\"time\":\"2023-03-07T11:49:34.43166899+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:49:34.431673947+08:00\",\"endTimestamp\":\"2023-03-07T11:49:34.488243344+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:49:34.431725058+08:00\",\"endTimestamp\":\"2023-03-07T11:49:34.488219812+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:49:34.43172746+08:00\",\"endTimestamp\":\"2023-03-07T11:49:34.431897324+08:00\"},{\"operation\":\"updateServerStatus\",\"success\":true,\"startTimestamp\":\"2023-03-07T11:49:34.431900425+08:00\",\"endTimestamp\":\"2023-03-07T11:49:34.488219319+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-07T11:48:44.331205128+08:00\",\"finishTimestamp\":\"2023-03-07T11:49:31.248298087+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"unknown\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-apiserver-5b7nn-74777b5468\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-03-07 11:49:32.434542009 +0800 CST m=+2993907.255809193",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16815829\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-03-07T03:48:42Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn",
          "xvip.cafe.sofastack.io/xvip_rs_prot_ocmpaastestcz30axvip0"
        ],
        "generateName": "hub-apiserver-5b7nn-",
        "labels": {
          "ali.EnableDefaultRoute": "true",
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/just-created": "1678160922202408006",
          "cafe.sofastack.io/pod-ip": "100.81.74.9",
          "cafe.sofastack.io/pod-number": "2",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/version": "hub-apiserver-5b7nn-74777b5468",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-apiserver",
          "controller-revision-hash": "hub-apiserver-5b7nn-74777b5468",
          "lifecycle.cafe.sofastack.io/finish-upgrade": "1678160974429841189",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-apiserver-5b7nn-3100a2d0-41b4-472c-9bb8-7dcedbe9af6d",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100081074009.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100081074009",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "10m0s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "8C16G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1678160922231885636",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-183536",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.81.74.9",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "f54dea9b-3f86-4881-8c16-34cd20e88565",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "hub-apiserver-5b7nn-5qqfk",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-apiserver-5b7nn",
            "uid": "3100a2d0-41b4-472c-9bb8-7dcedbe9af6d"
          }
        ],
        "resourceVersion": "32427669902",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-apiserver-5b7nn-5qqfk",
        "uid": "679a5f5c-a37b-47cd-aac4-039477afb3c6"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-apiserver",
              "--allow-privileged=true",
              "--authorization-mode=Node,RBAC",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--endpoint-reconciler-type=none",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--enable-admission-plugins=NodeRestriction",
              "--enable-bootstrap-token-auth=true",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--insecure-port=0",
              "--secure-port=8443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--service-account-key-file=/etc/kubernetes/pki/apiserver/sa.pub",
              "--service-account-signing-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--service-account-issuer=api",
              "--api-audiences=api",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--proxy-client-cert-file=/etc/kubernetes/pki/apiserver/proxy.crt",
              "--proxy-client-key-file=/etc/kubernetes/pki/apiserver/proxy.key",
              "--log-file=/home/admin/logs/apiserver.log",
              "--log-file-max-size=100",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "SN",
                "value": "f54dea9b-3f86-4881-8c16-34cd20e88565"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "5"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "single-request-reopen"
            },
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "817385102",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.public.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                },
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "53687091200",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "1c0e9646-93f6-49be-b68a-f56f750e8fc2"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:42Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:42Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:49:32Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-03-07T03:49:34Z",
            "lastTransitionTime": "2023-03-07T03:49:34Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:42Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:49:36Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:49:32Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:42Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-07T03:48:42Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://48722ae85be72e2cdb6428ae2af543a35ecd9ecf5d9b68b9c571c0648ace4803",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 0,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-07T03:49:31Z"
              }
            }
          }
        ],
        "hostIP": "100.81.208.115",
        "phase": "Running",
        "podIP": "100.81.74.9",
        "podIPs": [
          {
            "ip": "100.81.74.9"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-03-07T03:48:42Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-apiserver-5b7nn-r6xgm",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod7a38b23e-815f-43b9-bb0e-4655392fa7ba",
          "cafe.sofastack.io/available-conditions-ex": "{}",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "1ed984dca79f51acbd0efecd774d0500",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100069200113.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"0bfa53f867ee52a20000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"03cd07a40afaa181\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"035d146effe00fd5\"},{\"component\":\"default-scheduler\",\"span_id\":\"dd408479e6c94efc\"},{\"component\":\"cni-service\",\"span_id\":\"dd4428ee36c8ab5d\"},{\"component\":\"kubelet\",\"span_id\":\"6361d5142a2600bf\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"f9629ff5a1a565ad\"},{\"component\":\"naming-controller\",\"span_id\":\"0a1742ce4074da98\"}],\"start_at\":\"2023-03-06T20:38:03+08:00\",\"finish_at\":\"2023-03-06T20:38:48+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "paascore.alipay.com/upgrade-diff": "77771c94bd5e87ddeb6289afad6ee70e",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":5120,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.69.200.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.69.200.113\",\"securityDomain\":\"ALI_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-apiserver-5b7nn-74777b5468",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-03-06T20:38:03.543787198+08:00",
          "pod.beta1.sigma.ali/trace-id": "07e72db0-6009-4a01-9b5c-b7cf76fe2af7",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"07e72db0-6009-4a01-9b5c-b7cf76fe2af7\",\"service\":\"naming\",\"creationTimestamp\":\"2023-03-06T20:38:48.285011353+08:00\",\"completionTimestamp\":\"2023-03-06T20:38:48.358344663+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:48.285011353+08:00\",\"endTimestamp\":\"2023-03-06T20:38:48.358344663+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:38:48.285014533+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:48.285020179+08:00\",\"endTimestamp\":\"2023-03-06T20:38:48.358344124+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"07e72db0-6009-4a01-9b5c-b7cf76fe2af7\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-03-06T20:38:47.946855428+08:00\",\"FinishTimestamp\":\"2023-03-06T20:38:48.003571817+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"07e72db0-6009-4a01-9b5c-b7cf76fe2af7\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-06T20:38:51.851085375+08:00\",\"completionTimestamp\":\"2023-03-06T20:38:51.851269696+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:51.851085375+08:00\",\"endTimestamp\":\"2023-03-06T20:38:51.851269696+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:38:51.851087522+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:51.851090355+08:00\",\"endTimestamp\":\"2023-03-06T20:38:51.85126959+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:51.851125877+08:00\",\"endTimestamp\":\"2023-03-06T20:38:51.851255712+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:51.851127534+08:00\",\"endTimestamp\":\"2023-03-06T20:38:51.851255556+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-06T20:38:05.347610222+08:00\",\"finishTimestamp\":\"2023-03-06T20:38:46.634776725+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"unknown\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-apiserver-5b7nn-74777b5468\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-03-06 20:38:48.003635487 +0800 CST m=+2939262.824902693",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16828080\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-03-06T12:38:03Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn",
          "xvip.cafe.sofastack.io/xvip_rs_prot_ocmpaastestcz30axvip0"
        ],
        "generateName": "hub-apiserver-5b7nn-",
        "labels": {
          "ali.EnableDefaultRoute": "true",
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "100.69.200.113",
          "cafe.sofastack.io/pod-number": "4",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1678106594139910611",
          "cafe.sofastack.io/version": "hub-apiserver-5b7nn-74777b5468",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-apiserver1",
          "controller-revision-hash": "hub-apiserver-5b7nn-74777b5468",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-apiserver-5b7nn-3100a2d0-41b4-472c-9bb8-7dcedbe9af6d",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100069200113.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100069200113",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "10m0s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "8C16G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1678106283078326078",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-165552",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.69.200.113",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "6ac3a06f-6bfa-4955-a0cf-b8ffdbe2eb04",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "hub-apiserver-5b7nn-r6xgm",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-apiserver-5b7nn",
            "uid": "3100a2d0-41b4-472c-9bb8-7dcedbe9af6d"
          }
        ],
        "resourceVersion": "32427373352",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-apiserver-5b7nn-r6xgm",
        "uid": "7a38b23e-815f-43b9-bb0e-4655392fa7ba"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-apiserver",
              "--allow-privileged=true",
              "--authorization-mode=Node,RBAC",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--endpoint-reconciler-type=none",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--enable-admission-plugins=NodeRestriction",
              "--enable-bootstrap-token-auth=true",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--insecure-port=0",
              "--secure-port=8443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--service-account-key-file=/etc/kubernetes/pki/apiserver/sa.pub",
              "--service-account-signing-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--service-account-issuer=api",
              "--api-audiences=api",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--proxy-client-cert-file=/etc/kubernetes/pki/apiserver/proxy.crt",
              "--proxy-client-key-file=/etc/kubernetes/pki/apiserver/proxy.key",
              "--log-file=/home/admin/logs/apiserver.log",
              "--log-file-max-size=100",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "SN",
                "value": "6ac3a06f-6bfa-4955-a0cf-b8ffdbe2eb04"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "5"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "single-request-reopen"
            },
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "cv80367l05e",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.public.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                },
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "53687091200",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "c5006f59-a158-4061-9429-3c16d737ab42"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:03Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:03Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:48Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-03-06T12:38:48Z",
            "lastTransitionTime": "2023-03-06T12:38:48Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:03Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:50Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:47Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:03Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:03Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://af87c25c69173e366a164daf8539790d17ac5c926449f02294bee7af4932c899",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 0,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-06T12:38:46Z"
              }
            }
          }
        ],
        "hostIP": "100.81.235.90",
        "phase": "Running",
        "podIP": "100.69.200.113",
        "podIPs": [
          {
            "ip": "100.69.200.113"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-03-06T12:38:03Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "cluster-extension-apiserver-jqntf-vdnc4",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod2a80b4e9-2023-4c95-b20a-bb28a8da3851",
          "cafe.sofastack.io/pre-check-timestamp": "1678106201455700450",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"approveAndReset\",\"timestamp\":\"2023-03-06T12:36:41.45570045Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"4\",\"ephemeral-storage\":\"100Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"4\",\"ephemeral-storage\":\"100Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "9e44a37ca42d239142304071a5583805",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100083098040.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"734f52b65c3d69bf0000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"ac99eeb1a807c523\",\"delivery_type\":\"PodUpgrade\",\"status\":\"closed\",\"services\":[{\"component\":\"kubelet\",\"span_id\":\"a2ba2908f0c4cd11\"}],\"start_at\":\"2023-03-06T20:36:41+08:00\",\"finish_at\":\"2023-03-06T20:36:49+08:00\",\"extra_info\":{\"upgrade_containers\":\"apiserver\"}},{\"trace_id\":\"dd8e3482a5228ce80000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"95485ce79c9bcb51\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"2baa16a188f6e9bf\"},{\"component\":\"default-scheduler\",\"span_id\":\"3da30b0c36a74ef7\"},{\"component\":\"cni-service\",\"span_id\":\"d15e19c7ac6ef4f7\"},{\"component\":\"kubelet\",\"span_id\":\"5d57e1d3cd6cbcb4\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"c6e6d022fa57684b\"},{\"component\":\"naming-controller\",\"span_id\":\"03cc5fcd6618836a\"}],\"start_at\":\"2023-02-15T20:20:53+08:00\",\"finish_at\":\"2023-02-15T20:21:02+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "paascore.alipay.com/upgrade-diff": "aade8d86cc516ef329bbba9c64fefc00",
          "pod.beta1.alipay.com/request-action": "{\"action-type\":\"RequestUpgrade\",\"containers\":[\"apiserver\"],\"timestamp\":\"2023-03-06T20:36:41.819573062+08:00\"}",
          "pod.beta1.alipay.com/upgrade-reason": "apiserver:img",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":4096,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.83.98.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.83.98.40\",\"securityDomain\":\"ALIPAY_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "cluster-extension-apiserver-jqntf-694cf69554",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-02-15T20:20:53.992535406+08:00",
          "pod.beta1.sigma.ali/trace-id": "60f0d2b5-fc5e-45e0-9925-bb20cde908fb",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"60f0d2b5-fc5e-45e0-9925-bb20cde908fb\",\"service\":\"naming\",\"creationTimestamp\":\"2023-02-15T20:21:01.220856996+08:00\",\"completionTimestamp\":\"2023-02-15T20:21:01.295015006+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:21:01.220856996+08:00\",\"endTimestamp\":\"2023-02-15T20:21:01.295015006+08:00\",\"logs\":[{\"time\":\"2023-02-15T20:21:01.220860651+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:21:01.220866381+08:00\",\"endTimestamp\":\"2023-02-15T20:21:01.295014632+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"60f0d2b5-fc5e-45e0-9925-bb20cde908fb\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-02-15T20:21:00.90303777+08:00\",\"FinishTimestamp\":\"2023-02-15T20:21:00.984969584+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"60f0d2b5-fc5e-45e0-9925-bb20cde908fb\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-06T20:36:49.9861776+08:00\",\"completionTimestamp\":\"2023-03-06T20:36:49.98645243+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:36:49.9861776+08:00\",\"endTimestamp\":\"2023-03-06T20:36:49.98645243+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:36:49.986179633+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:36:49.986182598+08:00\",\"endTimestamp\":\"2023-03-06T20:36:49.986452321+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:36:49.986305064+08:00\",\"endTimestamp\":\"2023-03-06T20:36:49.986437094+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:36:49.986307432+08:00\",\"endTimestamp\":\"2023-03-06T20:36:49.986436924+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-06T20:36:43.81989254+08:00\",\"finishTimestamp\":\"2023-03-06T20:36:48.549262736+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"running\",\"action\":\"upgrade\",\"success\":true,\"message\":\"upgrade container success\",\"specHash\":\"cluster-extension-apiserver-jqntf-694cf69554\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-02-15 20:21:00.985056172 +0800 CST m=+1296595.806323377",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16855159\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-02-15T12:20:53Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn"
        ],
        "generateName": "cluster-extension-apiserver-jqntf-",
        "labels": {
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "100.83.98.40",
          "cafe.sofastack.io/pod-number": "2",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1678106210073206617",
          "cafe.sofastack.io/version": "cluster-extension-apiserver-jqntf-694cf69554",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "cluster-extension-apiserver",
          "controller-revision-hash": "cluster-extension-apiserver-jqntf-694cf69554",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "cluster-extension-apiserver-jqntf-30bd4911-a01c-4bd8-be1d-48dc67520374",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100083098040.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100083098040",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "30s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "4C8G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "operation.cafe.sofastack.io/inplaceset": "1677586282941133356",
          "paascore.alipay.com/adopted": "1676463653760431313",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-172647",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.83.98.40",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "ac2f784f-1516-40c5-99c1-e909326d7faa",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "cluster-extension-apiserver-jqntf-vdnc4",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "cluster-extension-apiserver-jqntf",
            "uid": "30bd4911-a01c-4bd8-be1d-48dc67520374"
          }
        ],
        "resourceVersion": "32427464188",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/cluster-extension-apiserver-jqntf-vdnc4",
        "uid": "2a80b4e9-2023-4c95-b20a-bb28a8da3851"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "./apiserver",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--disable-admission-plugins=MutatingAdmissionWebhook,NamespaceLifecycle,ValidatingAdmissionWebhook",
              "--feature-gates=APIPriorityAndFairness=false",
              "--secure-port=7443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--authorization-kubeconfig=/etc/kubernetes/pki/apiserver/cluster-extension-delegate.kubeconfig",
              "--authentication-kubeconfig=/etc/kubernetes/pki/apiserver/cluster-extension-delegate.kubeconfig",
              "--unified-identity-enabled=true",
              "--unified-identity-cert-file=/var/run/secrets/kubernetes.io/serviceaccount/app.crt",
              "--unified-identity-key-file=/var/run/secrets/kubernetes.io/serviceaccount/app.key",
              "--authorization-always-allow-paths=\"/metrics\"",
              "--audit-log-path=/home/admin/logs/audit.log",
              "--audit-policy-file=/etc/kubernetes/pki/apiserver/audit.yaml",
              "--audit-log-maxage=3",
              "--audit-log-maxsize=10240",
              "--log_file=/home/admin/logs/extension.log",
              "--log_file_max_size=200",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "SN",
                "value": "ac2f784f-1516-40c5-99c1-e909326d7faa"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "4"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "4"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "4"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              },
              {
                "name": "PARENT_SPEC_GENERATION",
                "value": "47"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver:v0.3.5",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "4",
                "ephemeral-storage": "100Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "4",
                "ephemeral-storage": "100Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "single-request-reopen"
            },
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "cv80367l0k0",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                },
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver:v0.3.1",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "107374182400",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "ae78a812-09e2-4d14-b0fc-7da87f6c68ae"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:36:49Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:20:54Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:01Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-02-15T12:21:02Z",
            "lastTransitionTime": "2023-02-15T12:21:02Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:20:54Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:36:50Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:00Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:20:54Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:20:54Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://fe6b5ead2e02cb12c05ec4a3448a72c9240d7343fd3a9ae8870f1584366a9339",
            "image": "reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver:v0.3.5",
            "imageID": "reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver@sha256:82f14d0ac1a325654256e7e82ee304a2270154bea62bbb34b2b02ae467ddce80",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 3,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-06T12:36:48Z"
              }
            }
          }
        ],
        "hostIP": "100.88.121.225",
        "phase": "Running",
        "podIP": "100.83.98.40",
        "podIPs": [
          {
            "ip": "100.83.98.40"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-02-15T12:20:54Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-apiserver-5b7nn-96q9d",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod54273f54-1962-4e2d-9289-5ac9d54f68b0",
          "cafe.sofastack.io/available-conditions-ex": "{}",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "efa23d74ff3ec6548a74af5e59660be0",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-011167047225.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"a260db5b22fadf9f0000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"951b86f608e1bdc8\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"22cf7cd2e61152b6\"},{\"component\":\"default-scheduler\",\"span_id\":\"64f017e50428db34\"},{\"component\":\"cni-service\",\"span_id\":\"1e70cecc624554d0\"},{\"component\":\"kubelet\",\"span_id\":\"fb0cd52c86213d5f\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"c8ea8a9e020450d2\"},{\"component\":\"naming-controller\",\"span_id\":\"4085b61781184da9\"}],\"start_at\":\"2023-03-06T20:38:03+08:00\",\"finish_at\":\"2023-03-06T20:38:48+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "paascore.alipay.com/upgrade-diff": "77771c94bd5e87ddeb6289afad6ee70e",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":5120,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"11.167.47.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"11.167.47.225\",\"securityDomain\":\"ALI_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-apiserver-5b7nn-74777b5468",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-03-06T20:38:04.108728427+08:00",
          "pod.beta1.sigma.ali/trace-id": "445ab0f0-1926-44f7-a2a0-15d2b3917dce",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"445ab0f0-1926-44f7-a2a0-15d2b3917dce\",\"service\":\"naming\",\"creationTimestamp\":\"2023-03-06T20:38:48.218011543+08:00\",\"completionTimestamp\":\"2023-03-06T20:38:48.28860073+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:48.218011543+08:00\",\"endTimestamp\":\"2023-03-06T20:38:48.28860073+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:38:48.218013944+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:48.218017061+08:00\",\"endTimestamp\":\"2023-03-06T20:38:48.288600474+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"445ab0f0-1926-44f7-a2a0-15d2b3917dce\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-03-06T20:38:47.868783502+08:00\",\"FinishTimestamp\":\"2023-03-06T20:38:47.940432971+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"445ab0f0-1926-44f7-a2a0-15d2b3917dce\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-06T20:38:49.236970305+08:00\",\"completionTimestamp\":\"2023-03-06T20:38:49.237389581+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:49.236970305+08:00\",\"endTimestamp\":\"2023-03-06T20:38:49.237389581+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:38:49.236972528+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:49.236975478+08:00\",\"endTimestamp\":\"2023-03-06T20:38:49.237389095+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:49.237137981+08:00\",\"endTimestamp\":\"2023-03-06T20:38:49.237363074+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:49.237141572+08:00\",\"endTimestamp\":\"2023-03-06T20:38:49.237362845+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-06T20:38:05.608176709+08:00\",\"finishTimestamp\":\"2023-03-06T20:38:46.523540378+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"unknown\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-apiserver-5b7nn-74777b5468\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-03-06 20:38:47.94046427 +0800 CST m=+2939262.761731459",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16810390\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-03-06T12:38:03Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn",
          "xvip.cafe.sofastack.io/xvip_rs_prot_ocmpaastestcz30axvip0"
        ],
        "generateName": "hub-apiserver-5b7nn-",
        "labels": {
          "ali.EnableDefaultRoute": "true",
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "11.167.47.225",
          "cafe.sofastack.io/pod-number": "3",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1678106777779342658",
          "cafe.sofastack.io/version": "hub-apiserver-5b7nn-74777b5468",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-apiserver1",
          "controller-revision-hash": "hub-apiserver-5b7nn-74777b5468",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-apiserver-5b7nn-3100a2d0-41b4-472c-9bb8-7dcedbe9af6d",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-011167047225.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-011167047225",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "10m0s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "8C16G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1678106283541691661",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-185944",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "11.167.47.225",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "41fded4a-7847-4f51-9984-b74e9a1fe876",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "hub-apiserver-5b7nn-96q9d",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-apiserver-5b7nn",
            "uid": "3100a2d0-41b4-472c-9bb8-7dcedbe9af6d"
          }
        ],
        "resourceVersion": "32427732582",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-apiserver-5b7nn-96q9d",
        "uid": "54273f54-1962-4e2d-9289-5ac9d54f68b0"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-apiserver",
              "--allow-privileged=true",
              "--authorization-mode=Node,RBAC",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--endpoint-reconciler-type=none",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--enable-admission-plugins=NodeRestriction",
              "--enable-bootstrap-token-auth=true",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--insecure-port=0",
              "--secure-port=8443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--service-account-key-file=/etc/kubernetes/pki/apiserver/sa.pub",
              "--service-account-signing-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--service-account-issuer=api",
              "--api-audiences=api",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--proxy-client-cert-file=/etc/kubernetes/pki/apiserver/proxy.crt",
              "--proxy-client-key-file=/etc/kubernetes/pki/apiserver/proxy.key",
              "--log-file=/home/admin/logs/apiserver.log",
              "--log-file-max-size=100",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "SN",
                "value": "41fded4a-7847-4f51-9984-b74e9a1fe876"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "5"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "single-request-reopen"
            },
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "cv80367l03t",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.public.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "53687091200",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "661af6aa-6a96-4f48-92cd-c30fc5708a10"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:03Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:04Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:48Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-03-06T12:38:48Z",
            "lastTransitionTime": "2023-03-06T12:38:48Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:04Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:49Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:47Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:04Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:04Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://298b5059574a8fb53cf6c68c19baa5c4ec37983313217e74e6203f19e701e3d8",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 0,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-06T12:38:46Z"
              }
            }
          }
        ],
        "hostIP": "100.69.120.27",
        "phase": "Running",
        "podIP": "11.167.47.225",
        "podIPs": [
          {
            "ip": "11.167.47.225"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-03-06T12:38:04Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-apiserver-5b7nn-8q8f9",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/podfaa2f600-5dfd-4d93-9b71-4c64e9caecdc",
          "cafe.sofastack.io/available-conditions-ex": "{}",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"5\",\"ephemeral-storage\":\"50Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "8b1b27b2513c4a8ea29a658dd5d4fb29",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100081038229.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"d137c5d2429e88050000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"c17c9d4855dda0a8\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"ffc63fd3fd150f46\"},{\"component\":\"default-scheduler\",\"span_id\":\"0e2771548b325f69\"},{\"component\":\"cni-service\",\"span_id\":\"b17e079186b4f74a\"},{\"component\":\"kubelet\",\"span_id\":\"cc2caae2a84fa9ae\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"5e3cb7da56ac6cb4\"},{\"component\":\"naming-controller\",\"span_id\":\"1ce18c3ecc4652dc\"}],\"start_at\":\"2023-03-06T20:37:59+08:00\",\"finish_at\":\"2023-03-06T20:38:49+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "paascore.alipay.com/upgrade-diff": "77771c94bd5e87ddeb6289afad6ee70e",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":5120,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.81.38.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.81.38.229\",\"securityDomain\":\"ALI_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-apiserver-5b7nn-74777b5468",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-03-06T20:38:00.034804858+08:00",
          "pod.beta1.sigma.ali/trace-id": "f44dae83-9ff8-4f6b-b66d-87013af53ef9",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"f44dae83-9ff8-4f6b-b66d-87013af53ef9\",\"service\":\"naming\",\"creationTimestamp\":\"2023-03-06T20:38:48.563104925+08:00\",\"completionTimestamp\":\"2023-03-06T20:38:48.631050159+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:48.563104925+08:00\",\"endTimestamp\":\"2023-03-06T20:38:48.631050159+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:38:48.563108372+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:48.563112028+08:00\",\"endTimestamp\":\"2023-03-06T20:38:48.631049707+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"f44dae83-9ff8-4f6b-b66d-87013af53ef9\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-03-06T20:38:48.304075873+08:00\",\"FinishTimestamp\":\"2023-03-06T20:38:48.377595563+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"f44dae83-9ff8-4f6b-b66d-87013af53ef9\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-06T20:38:50.412310867+08:00\",\"completionTimestamp\":\"2023-03-06T20:38:50.412519206+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:50.412310867+08:00\",\"endTimestamp\":\"2023-03-06T20:38:50.412519206+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:38:50.412314986+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:50.412322988+08:00\",\"endTimestamp\":\"2023-03-06T20:38:50.412519098+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:50.412374978+08:00\",\"endTimestamp\":\"2023-03-06T20:38:50.412503998+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:38:50.412376674+08:00\",\"endTimestamp\":\"2023-03-06T20:38:50.412503847+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-06T20:38:01.914265064+08:00\",\"finishTimestamp\":\"2023-03-06T20:38:47.927455324+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"unknown\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-apiserver-5b7nn-74777b5468\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-03-06 20:38:48.377660997 +0800 CST m=+2939263.198928187",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16804021\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-03-06T12:37:59Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn",
          "xvip.cafe.sofastack.io/xvip_rs_prot_ocmpaastestcz30axvip0"
        ],
        "generateName": "hub-apiserver-5b7nn-",
        "labels": {
          "ali.EnableDefaultRoute": "true",
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "100.81.38.229",
          "cafe.sofastack.io/pod-number": "1",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1678106707209084675",
          "cafe.sofastack.io/version": "hub-apiserver-5b7nn-74777b5468",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-apiserver1",
          "controller-revision-hash": "hub-apiserver-5b7nn-74777b5468",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-apiserver-5b7nn-3100a2d0-41b4-472c-9bb8-7dcedbe9af6d",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100081038229.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100081038229",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "10m0s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "8C16G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1678106279715515883",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-164734",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.81.38.229",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "912392d5-f42f-4fc1-8c68-8410a87bbb9d",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "hub-apiserver-5b7nn-8q8f9",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-apiserver-5b7nn",
            "uid": "3100a2d0-41b4-472c-9bb8-7dcedbe9af6d"
          }
        ],
        "resourceVersion": "32427349053",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-apiserver-5b7nn-8q8f9",
        "uid": "faa2f600-5dfd-4d93-9b71-4c64e9caecdc"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-apiserver",
              "--allow-privileged=true",
              "--authorization-mode=Node,RBAC",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--endpoint-reconciler-type=none",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--enable-admission-plugins=NodeRestriction",
              "--enable-bootstrap-token-auth=true",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--insecure-port=0",
              "--secure-port=8443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--service-account-key-file=/etc/kubernetes/pki/apiserver/sa.pub",
              "--service-account-signing-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--service-account-issuer=api",
              "--api-audiences=api",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--proxy-client-cert-file=/etc/kubernetes/pki/apiserver/proxy.crt",
              "--proxy-client-key-file=/etc/kubernetes/pki/apiserver/proxy.key",
              "--log-file=/home/admin/logs/apiserver.log",
              "--log-file-max-size=100",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "SN",
                "value": "912392d5-f42f-4fc1-8c68-8410a87bbb9d"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "5"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "5"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "5",
                "ephemeral-storage": "50Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            },
            {
              "name": "single-request-reopen"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "817382220",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.public.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                },
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "53687091200",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "20299db1-5b26-4e6e-85ba-d4d078156bde"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:37:59Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:00Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:48Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-03-06T12:38:49Z",
            "lastTransitionTime": "2023-03-06T12:38:49Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:00Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:51Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:48Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:00Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:38:00Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://360a923f5fa092c3373c92d8d24b339d0a3bd33ccce0145a0b487260149f065f",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 0,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-06T12:38:47Z"
              }
            }
          }
        ],
        "hostIP": "100.81.231.48",
        "phase": "Running",
        "podIP": "100.81.38.229",
        "podIPs": [
          {
            "ip": "100.81.38.229"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-03-06T12:38:00Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "cluster-extension-apiserver-jqntf-8mn2w",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod3ee65141-f99b-46c5-a24f-f114de36614d",
          "cafe.sofastack.io/pre-check-timestamp": "1678106201116566568",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"approveAndReset\",\"timestamp\":\"2023-03-06T12:36:41.116566568Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"apiserver\",\"Resources\":{\"limits\":{\"cpu\":\"4\",\"ephemeral-storage\":\"100Gi\",\"memory\":\"8Gi\"},\"requests\":{\"cpu\":\"4\",\"ephemeral-storage\":\"100Gi\",\"memory\":\"8Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "b29a0988bf33d6db246159461c03d076",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-011166253069.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"a3209aff0e8998760000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"542f18e78c41f0ab\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"387b3d461c05300d\"},{\"component\":\"default-scheduler\",\"span_id\":\"6c913c0350c00793\"},{\"component\":\"cni-service\",\"span_id\":\"2a66ed0cdf555a01\"},{\"component\":\"kubelet\",\"span_id\":\"e2e173a5fa0d5a88\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"8713fa8ae03dd369\"},{\"component\":\"naming-controller\",\"span_id\":\"f703bdb5e8a0f381\"}],\"start_at\":\"2023-02-15T20:21:37+08:00\",\"finish_at\":\"2023-02-15T20:22:00+08:00\",\"extra_info\":null},{\"trace_id\":\"596dc5df596a2fdb0000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"6b01058ff2c4db5c\",\"delivery_type\":\"PodUpgrade\",\"status\":\"closed\",\"services\":[{\"component\":\"kubelet\",\"span_id\":\"23a8fd51b7940ece\"}],\"start_at\":\"2023-03-06T20:36:41+08:00\",\"finish_at\":\"2023-03-06T20:37:00+08:00\",\"extra_info\":{\"upgrade_containers\":\"apiserver\"}}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "paascore.alipay.com/upgrade-diff": "ed8c8ba81db527180a6adc11c48b1c95",
          "pod.beta1.alipay.com/request-action": "{\"action-type\":\"RequestUpgrade\",\"containers\":[\"apiserver\"],\"timestamp\":\"2023-03-06T20:36:41.476999873+08:00\"}",
          "pod.beta1.alipay.com/upgrade-reason": "apiserver:img",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"apiserver\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":8589934592,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":4096,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"11.166.253.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"11.166.253.69\",\"securityDomain\":\"ALIPAY_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "cluster-extension-apiserver-jqntf-694cf69554",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-02-15T20:21:37.98900591+08:00",
          "pod.beta1.sigma.ali/trace-id": "7ceddedc-6d74-4e32-8ad2-53dba7ce3b37",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"7ceddedc-6d74-4e32-8ad2-53dba7ce3b37\",\"service\":\"naming\",\"creationTimestamp\":\"2023-02-15T20:21:58.475992074+08:00\",\"completionTimestamp\":\"2023-02-15T20:21:58.542383276+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:21:58.475992074+08:00\",\"endTimestamp\":\"2023-02-15T20:21:58.542383276+08:00\",\"logs\":[{\"time\":\"2023-02-15T20:21:58.475993642+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:21:58.475996333+08:00\",\"endTimestamp\":\"2023-02-15T20:21:58.542382792+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"7ceddedc-6d74-4e32-8ad2-53dba7ce3b37\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-02-15T20:21:58.22711954+08:00\",\"FinishTimestamp\":\"2023-02-15T20:21:58.313447628+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"7ceddedc-6d74-4e32-8ad2-53dba7ce3b37\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-03-06T20:37:01.236478197+08:00\",\"completionTimestamp\":\"2023-03-06T20:37:01.236824052+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:37:01.236478197+08:00\",\"endTimestamp\":\"2023-03-06T20:37:01.236824052+08:00\",\"logs\":[{\"time\":\"2023-03-06T20:37:01.236482825+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:37:01.236485623+08:00\",\"endTimestamp\":\"2023-03-06T20:37:01.236823802+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:37:01.236644137+08:00\",\"endTimestamp\":\"2023-03-06T20:37:01.236803572+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-03-06T20:37:01.236646566+08:00\",\"endTimestamp\":\"2023-03-06T20:37:01.236803423+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"apiserver\":{\"creationTimestamp\":\"2023-03-06T20:36:43.024862905+08:00\",\"finishTimestamp\":\"2023-03-06T20:36:58.909370042+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"running\",\"action\":\"upgrade\",\"success\":true,\"message\":\"upgrade container success\",\"specHash\":\"cluster-extension-apiserver-jqntf-694cf69554\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-02-15 20:21:58.313496217 +0800 CST m=+1296653.134763400",
          "sigma.ali/container-diskQuotaID": "{\"apiserver\":\"16805575\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202303060836185508747\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/enable-zclean": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-02-15T12:21:37Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn"
        ],
        "generateName": "cluster-extension-apiserver-jqntf-",
        "labels": {
          "a": "5",
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "11.166.253.69",
          "cafe.sofastack.io/pod-number": "4",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1678106221339975220",
          "cafe.sofastack.io/version": "cluster-extension-apiserver-jqntf-694cf69554",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "cluster-extension-apiserver",
          "controller-revision-hash": "cluster-extension-apiserver-jqntf-694cf69554",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "cluster-extension-apiserver-jqntf-30bd4911-a01c-4bd8-be1d-48dc67520374",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-011166253069.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-011166253069",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "30s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "4C8G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "operation.cafe.sofastack.io/inplaceset": "1677586282821177821",
          "paascore.alipay.com/adopted": "1676463697541280418",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-144146",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "11.166.253.69",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "e9a85417-88b2-41cb-800b-ea015b6dd5da",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202303060836185508747-1"
        },
        "name": "cluster-extension-apiserver-jqntf-8mn2w",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "cluster-extension-apiserver-jqntf",
            "uid": "30bd4911-a01c-4bd8-be1d-48dc67520374"
          }
        ],
        "resourceVersion": "32426974611",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/cluster-extension-apiserver-jqntf-8mn2w",
        "uid": "3ee65141-f99b-46c5-a24f-f114de36614d"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "./apiserver",
              "--client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--disable-admission-plugins=MutatingAdmissionWebhook,NamespaceLifecycle,ValidatingAdmissionWebhook",
              "--feature-gates=APIPriorityAndFairness=false",
              "--secure-port=7443",
              "--tls-cert-file=/etc/kubernetes/pki/apiserver/apiserver.crt",
              "--tls-private-key-file=/etc/kubernetes/pki/apiserver/apiserver.key",
              "--encryption-provider-config=/etc/kubernetes/pki/apiserver/kmi.yaml",
              "--authorization-kubeconfig=/etc/kubernetes/pki/apiserver/cluster-extension-delegate.kubeconfig",
              "--authentication-kubeconfig=/etc/kubernetes/pki/apiserver/cluster-extension-delegate.kubeconfig",
              "--unified-identity-enabled=true",
              "--unified-identity-cert-file=/var/run/secrets/kubernetes.io/serviceaccount/app.crt",
              "--unified-identity-key-file=/var/run/secrets/kubernetes.io/serviceaccount/app.key",
              "--authorization-always-allow-paths=\"/metrics\"",
              "--audit-log-path=/home/admin/logs/audit.log",
              "--audit-policy-file=/etc/kubernetes/pki/apiserver/audit.yaml",
              "--audit-log-maxage=3",
              "--audit-log-maxsize=10240",
              "--log_file=/home/admin/logs/extension.log",
              "--log_file_max_size=200",
              "--logtostderr=false",
              "--alsologtostderr",
              "--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt",
              "--etcd-certfile=/etc/kubernetes/pki/etcd/client.crt",
              "--etcd-keyfile=/etc/kubernetes/pki/etcd/client.key",
              "--etcd-servers=https://etcd1.ocmpass-eu95.alipay.net:7379,https://etcd2.ocmpass-eu95.alipay.net:7379,https://etcd3.ocmpass-eu95.alipay.net:7379"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ULOGFS_ZCLEAN_ENABLE",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "SN",
                "value": "e9a85417-88b2-41cb-800b-ea015b6dd5da"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "4"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "4"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "4"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              },
              {
                "name": "PARENT_SPEC_GENERATION",
                "value": "47"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver:v0.3.5",
            "imagePullPolicy": "IfNotPresent",
            "name": "apiserver",
            "resources": {
              "limits": {
                "cpu": "4",
                "ephemeral-storage": "100Gi",
                "memory": "8Gi"
              },
              "requests": {
                "cpu": "4",
                "ephemeral-storage": "100Gi",
                "memory": "8Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/pki/etcd/",
                "name": "etcd-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            },
            {
              "name": "single-request-reopen"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "cwtc358700g",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "etcd-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "etcd-credentials"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver:v0.3.1",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "107374182400",
                "ulogfs.k8s.alipay.com/enable-zclean": "true",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "37ca8805-0640-4b5b-bde6-c920dfcf42a1"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "4Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:37:00Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:38Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:58Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-02-15T12:22:00Z",
            "lastTransitionTime": "2023-02-15T12:22:00Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:38Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-03-06T12:37:02Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:57Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:38Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:21:38Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "docker://34ef2941ae87a00c7e7b80cec6a05944d4bf5b2bd7049f271d4e59c04245413c",
            "image": "reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver:v0.3.5",
            "imageID": "docker-pullable://reg.docker.alibaba-inc.com/ocmpaas/cluster-extension-apiserver@sha256:82f14d0ac1a325654256e7e82ee304a2270154bea62bbb34b2b02ae467ddce80",
            "lastState": {},
            "name": "apiserver",
            "ready": true,
            "restartCount": 3,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-03-06T12:36:58Z"
              }
            }
          }
        ],
        "hostIP": "11.166.79.130",
        "phase": "Running",
        "podIP": "11.166.253.69",
        "podIPs": [
          {
            "ip": "11.166.253.69"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-02-15T12:21:38Z"
      }
    }
  },
  {
    "group": "",
    "version": "v1",
    "kind": "Pod",
    "resource": "secrets",
    "resource_version": "",
    "name": "hub-controller-manager-fmvm2-9zhdw",
    "namespace": "ocmpaas",
    "object": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "annotations": {
          "alibabacloud.com/actual-pod-cgroup-path": "/sigma/pod1f77fa71-7218-4a4b-87aa-64f33e2e3db5",
          "cafe.sofastack.io/rulesets": "{\"ocmpaas-cz30a-test\":{\"passPreCheck\":false,\"passPostCheck\":false,\"ruleName\":\"\",\"state\":\"\",\"timestamp\":\"0001-01-01T00:00:00Z\"}}",
          "custom.k8s.alipay.com/original-resource": "{\"containers\":[{\"name\":\"controller-manager\",\"Resources\":{\"limits\":{\"cpu\":\"1\",\"ephemeral-storage\":\"20Gi\",\"memory\":\"2Gi\"},\"requests\":{\"cpu\":\"1\",\"ephemeral-storage\":\"20Gi\",\"memory\":\"2Gi\"}}}]}",
          "meta.k8s.alipay.com/last-spec-hash": "6d48a385bf75c37995c55f78d88263bf",
          "meta.k8s.alipay.com/pod-zappinfo": "{\"spec\":{\"appName\":\"ocmpaas\",\"zone\":\"CZ30A\",\"serverType\":\"DOCKER\",\"fqdn\":\"ocmpaas-cz30a-100088115163.eu95.alipay.net\",\"expectStatus\":\"\"},\"status\":{\"registered\":true,\"message\":\"\",\"status\":\"online\"}}",
          "meta.k8s.alipay.com/trace-context": "[{\"trace_id\":\"5f2e3f73de95c2b80000000000000000\",\"parent_id\":\"\",\"root_span_id\":\"1d0d9e5a3980a51c\",\"delivery_type\":\"PodCreate\",\"status\":\"closed\",\"services\":[{\"component\":\"cloud-scheduler\",\"span_id\":\"c216c9407ffb0b2b\"},{\"component\":\"default-scheduler\",\"span_id\":\"53e1cb6509016ae9\"},{\"component\":\"cni-service\",\"span_id\":\"7e9ff35e97f1525d\"},{\"component\":\"kubelet\",\"span_id\":\"cf7283b47a727310\"},{\"component\":\"zappinfo-controller\",\"span_id\":\"38bf790faf8a0e22\"},{\"component\":\"naming-controller\",\"span_id\":\"1c4fbbfcdbe56451\"}],\"start_at\":\"2023-02-15T20:54:36+08:00\",\"finish_at\":\"2023-02-15T20:55:34+08:00\",\"extra_info\":null}]",
          "orca.identity.alipay.com/serviceaccount": "true",
          "pod.beta1.sigma.ali/alloc-spec": "{\"containers\":[{\"name\":\"controller-manager\",\"resource\":{\"cpu\":{},\"gpu\":{\"shareMode\":\"exclusive\"}},\"hostConfig\":{\"cgroupParent\":\"/sigma\",\"diskQuotaMode\":\"\",\"memorySwap\":2147483648,\"pidsLimit\":32767,\"cpuBvtWarpNs\":2,\"memoryWmarkRatio\":95,\"cpuShares\":1024,\"oomScoreAdj\":-1}}]}",
          "pod.beta1.sigma.ali/hostname-template": "ocmpaas-cz30a-{{.IpAddress}}",
          "pod.beta1.sigma.ali/net-priority": "5",
          "pod.beta1.sigma.ali/network-status": "{\"ipam\":\"ais-ipam\",\"vlan\":\"701\",\"networkPrefixLen\":24,\"gateway\":\"100.88.115.247\",\"netType\":\"vlan\",\"sandboxId\":\"\",\"ip\":\"100.88.115.163\",\"securityDomain\":\"ALIPAY_TEST\"}",
          "pod.beta1.sigma.ali/pod-spec-hash": "hub-controller-manager-fmvm2-fd47d46b4",
          "pod.beta1.sigma.ali/scheduler-update-time": "2023-02-15T20:54:36.9643415+08:00",
          "pod.beta1.sigma.ali/trace-id": "1874f109-f692-4a92-83d1-8a1c719fc04d",
          "pod.beta1.sigma.ali/trace-naming": "{\"id\":\"1874f109-f692-4a92-83d1-8a1c719fc04d\",\"service\":\"naming\",\"creationTimestamp\":\"2023-02-15T20:55:33.380714906+08:00\",\"completionTimestamp\":\"2023-02-15T20:55:33.447071536+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:55:33.380714906+08:00\",\"endTimestamp\":\"2023-02-15T20:55:33.447071536+08:00\",\"logs\":[{\"time\":\"2023-02-15T20:55:33.380720091+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:55:33.380723743+08:00\",\"endTimestamp\":\"2023-02-15T20:55:33.447071139+08:00\"}]}}",
          "pod.beta1.sigma.ali/trace-podfqdn": "{\"TraceID\":\"1874f109-f692-4a92-83d1-8a1c719fc04d\",\"Service\":\"podfqdn\",\"Operation\":\"AddResourceRecord\",\"Error\":false,\"Message\":\"\",\"StartTimestamp\":\"2023-02-15T20:55:33.058063058+08:00\",\"FinishTimestamp\":\"2023-02-15T20:55:33.150272808+08:00\",\"Logs\":{\"error\":\"\"}}",
          "pod.beta1.sigma.ali/trace-zappinfo": "{\"id\":\"1874f109-f692-4a92-83d1-8a1c719fc04d\",\"service\":\"zappinfo\",\"creationTimestamp\":\"2023-02-15T20:55:36.823157036+08:00\",\"completionTimestamp\":\"2023-02-15T20:55:36.823363035+08:00\",\"executionCount\":1,\"span\":{\"operation\":\"ReconcilerHelper.Reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:55:36.823157036+08:00\",\"endTimestamp\":\"2023-02-15T20:55:36.823363035+08:00\",\"logs\":[{\"time\":\"2023-02-15T20:55:36.823160287+08:00\",\"fields\":[{\"key\":\"action\",\"value\":\"update\"}]}],\"children\":[{\"operation\":\"reconcile\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:55:36.823164056+08:00\",\"endTimestamp\":\"2023-02-15T20:55:36.823362904+08:00\",\"children\":[{\"operation\":\"zappinfo.update\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:55:36.823207826+08:00\",\"endTimestamp\":\"2023-02-15T20:55:36.823348917+08:00\",\"children\":[{\"operation\":\"getPodZappinfoMetaSpec\",\"success\":true,\"startTimestamp\":\"2023-02-15T20:55:36.823209193+08:00\",\"endTimestamp\":\"2023-02-15T20:55:36.823348759+08:00\"}]}]}]}}",
          "pod.beta1.sigma.ali/update-status": "{\"statuses\":{\"controller-manager\":{\"creationTimestamp\":\"2023-02-15T20:54:39.055910513+08:00\",\"finishTimestamp\":\"2023-02-15T20:55:31.83004626+08:00\",\"retryCount\":0,\"currentState\":\"running\",\"lastState\":\"unknown\",\"action\":\"start\",\"success\":true,\"message\":\"create start and post start success\",\"specHash\":\"hub-controller-manager-fmvm2-fd47d46b4\"}}}",
          "pod.k8s.alipay.com/auto-eviction": "true",
          "pod.k8s.alipay.com/fqdn-registered-timestamp": "2023-02-15 20:55:33.150320726 +0800 CST m=+1298667.971587914",
          "sigma.ali/container-diskQuotaID": "{\"controller-manager\":\"17794788\"}",
          "trace.cafe.sofastack.io/distribution-info": "{\"Stage\":\"1\",\"Id\":\"RELEASE_202302150851427168715\"}",
          "ulogfs.k8s.alipay.com/biz-disk-quota-repaired": "true",
          "ulogfs.k8s.alipay.com/inject": "enabled"
        },
        "creationTimestamp": "2023-02-15T12:54:36Z",
        "finalizers": [
          "finalizer.k8s.alipay.com/zappinfo",
          "protection-delete.pod.sigma.ali/naming-registered",
          "pod.beta1.sigma.ali/cni-allocated",
          "finalizers.k8s.alipay.com/pod-fqdn"
        ],
        "generateName": "hub-controller-manager-fmvm2-",
        "labels": {
          "alibabacloud.com/quota-name": "ocmpaas-test-sigmaguaranteed-daily",
          "cafe.sofastack.io/app-instance-group": "",
          "cafe.sofastack.io/app-instance-group-name": "",
          "cafe.sofastack.io/cell": "CZ30A",
          "cafe.sofastack.io/control": "true",
          "cafe.sofastack.io/creator": "huanyu",
          "cafe.sofastack.io/deploy-type": "workload",
          "cafe.sofastack.io/global-tenant": "MAIN_SITE",
          "cafe.sofastack.io/pod-ip": "100.88.115.163",
          "cafe.sofastack.io/pod-number": "1",
          "cafe.sofastack.io/pre-check": "false",
          "cafe.sofastack.io/service-available": "1676465737004028798",
          "cafe.sofastack.io/version": "hub-controller-manager-fmvm2-fd47d46b4",
          "cluster.x-k8s.io/cluster-name": "eu95",
          "component": "hub-controller-manager",
          "controller-revision-hash": "hub-controller-manager-fmvm2-fd47d46b4",
          "meta.k8s.alipay.com/app-env": "TEST",
          "meta.k8s.alipay.com/biz-group": "ocmpaas",
          "meta.k8s.alipay.com/biz-group-id": "hub-controller-manager-fmvm2-94cb7efb-f3b1-4715-b86a-808c1d052bd5",
          "meta.k8s.alipay.com/biz-name": "cloudprovision",
          "meta.k8s.alipay.com/delivery-workload": "paascore-cafeext",
          "meta.k8s.alipay.com/fqdn": "ocmpaas-cz30a-100088115163.eu95.alipay.net",
          "meta.k8s.alipay.com/hostname": "ocmpaas-cz30a-100088115163",
          "meta.k8s.alipay.com/migration-level": "L2",
          "meta.k8s.alipay.com/min-replicas": "1",
          "meta.k8s.alipay.com/original-pod-namespace": "ocmpaas",
          "meta.k8s.alipay.com/priority": "production",
          "meta.k8s.alipay.com/qoc-class": "ProdGeneral",
          "meta.k8s.alipay.com/qos-class": "Prod",
          "meta.k8s.alipay.com/replicas": "1",
          "meta.k8s.alipay.com/schedule-time-limit": "30s",
          "meta.k8s.alipay.com/situation": "normal",
          "meta.k8s.alipay.com/slo-resource": "1C2G",
          "meta.k8s.alipay.com/slo-scale": "10",
          "meta.k8s.alipay.com/zone": "CZ30A",
          "paascore.alipay.com/adopted": "1676465676637386740",
          "sigma.ali/app-name": "ocmpaas",
          "sigma.ali/deploy-unit": "ocmpaas-test",
          "sigma.ali/force-update-quota-name": "20230308-153817",
          "sigma.ali/instance-group": "ocmpaassqa",
          "sigma.ali/ip": "100.88.115.163",
          "sigma.ali/qos": "SigmaBurstable",
          "sigma.ali/site": "eu95",
          "sigma.ali/sn": "9fd439fe-d0c3-4532-a147-ff53251bbc20",
          "strategy.cafe.sofastack.io/batch-index": "RELEASE_202302150851427168715-1"
        },
        "name": "hub-controller-manager-fmvm2-9zhdw",
        "namespace": "ocmpaas",
        "ownerReferences": [
          {
            "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
            "blockOwnerDeletion": true,
            "controller": true,
            "kind": "InPlaceSet",
            "name": "hub-controller-manager-fmvm2",
            "uid": "94cb7efb-f3b1-4715-b86a-808c1d052bd5"
          }
        ],
        "resourceVersion": "32427142432",
        "selfLink": "/api/v1/namespaces/ocmpaas/pods/hub-controller-manager-fmvm2-9zhdw",
        "uid": "1f77fa71-7218-4a4b-87aa-64f33e2e3db5"
      },
      "spec": {
        "affinity": {
          "nodeAffinity": {
            "requiredDuringSchedulingIgnoredDuringExecution": {
              "nodeSelectorTerms": [
                {
                  "matchExpressions": [
                    {
                      "key": "sigma.ali/is-over-quota",
                      "operator": "In",
                      "values": [
                        "true"
                      ]
                    }
                  ]
                }
              ]
            }
          }
        },
        "automountServiceAccountToken": true,
        "containers": [
          {
            "command": [
              "kube-controller-manager",
              "--kubeconfig=/etc/kubernetes/config/hub-controller-manager.kubeconfig",
              "--authentication-kubeconfig=/etc/kubernetes/config/hub-controller-manager.kubeconfig",
              "--authorization-kubeconfig=/etc/kubernetes/config/hub-controller-manager.kubeconfig",
              "--requestheader-client-ca-file=/etc/kubernetes/pki/apiserver/ca.crt",
              "--bind-address=127.0.0.1",
              "--cluster-name=hub-cluster",
              "--controllers=*",
              "--leader-elect=true",
              "--port=0",
              "--service-account-private-key-file=/etc/kubernetes/pki/apiserver/sa.key",
              "--use-service-account-credentials=true",
              "--cluster-signing-cert-file=/etc/kubernetes/config/ca.crt",
              "--cluster-signing-key-file=/etc/kubernetes/config/ca.key"
            ],
            "env": [
              {
                "name": "ULOGFS_ENABLED",
                "value": "true"
              },
              {
                "name": "ILOGTAIL_PODNAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ILOGTAIL_ENV",
                "value": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}"
              },
              {
                "name": "container",
                "value": "placeholder"
              },
              {
                "name": "ALIPAY_POD_NAME",
                "valueFrom": {
                  "fieldRef": {
                    "apiVersion": "v1",
                    "fieldPath": "metadata.name"
                  }
                }
              },
              {
                "name": "ALIPAY_POD_NAMESPACE",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ENV",
                "value": "TEST"
              },
              {
                "name": "ALIPAY_APP_APPNAME",
                "value": "ocmpaas"
              },
              {
                "name": "ALIPAY_APP_ZONE",
                "value": "CZ30A"
              },
              {
                "name": "SN",
                "value": "9fd439fe-d0c3-4532-a147-ff53251bbc20"
              },
              {
                "name": "KUBERNETES_SERVICE_HOST",
                "value": "apiserver.sigma-eu95.svc.alipay.net"
              },
              {
                "name": "KUBERNETES_SERVICE_PORT",
                "value": "6443"
              },
              {
                "name": "ALIPAY_SIGMA_CPUMODE",
                "value": "cpushare"
              },
              {
                "name": "SIGMA_MAX_PROCESSORS_LIMIT",
                "value": "1"
              },
              {
                "name": "AJDK_MAX_PROCESSORS_LIMIT",
                "value": "1"
              },
              {
                "name": "LEGACY_CONTAINER_SIZE_CPU_COUNT",
                "value": "1"
              },
              {
                "name": "ali_run_mode",
                "value": "alipay_container"
              }
            ],
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imagePullPolicy": "IfNotPresent",
            "name": "controller-manager",
            "resources": {
              "limits": {
                "cpu": "1",
                "ephemeral-storage": "20Gi",
                "memory": "2Gi"
              },
              "requests": {
                "cpu": "1",
                "ephemeral-storage": "20Gi",
                "memory": "2Gi"
              }
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "volumeMounts": [
              {
                "mountPath": "/etc/kubernetes/config/",
                "name": "hub-controller-manager-kubeconfig",
                "readOnly": true
              },
              {
                "mountPath": "/etc/kubernetes/pki/apiserver/",
                "name": "hub-apiserver-credentials",
                "readOnly": true
              },
              {
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                "name": "default-token-77wms",
                "readOnly": true
              },
              {
                "mountPath": "/home/admin/logs",
                "name": "ulogfs-volume"
              },
              {
                "mountPath": "/dev/shm",
                "name": "shm"
              },
              {
                "mountPath": "/lib/libsysconf-alipay.so",
                "name": "cpushare-volume",
                "readOnly": true
              },
              {
                "mountPath": "/etc/route.tmpl",
                "name": "router-volume",
                "readOnly": true
              }
            ]
          }
        ],
        "dnsConfig": {
          "options": [
            {
              "name": "attempts",
              "value": "2"
            },
            {
              "name": "timeout",
              "value": "2"
            },
            {
              "name": "single-request-reopen"
            }
          ],
          "searches": [
            "ocmpaas.svc.eu95.alipay.net",
            "svc.eu95.alipay.net",
            "eu95.alipay.net"
          ]
        },
        "dnsPolicy": "ClusterFirst",
        "enableServiceLinks": true,
        "imagePullSecrets": [
          {
            "name": "sigma-regcred"
          }
        ],
        "nodeName": "217339291",
        "priority": 0,
        "readinessGates": [
          {
            "conditionType": "cafe.sofastack.io/service-ready"
          },
          {
            "conditionType": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "conditionType": "NamingRegistered"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "ocmpaas-test",
        "serviceAccountName": "ocmpaas-test",
        "terminationGracePeriodSeconds": 30,
        "tolerations": [
          {
            "effect": "NoSchedule",
            "key": "sigma.ali/is-over-quota",
            "operator": "Equal",
            "value": "true"
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/not-ready",
            "operator": "Exists",
            "tolerationSeconds": 300
          },
          {
            "effect": "NoExecute",
            "key": "node.kubernetes.io/unreachable",
            "operator": "Exists",
            "tolerationSeconds": 300
          }
        ],
        "volumes": [
          {
            "name": "hub-controller-manager-kubeconfig",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-controller-manager-kubeconfig"
            }
          },
          {
            "name": "hub-apiserver-credentials",
            "secret": {
              "defaultMode": 420,
              "secretName": "hub-apiserver-credentials"
            }
          },
          {
            "name": "default-token-77wms",
            "secret": {
              "defaultMode": 420,
              "items": [
                {
                  "key": "namespace",
                  "path": "namespace"
                },
                {
                  "key": "token",
                  "path": "token"
                },
                {
                  "key": "ca.crt",
                  "path": "ca.crt"
                },
                {
                  "key": "sa-ca.crt",
                  "path": "sa-ca.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "app.crt"
                },
                {
                  "key": "tls.crt",
                  "path": "tls.crt"
                },
                {
                  "key": "tls.key",
                  "path": "app.key"
                },
                {
                  "key": "tls.key",
                  "path": "tls.key"
                }
              ],
              "secretName": "ocmpaas-test-token-vf94z"
            }
          },
          {
            "csi": {
              "driver": "ulogfs.csi.alipay.com",
              "volumeAttributes": {
                "app.container/image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
                "sigma.ali/app-name": "ocmpaas",
                "sigma.ali/qos": "",
                "sigma.ali/site": "eu95",
                "ulogfs.k8s.alipay.com/disk-quota": "21474836480",
                "ulogfs.k8s.alipay.com/enable-zclean": "",
                "ulogfs.k8s.alipay.com/high-priority": "",
                "ulogfs.k8s.alipay.com/ilogtail": "{\"Appname\":\"ocmpaas\",\"LogAppname\":\"\",\"Idcname\":\"eu95\",\"Apppath\":\"/home/admin/logs\",\"Taglist\":{}}",
                "ulogfs.k8s.alipay.com/lite": "false",
                "ulogfs.k8s.alipay.com/low-priority": "",
                "ulogfs.k8s.alipay.com/ulogfs-preferred-protocol": "fuse",
                "ulogfs.k8s.alipay.com/ulogfs-volume-type": "ulogfs",
                "ulogfs.k8s.alipay.com/volumeid": "f59decd2-f245-4183-82a5-33a1f68bca76"
              }
            },
            "name": "ulogfs-volume"
          },
          {
            "emptyDir": {
              "medium": "Memory",
              "sizeLimit": "1Gi"
            },
            "name": "shm"
          },
          {
            "hostPath": {
              "path": "/lib/libsysconf-alipay.so",
              "type": "File"
            },
            "name": "cpushare-volume"
          },
          {
            "hostPath": {
              "path": "/opt/ali-iaas/env_create/alipay_route.tmpl",
              "type": "File"
            },
            "name": "router-volume"
          }
        ]
      },
      "status": {
        "conditions": [
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:54:36Z",
            "status": "True",
            "type": "cafe.sofastack.io/service-ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:54:37Z",
            "status": "True",
            "type": "IPAllocated"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:55:33Z",
            "reason": "NamingRegisterSucceeded",
            "status": "True",
            "type": "NamingRegistered"
          },
          {
            "lastProbeTime": "2023-02-15T12:55:34Z",
            "lastTransitionTime": "2023-02-15T12:55:34Z",
            "status": "True",
            "type": "readinessgate.k8s.alipay.com/zappinfo"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:54:37Z",
            "status": "True",
            "type": "Initialized"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:55:35Z",
            "status": "True",
            "type": "Ready"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:55:32Z",
            "status": "True",
            "type": "ContainersReady"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:54:37Z",
            "status": "False",
            "type": "ContainerDiskPressure"
          },
          {
            "lastProbeTime": null,
            "lastTransitionTime": "2023-02-15T12:54:37Z",
            "status": "True",
            "type": "PodScheduled"
          }
        ],
        "containerStatuses": [
          {
            "containerID": "pouch://a9964dbbad595256ee7091857528f41aa7fa56bd2ab887e5ed3a71041fb7d8de",
            "image": "reg.docker.alibaba-inc.com/ant-iac/kubernetes:v1.20.1-ee539729",
            "imageID": "reg.docker.alibaba-inc.com/ant-iac/kubernetes@sha256:ce26222470ff2b885084e153400fdb092dfa2764b0a7f94fedc06164d2c2db5d",
            "lastState": {},
            "name": "controller-manager",
            "ready": true,
            "restartCount": 0,
            "started": true,
            "state": {
              "running": {
                "startedAt": "2023-02-15T12:55:31Z"
              }
            }
          }
        ],
        "hostIP": "100.83.13.174",
        "phase": "Running",
        "podIP": "100.88.115.163",
        "podIPs": [
          {
            "ip": "100.88.115.163"
          }
        ],
        "qosClass": "Guaranteed",
        "startTime": "2023-02-15T12:54:37Z"
      }
    }
  }
]

export const clusterList = {
  "apiVersion": "v1",
  "items": [
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster1",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster1.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    },
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster2",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster2.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    },
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster3",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster3.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    },
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster4",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster4.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    },
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster5",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster5.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    },
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster6",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster6.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    },
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster7",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster7.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    },
    {
      "apiVersion": "cluster.alipay-addon.open-cluster-management.io/v1",
      "kind": "ClusterExtension",
      "metadata": {
        "creationTimestamp": "2022-03-18T03:57:01Z",
        "name": "cluster8",
        "resourceVersion": "39470868",
        "uid": "e11ff67f-3242-41df-89a0-bd78d1ed0aee"
      },
      "spec": {
        "access": {
          "caBundle": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHakNDQWdLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREErTVNjd0VnWURWUVFLRXd0QmJHbGkKWVdKaExrbHVZekFSQmdOVkJBb1RDa0ZzYVhCaGVTNUpibU14RXpBUkJnTlZCQU1UQ2t0MVltVnlibVYwWlhNdwpIaGNOTWpJd01qRTFNVEl4T1RNeldoY05Nekl3TWpFek1USXhPVE16V2pBK01TY3dFZ1lEVlFRS0V3dEJiR2xpCllXSmhMa2x1WXpBUkJnTlZCQW9UQ2tGc2FYQmhlUzVKYm1NeEV6QVJCZ05WQkFNVENrdDFZbVZ5Ym1WMFpYTXcKZ2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRREhXWG13cXVka0lsN3VGYVJKVkoybgphNzdBOGwwN0tPYmJkZ21QMkdjREFkZFNlOWZCSk9RWndsS05sNGVLWnp6VERpVVR4MlJXajJVZWFKWUg3VFEwCnlXU2pxZTRMZ3hXVzI0ZHp4YVlkM0VhcVZkUkQ0akxTWU5PNWRvNEh6dUxQOVFBYktyQmtZVEVoM0ZSTXZMOVYKMGI5eTZsWjAxalJKMnNKOEhUcTNzQ1duRlBUc1Uyb0Z1YTYrK2dWYk1Kb0Ywdk1CV05qY01GNHluMVlvQWVOSQoyUXo1QUEybDRHSmhWanVGZkpBKzBZVDVIcDU1TnIzYWVHK1hzSmhoOFUvL0F2a1JYSUdvTjFPSnN5ckpDYkxnCjU1QlJ5UllQMWwxZzFHc2ZnMFNld00xeEcyZ0JIamRqQ25WL0ZsSlczeXpPOC9haXhBN2lFMUppUlFVTFBMd3gKQWdNQkFBR2pJekFoTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxRwpTSWIzRFFFQkN3VUFBNElCQVFCbGdveC9MajFKWGZBUHdnMEMwNUJFcGNBa1o4THhpUlN2SUhpdENFWGFOMWxZCk1URHVyV29NTVdrWDhqdDNWb2hGL1dQQ1FwVHNmQnd5UnNNQSszUDJyamhwaTUvTVpISjhtRUdwZ2Q4YTRwK3cKOFVLaC9OZG9KKzQ3UEtXNlJFdDFwbFdMWWJrQm84V0VVMmQrOXRseXNDTCsrZnR4NTc4eXYraWhDTVhRV2hoagpycHZraGFITFZ4U3FYZ2hKdTkrTzMraHB1eE9SaGZaMVZVZU5ydXo1VUYyd3drU0Jkc1NKa0pObGhUZjN1Z2NVCndQcDhZOTZncFBOaEFTWE9uNWNXVHVSR0NybGRqT2VRUUNFN1NHSEpUM2RKRldPSHdjV3hTSGVKQU9IQkhndU4KK0RNMVV2aDNWbk0zdElIUEh3ZFpOdUorZ3ZYMUFEcGZYa0JiVFlQOQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
          "credential": {
            "type": "UnifiedIdentity"
          },
          "endpoint": "https://apiserver.cluster8.svc.alipay.net:6443"
        },
        "provider": "SigmaBoss"
      },
      "status": { "healthy": "true", "delay": "15ms", "node": 99 }
    }
  ],
  "kind": "List",
  "metadata": {
    "resourceVersion": ""
  }
}

export const searchData = {
  "kind": "UniResource",
  "apiVersion": "search.karbour.com/v1beta1",
  "objects": [
    {
      "apiVersion": "v1",
      "kind": "Namespace",
      "metadata": {
        "name": "test1",
        "namespace": ""
      }
    },
    {
      "apiVersion": "v1",
      "kind": "Namespace",
      "metadata": {
        "name": "test2",
        "namespace": ""
      }
    },
    {
      "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
      "kind": "CafeDeployment",
      "metadata": {
        "name": "test1",
        "namespace": "test1"
      }
    },
    {
      "apiVersion": "apps.cafe.cloud.alipay.com/v1alpha1",
      "kind": "CafeDeployment",
      "metadata": {
        "name": "test2",
        "namespace": "test2"
      }
    },
    {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "name": "test1",
        "namespace": "test1"
      }
    },
    {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "name": "test2",
        "namespace": "test2"
      }
    }
  ]
}
