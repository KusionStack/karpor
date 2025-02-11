apiVersion: v1
kind: Namespace
metadata:
  name: karpor
spec:
  finalizers:
  - kubernetes
{{- if eq .Mode "pull" }}
---
apiVersion: v1
data:
  config: |-
    apiVersion: v1
    clusters:
        - cluster:
              insecure-skip-tls-verify: true
              server: {{ .ExternalEndpoint }}
          name: karpor
    contexts:
        - context:
              cluster: karpor
              user: {{ .ClusterName }}
          name: default
    current-context: default
    kind: Config
    users:
        - name: {{ .ClusterName }}
          user:
              client-certificate-data: {{ .CaCert }}
              client-key-data: {{ .CaKey }}
kind: ConfigMap
metadata:
  name: karpor-kubeconfig
  namespace: karpor
{{- end }}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: karpor-agent
  namespace: karpor
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app.kubernetes.io/component: karpor-agent
      app.kubernetes.io/instance: karpor
      app.kubernetes.io/name: karpor
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/component: karpor-agent
        app.kubernetes.io/instance: karpor
        app.kubernetes.io/name: karpor
    spec:
      containers:
      - args:
        - agent
        - --elastic-search-addresses={{ range .StorageAddresses }}{{.}} {{ end }}
        - --cluster-name={{ .ClusterName }}
        - --cluster-mode={{ .ClusterMode }}
        command:
        - /karpor
{{- if eq .Mode "pull" }}
        env:
        - name: KUBECONFIG
          value: /etc/karpor/config
{{- end }}
        image: kusionstack/karpor:v0.5.9
        imagePullPolicy: IfNotPresent
        name: karpor-agent
        ports:
        - containerPort: 7443
          protocol: TCP
        resources:
{{- if eq .Level 3 }}
          limits:
            cpu: 1
            ephemeral-storage: 20Gi
            memory: 2Gi
          requests:
            cpu: 500m
            ephemeral-storage: 4Gi
            memory: 512Mi
{{- else if eq .Level 2 }}
          limits:
            cpu: 500m
            ephemeral-storage: 10Gi
            memory: 1Gi
          requests:
            cpu: 250m
            ephemeral-storage: 2Gi
            memory: 256Mi
{{- else }}
          limits:
            cpu: 250m
            ephemeral-storage: 5Gi
            memory: 500Mi
          requests:
            cpu: 125m
            ephemeral-storage: 1Gi
            memory: 128Mi
{{- end }}
{{- if eq .Mode "pull" }}
        volumeMounts:
        - mountPath: /etc/karpor/
          name: karpor-kubeconfig
{{- end }}
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
{{- if eq .Mode "pull" }}
      volumes:
      - configMap:
          defaultMode: 420
          name: karpor-kubeconfig
        name: karpor-kubeconfig
{{- end }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: karpor
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: default
  namespace: karpor
