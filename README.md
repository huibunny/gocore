# gocore


## consul

### kv

In consul kv must contain the below options:

```bash

consul: 
  # health check api
  checkapi: '/healthz'
  # health check interval, seconds
  interval: 3
  # health check timeout, seconds
  timeout: 5


```

### reference

* [A consul key/value tool for golang](https://golangexample.com/a-consul-key-value-tool-for-golang/)

* [Go configuration with fangs](github.com/spf13/viper)
