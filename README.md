# s6_exporter
[![CircleCI](https://circleci.com/gh/imgix/s6_exporter.svg?style=svg)](https://circleci.com/gh/imgix/s6_exporter)

Prometheus Exporter exposes the following metrics for each s6 service:

```
# HELP s6\_exporter\_build\_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which s6\_exporter was built.
# TYPE s6\_exporter\_build\_info gauge
s6\_exporter\_build\_info{branch="master",goversion="go1.6.1",revision="a0610c4",version="v0.0.1"} 1
# HELP s6\_service\_state\_change\_timestamp\_seconds Unix timestamp of service's last state change.
# TYPE s6\_service\_state\_change\_timestamp\_seconds gauge
s6\_service\_state\_change\_timestamp\_seconds{service="foo"} 1.474291919e+09
# HELP s6\_service\_up State of s6 service, 1 = up, 0 = down
# TYPE s6\_service\_up gauge
s6\_service\_up{service="foo"} 1
# HELP s6\_service\_wanted\_up Wanted state of s6 service, 1 = up, 0 = down
# TYPE s6\_service\_wanted\_up gauge
s6\_service\_wanted\_up{service="foo"} 1
```

Additionally to that, the exporter exposes a global error counter:
```
# HELP s6\_exporter\_errors\_total Total number of errors s6\_exporter encountered
# TYPE s6\_exporter\_errors\_total counter
s6\_exporter\_errors\_total 0
```

The exporter can be configured by the following command line flags:

```
Usage of ./s6\_exporter:
  -d string
      Path to service directory (default "/etc/service")
  -h string
      Address to expose prometheus metrics on (default ":9164")
  -s string
      svstat binary name (default "s6-svstat")
  -version
      Print version information.
```
