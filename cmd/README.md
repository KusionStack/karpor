# `/cmd`

Main applications for `karbour`.

It's common to have a small `main` function that imports and invokes the code from the `/pkg` directories and nothing else.

Examples:

* https://github.com/vmware-tanzu/velero/tree/main/cmd (just a really small `main` function with everything else in packages)
* https://github.com/moby/moby/tree/master/cmd
* https://github.com/prometheus/prometheus/tree/main/cmd
* https://github.com/influxdata/influxdb/tree/master/cmd
* https://github.com/kubernetes/kubernetes/tree/master/cmd
* https://github.com/dapr/dapr/tree/master/cmd
* https://github.com/ethereum/go-ethereum/tree/master/cmd
