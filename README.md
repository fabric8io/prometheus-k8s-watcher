# Prometheus Kubernetes Watcher

This utility watches Kubernetes & generates the appropriate
config files for Prometheus. It can watch Kubernetes resources, e.g.
nodes, & update config files on the fly.

## Usage
```
Usage of prometheus-k8s-watcher:
  -alsologtostderr=false: log to standard error as well as files
  -bearer-token-file="/var/run/secrets/kubernetes.io/serviceaccount/token": The file containing the bearer token.
  -insecure=false: If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure.
  -log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
  -log_dir="": If non-empty, write log files in this directory
  -logtostderr=false: log to standard error instead of files
  -master="https://kubernetes.default.svc.cluster.local": The URL of the Kubernetes API server
  -node-read-only-port=10255: The port that metrics can be retrieved from the nodes.
  -nodes-file="/etc/prometheus/config.d/nodes.yml": The file to write the node targets to.
  -stderrthreshold=0: logs at or above this threshold go to stderr
  -v=0: log level for V logs
  -vmodule=: comma-separated list of pattern=N settings for file-filtered logging
```
