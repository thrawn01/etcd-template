# Introduction
etcd-template is a simple implementation of confd or consul-template which
serves a very specific use case for mailgun. I hope this is a transitional tool
as while we implement the auto pilot pattern for our contanerized micro
services.

# Usage
```
Usage: etcd-template [OPTIONS]  <etcd-path> <template-dir> [output-dir]

Read mailgun compatable etcd dictionaries from etcd and generate files from a
 template

Arguments:
  etcd-path      The etcd path to the key where our config is stored
  template-dir   The directory where template files suffixed with .tpl are located
  output-dir     Output directory for generated files (Defaults to template-dir if
                 not provided)

Options:
  -w, --watch            Watches the specified etcd key for changes and regenerates
                         templates if the key value changes
  -e, --etcd-endpoints   A Comma Separated list of etcd server endpoints
                         (Default=http://localhost:2379, Env=ETCD_ENDPOINTS)
  -h, --help             Display this help message and exit
```

# Single Shot Example
Given a file called `examples/test.conf.tpl` with the contents
```
key1={{ .value1 }}
key2={{ .value2 }}
key3=3
```
Set an etcd key with the appropriate json
```bash
$ etcdctl set /mailgun/configs/ord/service1 '{"value1": 1, "value2": 2}'
```
Now run etcd-template
```bash
$ etcd-template /mailgun/configs/ord/service1 examples
```
A file called `examples/test.conf` should have been generated
```bash
$ cat examples/test.conf
key1=1
key2=2
key3=3
```

# Watch for config changes
Use the `--watch` optional argument to have `etcd-template` generate new files
from the templates when a config in etcd changes.


# Installation
For local installation and testing template generation
```bash
$ go install github.com/thrawn01/etcd-template/cmd/./...
```

Install in a container with `Dockerfile`
```
ENV ETCD_TEMPLATE_VERSION v0.1
RUN curl -Lso /tmp/etcd-template.tar.gz \
    "https://github.com/thrawn01/etcd-template/releases/download/${ETCD_TEMPLATE_VERSION}/etcd-template-${ETCD_TEMPLATE_VERSION}.tar.gz" \
    && tar zxf /tmp/etcd-template.tar.gz -C /usr/local/bin \
    && rm /tmp/etcd-template.tar.gz
```
