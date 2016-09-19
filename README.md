# s6_exporter
[![CircleCI](https://circleci.com/gh/imgix/s6_exporter.svg?style=svg)](https://circleci.com/gh/imgix/s6_exporter)

Prometheus Exporter exposes the following metrics for each s6 service:

```
# HELP s6_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which s6_exporter was built.
# TYPE s6_exporter_build_info gauge
s6_exporter_build_info{branch="master",goversion="go1.6.1",revision="a0610c4",version="v0.0.1"} 1
# HELP s6_service_state_change_timestamp_seconds Unix timestamp of service's last state change.
# TYPE s6_service_state_change_timestamp_seconds gauge
s6_service_state_change_timestamp_seconds{service="foo"} 1.474291919e+09
# HELP s6_service_up State of s6 service, 1 = up, 0 = down
# TYPE s6_service_up gauge
s6_service_up{service="foo"} 1
# HELP s6_service_wanted_up Wanted state of s6 service, 1 = up, 0 = down
# TYPE s6_service_wanted_up gauge
s6_service_wanted_up{service="foo"} 1
```

Additionally to that, the exporter exposes a global error counter:
```
# HELP s6_exporter_errors_total Total number of errors s6_exporter encountered
# TYPE s6_exporter_errors_total counter
s6_exporter_errors_total 0
```

The exporter can be configured by the following command line flags:

```
Usage of ./s6_exporter:
  -d string
      Path to service directory (default "/etc/service")
  -h string
      Address to expose prometheus metrics on (default ":9164")
  -s string
      svstat binary name (default "s6-svstat")
  -version
      Print version information.
```
