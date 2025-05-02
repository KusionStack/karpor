## karpor

Launch an API server

### Synopsis

Launch an API server

```
karpor [flags]
```

### Options

```
      --admission-control-config-file string                    File with admission control configuration.
      --advertise-address ip                                    The IP address on which to advertise the apiserver to members of the cluster. This address must be reachable by the rest of the cluster. If blank, the --bind-address will be used. If --bind-address is unspecified, the host's default interface will be used.
      --ai-auth-token string                                    The ai auth token
      --ai-backend string                                       The ai backend (default "openai")
      --ai-base-url string                                      The ai base url
      --ai-http-proxy string                                    The ai http proxy
      --ai-https-proxy string                                   The ai https proxy
      --ai-model string                                         The ai model (default "gpt-3.5-turbo")
      --ai-no-proxy string                                      The ai no-proxy
      --ai-proxy-enabled                                        The ai proxy enable
      --ai-temperature float32                                  The ai temperature (default 1)
      --ai-top-p float32                                        The ai top-p (default 1)
      --anonymous-auth                                          Enables anonymous requests to the secure port of the API server. Requests that are not rejected by another authentication method are treated as anonymous requests. Anonymous requests have a username of system:anonymous, and a group name of system:unauthenticated. (default true)
      --api-audiences strings                                   Identifiers of the API. The service account token authenticator will validate that tokens used against the API are bound to at least one of these audiences. If the --service-account-issuer flag is configured and this flag is not, this field defaults to a single element list containing the issuer URL.
      --audit-log-batch-buffer-size int                         The size of the buffer to store events before batching and writing. Only used in batch mode. (default 10000)
      --audit-log-batch-max-size int                            The maximum size of a batch. Only used in batch mode. (default 1)
      --audit-log-batch-max-wait duration                       The amount of time to wait before force writing the batch that hadn't reached the max size. Only used in batch mode.
      --audit-log-batch-throttle-burst int                      Maximum number of requests sent at the same moment if ThrottleQPS was not utilized before. Only used in batch mode.
      --audit-log-batch-throttle-enable                         Whether batching throttling is enabled. Only used in batch mode.
      --audit-log-batch-throttle-qps float32                    Maximum average number of batches per second. Only used in batch mode.
      --audit-log-compress                                      If set, the rotated log files will be compressed using gzip.
      --audit-log-format string                                 Format of saved audits. "legacy" indicates 1-line text format for each event. "json" indicates structured json format. Known formats are legacy,json. (default "json")
      --audit-log-maxage int                                    The maximum number of days to retain old audit log files based on the timestamp encoded in their filename.
      --audit-log-maxbackup int                                 The maximum number of old audit log files to retain. Setting a value of 0 will mean there's no restriction on the number of files.
      --audit-log-maxsize int                                   The maximum size in megabytes of the audit log file before it gets rotated.
      --audit-log-mode string                                   Strategy for sending audit events. Blocking indicates sending events should block server responses. Batch causes the backend to buffer and write events asynchronously. Known modes are batch,blocking,blocking-strict. (default "blocking")
      --audit-log-path string                                   If set, all requests coming to the apiserver will be logged to this file.  '-' means standard out.
      --audit-log-truncate-enabled                              Whether event and batch truncating is enabled.
      --audit-log-truncate-max-batch-size int                   Maximum size of the batch sent to the underlying backend. Actual serialized size can be several hundreds of bytes greater. If a batch exceeds this limit, it is split into several batches of smaller size. (default 10485760)
      --audit-log-truncate-max-event-size int                   Maximum size of the audit event sent to the underlying backend. If the size of an event is greater than this number, first request and response are removed, and if this doesn't reduce the size enough, event is discarded. (default 102400)
      --audit-log-version string                                API group and version used for serializing audit events written to log. (default "audit.k8s.io/v1")
      --audit-policy-file string                                Path to the file that defines the audit policy configuration.
      --audit-webhook-batch-buffer-size int                     The size of the buffer to store events before batching and writing. Only used in batch mode. (default 10000)
      --audit-webhook-batch-max-size int                        The maximum size of a batch. Only used in batch mode. (default 400)
      --audit-webhook-batch-max-wait duration                   The amount of time to wait before force writing the batch that hadn't reached the max size. Only used in batch mode. (default 30s)
      --audit-webhook-batch-throttle-burst int                  Maximum number of requests sent at the same moment if ThrottleQPS was not utilized before. Only used in batch mode. (default 15)
      --audit-webhook-batch-throttle-enable                     Whether batching throttling is enabled. Only used in batch mode. (default true)
      --audit-webhook-batch-throttle-qps float32                Maximum average number of batches per second. Only used in batch mode. (default 10)
      --audit-webhook-config-file string                        Path to a kubeconfig formatted file that defines the audit webhook configuration.
      --audit-webhook-initial-backoff duration                  The amount of time to wait before retrying the first failed request. (default 10s)
      --audit-webhook-mode string                               Strategy for sending audit events. Blocking indicates sending events should block server responses. Batch causes the backend to buffer and write events asynchronously. Known modes are batch,blocking,blocking-strict. (default "batch")
      --audit-webhook-truncate-enabled                          Whether event and batch truncating is enabled.
      --audit-webhook-truncate-max-batch-size int               Maximum size of the batch sent to the underlying backend. Actual serialized size can be several hundreds of bytes greater. If a batch exceeds this limit, it is split into several batches of smaller size. (default 10485760)
      --audit-webhook-truncate-max-event-size int               Maximum size of the audit event sent to the underlying backend. If the size of an event is greater than this number, first request and response are removed, and if this doesn't reduce the size enough, event is discarded. (default 102400)
      --audit-webhook-version string                            API group and version used for serializing audit events written to webhook. (default "audit.k8s.io/v1")
      --authorization-mode strings                              Ordered list of plug-ins to do authorization on secure port. Comma-delimited list of: AlwaysAllow,AlwaysDeny,ABAC,Webhook,RBAC,Node. (default [RBAC])
      --authorization-policy-file string                        File with authorization policy in json line by line format, used with --authorization-mode=ABAC, on the secure port.
      --authorization-webhook-cache-authorized-ttl duration     The duration to cache 'authorized' responses from the webhook authorizer. (default 5m0s)
      --authorization-webhook-cache-unauthorized-ttl duration   The duration to cache 'unauthorized' responses from the webhook authorizer. (default 30s)
      --authorization-webhook-config-file string                File with webhook configuration in kubeconfig format, used with --authorization-mode=Webhook. The API server will query the remote service to determine access on the API server's secure port.
      --authorization-webhook-version string                    The API version of the authorization.k8s.io SubjectAccessReview to send to and expect from the webhook. (default "v1beta1")
      --bind-address ip                                         The IP address on which to listen for the --secure-port port. The associated interface(s) must be reachable by the rest of the cluster, and by CLI/web clients. If blank or an unspecified address (0.0.0.0 or ::), all interfaces will be used. (default 0.0.0.0)
      --cert-dir string                                         The directory where the TLS certs are located. If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored. (default "apiserver.local.config/certificates")
      --client-ca-file string                                   If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate.
      --contention-profiling                                    Enable lock contention profiling, if profiling is enabled
      --cors-allowed-origins strings                            List of allowed origins for CORS, comma separated.  An allowed origin can be a regular expression to support subdomain matching. If this list is empty CORS will not be enabled. (default [.*])
      --delete-collection-workers int                           Number of workers spawned for DeleteCollection call. These are used to speed up namespace cleanup. (default 1)
      --disable-admission-plugins strings                       admission plugins that should be disabled although they are in the default enabled plugins list (NamespaceLifecycle, MutatingAdmissionWebhook, ValidatingAdmissionPolicy, ValidatingAdmissionWebhook). Comma-delimited list of admission plugins: MutatingAdmissionWebhook, NamespaceLifecycle, ValidatingAdmissionPolicy, ValidatingAdmissionWebhook. The order of plugins in this flag does not matter. (default [MutatingAdmissionWebhook,NamespaceLifecycle,ValidatingAdmissionWebhook,ValidatingAdmissionPolicy])
      --egress-selector-config-file string                      File with apiserver egress selector configuration.
      --elastic-search-addresses strings                        The elastic search address
      --elastic-search-password string                          The elastic search password
      --elastic-search-username string                          The elastic search username
      --enable-admission-plugins strings                        admission plugins that should be enabled in addition to default enabled ones (NamespaceLifecycle, MutatingAdmissionWebhook, ValidatingAdmissionPolicy, ValidatingAdmissionWebhook). Comma-delimited list of admission plugins: MutatingAdmissionWebhook, NamespaceLifecycle, ValidatingAdmissionPolicy, ValidatingAdmissionWebhook. The order of plugins in this flag does not matter.
      --enable-garbage-collector                                Enables the generic garbage collector. MUST be synced with the corresponding flag of the kube-controller-manager. (default true)
      --enable-priority-and-fairness                            If true and the APIPriorityAndFairness feature gate is enabled, replace the max-in-flight handler with an enhanced one that queues and dispatches with priority and fairness (default true)
      --enable-rbac                                             trun on to enable RBAC authorization
      --encryption-provider-config string                       The file containing configuration for encryption providers to be used for storing secrets in etcd
      --encryption-provider-config-automatic-reload             Determines if the file set by --encryption-provider-config should be automatically reloaded if the disk contents change. Setting this to true disables the ability to uniquely identify distinct KMS plugins via the API server healthz endpoints.
      --etcd-cafile string                                      SSL Certificate Authority file used to secure etcd communication.
      --etcd-certfile string                                    SSL certification file used to secure etcd communication.
      --etcd-compaction-interval duration                       The interval of compaction requests. If 0, the compaction request from apiserver is disabled. (default 5m0s)
      --etcd-count-metric-poll-period duration                  Frequency of polling etcd for number of resources per type. 0 disables the metric collection. (default 1m0s)
      --etcd-db-metric-poll-interval duration                   The interval of requests to poll etcd and update metric. 0 disables the metric collection (default 30s)
      --etcd-healthcheck-timeout duration                       The timeout to use when checking etcd health. (default 2s)
      --etcd-keyfile string                                     SSL key file used to secure etcd communication.
      --etcd-prefix string                                      The prefix to prepend to all resource paths in etcd. (default "/registry/karpor")
      --etcd-readycheck-timeout duration                        The timeout to use when checking etcd readiness (default 2s)
      --etcd-servers strings                                    List of etcd servers to connect with (scheme://ip:port), comma separated.
      --etcd-servers-overrides strings                          Per-resource etcd servers overrides, comma separated. The individual override format: group/resource#servers, where servers are URLs, semicolon separated. Note that this applies only to resources compiled into this server binary. 
      --external-hostname string                                The hostname to use when generating externalized URLs for this master (e.g. Swagger API Docs or OpenID Discovery).
      --feature-gates mapStringBool                             A set of key=value pairs that describe feature gates for alpha/experimental features. Options are:
                                                                APIListChunking=true|false (BETA - default=true)
                                                                APIPriorityAndFairness=true|false (BETA - default=true)
                                                                APIResponseCompression=true|false (BETA - default=true)
                                                                APISelfSubjectReview=true|false (ALPHA - default=false)
                                                                APIServerIdentity=true|false (BETA - default=true)
                                                                APIServerTracing=true|false (ALPHA - default=false)
                                                                AggregatedDiscoveryEndpoint=true|false (ALPHA - default=false)
                                                                AllAlpha=true|false (ALPHA - default=false)
                                                                AllBeta=true|false (BETA - default=false)
                                                                AnyVolumeDataSource=true|false (BETA - default=true)
                                                                AppArmor=true|false (BETA - default=true)
                                                                CPUManagerPolicyAlphaOptions=true|false (ALPHA - default=false)
                                                                CPUManagerPolicyBetaOptions=true|false (BETA - default=true)
                                                                CPUManagerPolicyOptions=true|false (BETA - default=true)
                                                                CSIMigrationPortworx=true|false (BETA - default=false)
                                                                CSIMigrationRBD=true|false (ALPHA - default=false)
                                                                CSINodeExpandSecret=true|false (ALPHA - default=false)
                                                                CSIVolumeHealth=true|false (ALPHA - default=false)
                                                                ComponentSLIs=true|false (ALPHA - default=false)
                                                                ContainerCheckpoint=true|false (ALPHA - default=false)
                                                                CronJobTimeZone=true|false (BETA - default=true)
                                                                CrossNamespaceVolumeDataSource=true|false (ALPHA - default=false)
                                                                CustomCPUCFSQuotaPeriod=true|false (ALPHA - default=false)
                                                                CustomResourceValidationExpressions=true|false (BETA - default=true)
                                                                DisableCloudProviders=true|false (ALPHA - default=false)
                                                                DisableKubeletCloudCredentialProviders=true|false (ALPHA - default=false)
                                                                DownwardAPIHugePages=true|false (BETA - default=true)
                                                                DynamicResourceAllocation=true|false (ALPHA - default=false)
                                                                EventedPLEG=true|false (ALPHA - default=false)
                                                                ExpandedDNSConfig=true|false (BETA - default=true)
                                                                ExperimentalHostUserNamespaceDefaulting=true|false (BETA - default=false)
                                                                GRPCContainerProbe=true|false (BETA - default=true)
                                                                GracefulNodeShutdown=true|false (BETA - default=true)
                                                                GracefulNodeShutdownBasedOnPodPriority=true|false (BETA - default=true)
                                                                HPAContainerMetrics=true|false (ALPHA - default=false)
                                                                HPAScaleToZero=true|false (ALPHA - default=false)
                                                                HonorPVReclaimPolicy=true|false (ALPHA - default=false)
                                                                IPTablesOwnershipCleanup=true|false (ALPHA - default=false)
                                                                InTreePluginAWSUnregister=true|false (ALPHA - default=false)
                                                                InTreePluginAzureDiskUnregister=true|false (ALPHA - default=false)
                                                                InTreePluginAzureFileUnregister=true|false (ALPHA - default=false)
                                                                InTreePluginGCEUnregister=true|false (ALPHA - default=false)
                                                                InTreePluginOpenStackUnregister=true|false (ALPHA - default=false)
                                                                InTreePluginPortworxUnregister=true|false (ALPHA - default=false)
                                                                InTreePluginRBDUnregister=true|false (ALPHA - default=false)
                                                                InTreePluginvSphereUnregister=true|false (ALPHA - default=false)
                                                                JobMutableNodeSchedulingDirectives=true|false (BETA - default=true)
                                                                JobPodFailurePolicy=true|false (BETA - default=true)
                                                                JobReadyPods=true|false (BETA - default=true)
                                                                KMSv2=true|false (ALPHA - default=false)
                                                                KubeletInUserNamespace=true|false (ALPHA - default=false)
                                                                KubeletPodResources=true|false (BETA - default=true)
                                                                KubeletPodResourcesGetAllocatable=true|false (BETA - default=true)
                                                                KubeletTracing=true|false (ALPHA - default=false)
                                                                LegacyServiceAccountTokenTracking=true|false (ALPHA - default=false)
                                                                LocalStorageCapacityIsolationFSQuotaMonitoring=true|false (ALPHA - default=false)
                                                                LogarithmicScaleDown=true|false (BETA - default=true)
                                                                MatchLabelKeysInPodTopologySpread=true|false (ALPHA - default=false)
                                                                MaxUnavailableStatefulSet=true|false (ALPHA - default=false)
                                                                MemoryManager=true|false (BETA - default=true)
                                                                MemoryQoS=true|false (ALPHA - default=false)
                                                                MinDomainsInPodTopologySpread=true|false (BETA - default=false)
                                                                MinimizeIPTablesRestore=true|false (ALPHA - default=false)
                                                                MultiCIDRRangeAllocator=true|false (ALPHA - default=false)
                                                                NetworkPolicyStatus=true|false (ALPHA - default=false)
                                                                NodeInclusionPolicyInPodTopologySpread=true|false (BETA - default=true)
                                                                NodeOutOfServiceVolumeDetach=true|false (BETA - default=true)
                                                                NodeSwap=true|false (ALPHA - default=false)
                                                                OpenAPIEnums=true|false (BETA - default=true)
                                                                OpenAPIV3=true|false (BETA - default=true)
                                                                PDBUnhealthyPodEvictionPolicy=true|false (ALPHA - default=false)
                                                                PodAndContainerStatsFromCRI=true|false (ALPHA - default=false)
                                                                PodDeletionCost=true|false (BETA - default=true)
                                                                PodDisruptionConditions=true|false (BETA - default=true)
                                                                PodHasNetworkCondition=true|false (ALPHA - default=false)
                                                                PodSchedulingReadiness=true|false (ALPHA - default=false)
                                                                ProbeTerminationGracePeriod=true|false (BETA - default=true)
                                                                ProcMountType=true|false (ALPHA - default=false)
                                                                ProxyTerminatingEndpoints=true|false (BETA - default=true)
                                                                QOSReserved=true|false (ALPHA - default=false)
                                                                ReadWriteOncePod=true|false (ALPHA - default=false)
                                                                RecoverVolumeExpansionFailure=true|false (ALPHA - default=false)
                                                                RemainingItemCount=true|false (BETA - default=true)
                                                                RetroactiveDefaultStorageClass=true|false (BETA - default=true)
                                                                RotateKubeletServerCertificate=true|false (BETA - default=true)
                                                                SELinuxMountReadWriteOncePod=true|false (ALPHA - default=false)
                                                                SeccompDefault=true|false (BETA - default=true)
                                                                ServerSideFieldValidation=true|false (BETA - default=true)
                                                                SizeMemoryBackedVolumes=true|false (BETA - default=true)
                                                                StatefulSetAutoDeletePVC=true|false (ALPHA - default=false)
                                                                StatefulSetStartOrdinal=true|false (ALPHA - default=false)
                                                                StorageVersionAPI=true|false (ALPHA - default=false)
                                                                StorageVersionHash=true|false (BETA - default=true)
                                                                TopologyAwareHints=true|false (BETA - default=true)
                                                                TopologyManager=true|false (BETA - default=true)
                                                                TopologyManagerPolicyAlphaOptions=true|false (ALPHA - default=false)
                                                                TopologyManagerPolicyBetaOptions=true|false (BETA - default=false)
                                                                TopologyManagerPolicyOptions=true|false (ALPHA - default=false)
                                                                UserNamespacesStatelessPodsSupport=true|false (ALPHA - default=false)
                                                                ValidatingAdmissionPolicy=true|false (ALPHA - default=false)
                                                                VolumeCapacityPriority=true|false (ALPHA - default=false)
                                                                WinDSR=true|false (ALPHA - default=false)
                                                                WinOverlay=true|false (BETA - default=true)
                                                                WindowsHostNetwork=true|false (ALPHA - default=true) (default APIPriorityAndFairness=true)
      --github-badge                                            whether to display the github badge
      --goaway-chance float                                     To prevent HTTP/2 clients from getting stuck on a single apiserver, randomly close a connection (GOAWAY). The client's other in-flight requests won't be affected, and the client will reconnect, likely landing on a different apiserver after going through the load balancer again. This argument sets the fraction of requests that will be sent a GOAWAY. Clusters with single apiservers, or which don't use a load balancer, should NOT enable this. Min is 0 (off), Max is .02 (1/50 requests); .001 (1/1000) is a recommended starting point.
  -h, --help                                                    help for karpor
      --high-availability                                       whether to use high-availability feature.
      --http2-max-streams-per-connection int                    The limit that the server gives to clients for the maximum number of streams in an HTTP/2 connection. Zero means to use golang's default. (default 1000)
      --lease-reuse-duration-seconds int                        The time in seconds that each lease is reused. A lower value could avoid large number of objects reusing the same lease. Notice that a too small value may cause performance problems at storage layer. (default 60)
      --livez-grace-period duration                             This option represents the maximum amount of time it should take for apiserver to complete its startup sequence and become live. From apiserver's start time to when this amount of time has elapsed, /livez will assume that unfinished post-start hooks will complete successfully and therefore return true.
      --max-mutating-requests-inflight int                      This and --max-requests-inflight are summed to determine the server's total concurrency limit (which must be positive) if --enable-priority-and-fairness is true. Otherwise, this flag limits the maximum number of mutating requests in flight, or a zero value disables the limit completely. (default 200)
      --max-requests-inflight int                               This and --max-mutating-requests-inflight are summed to determine the server's total concurrency limit (which must be positive) if --enable-priority-and-fairness is true. Otherwise, this flag limits the maximum number of non-mutating requests in flight, or a zero value disables the limit completely. (default 400)
      --min-request-timeout int                                 An optional field indicating the minimum number of seconds a handler must keep a request open before timing it out. Currently only honored by the watch request handler, which picks a randomized value above this number as the connection timeout, to spread out load. (default 1800)
      --permit-address-sharing                                  If true, SO_REUSEADDR will be used when binding the port. This allows binding to wildcard IPs like 0.0.0.0 and specific IPs in parallel, and it avoids waiting for the kernel to release sockets in TIME_WAIT state. [default=false]
      --permit-port-sharing                                     If true, SO_REUSEPORT will be used when binding the port, which allows more than one instance to bind on the same address and port. [default=false]
      --profiling                                               Enable profiling via web interface host:port/debug/pprof/ (default true)
      --read-only-mode                                          turn on the read only mode
      --request-timeout duration                                An optional field indicating the duration a handler must keep a request open before timing it out. This is the default request timeout for requests but may be overridden by flags such as --min-request-timeout for specific types of requests. (default 1m0s)
      --requestheader-allowed-names strings                     List of client certificate common names to allow to provide usernames in headers specified by --requestheader-username-headers. If empty, any client certificate validated by the authorities in --requestheader-client-ca-file is allowed.
      --requestheader-client-ca-file string                     Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers. WARNING: generally do not depend on authorization being already done for incoming requests.
      --requestheader-extra-headers-prefix strings              List of request header prefixes to inspect. X-Remote-Extra- is suggested.
      --requestheader-group-headers strings                     List of request headers to inspect for groups. X-Remote-Group is suggested.
      --requestheader-username-headers strings                  List of request headers to inspect for usernames. X-Remote-User is common.
      --search-storage-type string                              The search storage type
      --secure-port int                                         The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all. (default 443)
      --service-account-extend-token-expiration                 Turns on projected service account expiration extension during token generation, which helps safe transition from legacy token to bound service account token feature. If this flag is enabled, admission injected tokens would be extended up to 1 year to prevent unexpected failure during transition, ignoring value of service-account-max-token-expiration. (default true)
      --service-account-issuer stringArray                      Identifier of the service account token issuer. The issuer will assert this identifier in "iss" claim of issued tokens. This value is a string or URI. If this option is not a valid URI per the OpenID Discovery 1.0 spec, the ServiceAccountIssuerDiscovery feature will remain disabled, even if the feature gate is set to true. It is highly recommended that this value comply with the OpenID spec: https://openid.net/specs/openid-connect-discovery-1_0.html. In practice, this means that service-account-issuer must be an https URL. It is also highly recommended that this URL be capable of serving OpenID discovery documents at {service-account-issuer}/.well-known/openid-configuration. When this flag is specified multiple times, the first is used to generate tokens and all are used to determine which issuers are accepted.
      --service-account-jwks-uri string                         Overrides the URI for the JSON Web Key Set in the discovery doc served at /.well-known/openid-configuration. This flag is useful if the discovery docand key set are served to relying parties from a URL other than the API server's external (as auto-detected or overridden with external-hostname). 
      --service-account-key-file stringArray                    File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify ServiceAccount tokens. The specified file can contain multiple keys, and the flag can be specified multiple times with different files. If unspecified, --tls-private-key-file is used. Must be specified when --service-account-signing-key-file is provided
      --service-account-lookup                                  If true, validate ServiceAccount tokens exist in etcd as part of authentication. (default true)
      --service-account-max-token-expiration duration           The maximum validity duration of a token created by the service account token issuer. If an otherwise valid TokenRequest with a validity duration larger than this value is requested, a token will be issued with a validity duration of this value.
      --service-account-signing-key-file string                 Path to the file that contains the current private key of the service account token issuer. The issuer will sign issued ID tokens with this private key.
      --shutdown-delay-duration duration                        Time to delay the termination. During that time the server keeps serving requests normally. The endpoints /healthz and /livez will return success, but /readyz immediately returns failure. Graceful termination starts after this delay has elapsed. This can be used to allow load balancer to stop sending traffic to this server.
      --shutdown-send-retry-after                               If true the HTTP Server will continue listening until all non long running request(s) in flight have been drained, during this window all incoming requests will be rejected with a status code 429 and a 'Retry-After' response header, in addition 'Connection: close' response header is set in order to tear down the TCP connection when idle.
      --storage-backend string                                  The storage backend for persistence. Options: 'etcd3' (default).
      --storage-media-type string                               The media type to use to store objects in storage. Some resources or storage backends may only support a specific media type and will ignore this setting. Supported media types: [application/json, application/yaml, application/vnd.kubernetes.protobuf] (default "application/json")
      --strict-transport-security-directives strings            List of directives for HSTS, comma separated. If this list is empty, then HSTS directives will not be added. Example: 'max-age=31536000,includeSubDomains,preload'
      --tls-cert-file string                                    File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert). If HTTPS serving is enabled, and --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir. (default "apiserver.local.config/certificates/apiserver.crt")
      --tls-cipher-suites strings                               Comma-separated list of cipher suites for the server. If omitted, the default Go cipher suites will be used. 
                                                                Preferred values: TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256, TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA, TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA, TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256. 
                                                                Insecure values: TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, TLS_ECDHE_ECDSA_WITH_RC4_128_SHA, TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA, TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256, TLS_ECDHE_RSA_WITH_RC4_128_SHA, TLS_RSA_WITH_3DES_EDE_CBC_SHA, TLS_RSA_WITH_AES_128_CBC_SHA, TLS_RSA_WITH_AES_128_CBC_SHA256, TLS_RSA_WITH_AES_128_GCM_SHA256, TLS_RSA_WITH_AES_256_CBC_SHA, TLS_RSA_WITH_AES_256_GCM_SHA384, TLS_RSA_WITH_RC4_128_SHA.
      --tls-min-version string                                  Minimum TLS version supported. Possible values: VersionTLS10, VersionTLS11, VersionTLS12, VersionTLS13
      --tls-private-key-file string                             File containing the default x509 private key matching --tls-cert-file. (default "apiserver.local.config/certificates/apiserver.key")
      --tls-sni-cert-key namedCertKey                           A pair of x509 certificate and private key file paths, optionally suffixed with a list of domain patterns which are fully qualified domain names, possibly with prefixed wildcard segments. The domain patterns also allow IP addresses, but IPs should only be used if the apiserver has visibility to the IP address requested by a client. If no domain patterns are provided, the names of the certificate are extracted. Non-wildcard matches trump over wildcard matches, explicit domain patterns trump over extracted names. For multiple key/certificate pairs, use the --tls-sni-cert-key multiple times. Examples: "example.crt,example.key" or "foo.crt,foo.key:*.foo.com,foo.com". (default [])
      --tracing-config-file string                              File with apiserver tracing configuration.
  -V, --version                                                 Print version and exit
      --watch-cache                                             Enable watch caching in the apiserver (default true)
      --watch-cache-sizes strings                               Watch cache size settings for some resources (pods, nodes, etc.), comma separated. The individual setting format: resource[.group]#size, where resource is lowercase plural (no version), group is omitted for resources of apiVersion v1 (the legacy core API) and included for others, and size is a number. This option is only meaningful for resources built into the apiserver, not ones defined by CRDs or aggregated from external servers, and is only consulted if the watch-cache is enabled. The only meaningful size setting to supply here is zero, which means to disable watch caching for the associated resource; all non-zero values are equivalent and mean to not disable watch caching for that resource
```

### SEE ALSO

* [karpor syncer](karpor_syncer.md)	 - start a resource syncer to sync resource from clusters

###### Auto generated by spf13/cobra on 2-May-2025
