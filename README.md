# JRPCServer

JRPCServer is a simple library which allows you to rapidly expose a JSON response based RPC service which includes support for prometheus expont on '/metrics', registering with a Consul service if desired, and a health endpoint on '/health'.

## A Simple Example

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/riomhaire/jrpcserver/infrastructure/api/rpc"
	"github.com/riomhaire/jrpcserver/model"
	"github.com/riomhaire/jrpcserver/model/jrpcerror"
	"github.com/riomhaire/jrpcserver/usecases/defaultcommand"
)

// A simple example 'helloworld' program to show how use the framework
//
// to execute (after compiling):
//      helloworld -port=9999 -consul=consul:8500
//
func main() {

	name := flag.String("name", "helloworld", "name of service")
	path := flag.String("path", "/api/v1/helloworld", "Path to which to  <path>/<command>  points to action")
	consulHost := flag.String("consul", "", "consul host usually something like 'localhost:8500'. Leave blank if not required")
	port := flag.Int("port", 9999, "port to use")
	flag.Parse()

	config := model.DefaultConfiguration{
		Server: model.ServerConfig{
			ServiceName: *name,
			BaseURI:     *path,
			Port:        *port,
			Commands:    Commands(),
			Consul:      *consulHost,
		},
	}
	rpc.StartAPI(config) // Start service -  wont return
}

func Commands() []model.JRPCCommand {
	commands := make([]model.JRPCCommand, 0)

	commands = append(commands, model.JRPCCommand{"example.helloworld", HelloWorldCommand, false})
	commands = append(commands, model.JRPCCommand{"system.commands", defaultcommand.ListCommandsCommand, false})
	return commands
}

func HelloWorldCommand(config interface{}, metadata map[string]string, payload io.ReadCloser) (interface{}, jrpcerror.JrpcError) {
	data, err := ioutil.ReadAll(payload)
	if err != nil {
		return "", jrpcerror.JrpcError{500, err.Error()}
	} else {
		fmt.Println(string(data))
		response := fmt.Sprintf("Hello %v", string(data))

		return response, jrpcerror.JrpcError{}
	}
}

```

Hopefully this is fairly understandable (and is included within the examples/helloworld ). Essentially all you have to do is provide a list of functions and their command names and the library handles the rest.

If you call (using VSCode Rest syntax)

```curl
POST http://localhost:9999/api/v1/helloworld/example.helloworld
Content-Type: text/plain

fred
```

You would receive something like:

```text
HTTP/1.1 200 OK
Content-Type: application/json
Vary: Origin
X-Worker: warband
X-Worker-Version: UNKNOWN
Date: Mon, 17 Sep 2018 20:14:27 GMT
Content-Length: 55

{
  "code": 0,
  "error": "",
  "value": "Hello fred"
}
```

If you call:

```curl
GET http://localhost:9999/metrics
```

You would receive something like:

```text
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Length: 1747
Content-Type: text/plain; version=0.0.4
Vary: Origin
X-Worker: warband
X-Worker-Version: UNKNOWN
Date: Mon, 17 Sep 2018 20:15:29 GMT

# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 9
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.11"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.973248e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.973248e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.443459e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 672
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.234368e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.973248e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.3520768e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.063808e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 4872
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 0
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6584576e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 5544
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 13824
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 32528
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 32768
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.473924e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.055349e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 524288
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 524288
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.1891192e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 11
# HELP http_request_duration_microseconds The HTTP request latencies in microseconds.
# TYPE http_request_duration_microseconds summary
http_request_duration_microseconds{handler="prometheus",quantile="0.5"} NaN
http_request_duration_microseconds{handler="prometheus",quantile="0.9"} NaN
http_request_duration_microseconds{handler="prometheus",quantile="0.99"} NaN
http_request_duration_microseconds_sum{handler="prometheus"} 4973.204
http_request_duration_microseconds_count{handler="prometheus"} 1
# HELP http_request_size_bytes The HTTP request sizes in bytes.
# TYPE http_request_size_bytes summary
http_request_size_bytes{handler="prometheus",quantile="0.5"} NaN
http_request_size_bytes{handler="prometheus",quantile="0.9"} NaN
http_request_size_bytes{handler="prometheus",quantile="0.99"} NaN
http_request_size_bytes_sum{handler="prometheus"} 130
http_request_size_bytes_count{handler="prometheus"} 1
# HELP http_requests_total Total number of HTTP requests made.
# TYPE http_requests_total counter
http_requests_total{code="200",handler="prometheus",method="get"} 1
# HELP http_response_size_bytes The HTTP response sizes in bytes.
# TYPE http_response_size_bytes summary
http_response_size_bytes{handler="prometheus",quantile="0.5"} NaN
http_response_size_bytes{handler="prometheus",quantile="0.9"} NaN
http_response_size_bytes{handler="prometheus",quantile="0.99"} NaN
http_response_size_bytes_sum{handler="prometheus"} 1605
http_response_size_bytes_count{handler="prometheus"} 1
# HELP negroni_request_duration_milliseconds How long it took to process the request, partitioned by status code, method and HTTP path.
# TYPE negroni_request_duration_milliseconds histogram
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/health",service="helloworld",le="300"} 1
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/health",service="helloworld",le="1200"} 1
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/health",service="helloworld",le="5000"} 1
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/health",service="helloworld",le="+Inf"} 1
negroni_request_duration_milliseconds_sum{code="",method="GET",path="/health",service="helloworld"} 0.028667
negroni_request_duration_milliseconds_count{code="",method="GET",path="/health",service="helloworld"} 1
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/metrics",service="helloworld",le="300"} 1
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/metrics",service="helloworld",le="1200"} 1
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/metrics",service="helloworld",le="5000"} 1
negroni_request_duration_milliseconds_bucket{code="",method="GET",path="/metrics",service="helloworld",le="+Inf"} 1
negroni_request_duration_milliseconds_sum{code="",method="GET",path="/metrics",service="helloworld"} 0.026118
negroni_request_duration_milliseconds_count{code="",method="GET",path="/metrics",service="helloworld"} 1
negroni_request_duration_milliseconds_bucket{code="",method="POST",path="/api/v1/helloworld/example.helloworld",service="helloworld",le="300"} 12
negroni_request_duration_milliseconds_bucket{code="",method="POST",path="/api/v1/helloworld/example.helloworld",service="helloworld",le="1200"} 12
negroni_request_duration_milliseconds_bucket{code="",method="POST",path="/api/v1/helloworld/example.helloworld",service="helloworld",le="5000"} 12
negroni_request_duration_milliseconds_bucket{code="",method="POST",path="/api/v1/helloworld/example.helloworld",service="helloworld",le="+Inf"} 12
negroni_request_duration_milliseconds_sum{code="",method="POST",path="/api/v1/helloworld/example.helloworld",service="helloworld"} 0.10833
negroni_request_duration_milliseconds_count{code="",method="POST",path="/api/v1/helloworld/example.helloworld",service="helloworld"} 12
# HELP negroni_requests_total How many HTTP requests processed, partitioned by status code, method and HTTP path.
# TYPE negroni_requests_total counter
negroni_requests_total{code="",method="GET",path="/health",service="helloworld"} 1
negroni_requests_total{code="",method="GET",path="/metrics",service="helloworld"} 1
negroni_requests_total{code="",method="POST",path="/api/v1/helloworld/example.helloworld",service="helloworld"} 12
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.01
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1024
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 73
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 9.555968e+06
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.53721438365e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 8.08869888e+08
```

If you call the health endpoint:

```curl
GET http://localhost:9999/health
```

You would receive something like:

```text
HTTP/1.1 200 OK
Content-Type: application/json
Vary: Origin
X-Worker: warband
X-Worker-Version: UNKNOWN
Date: Mon, 17 Sep 2018 20:01:26 GMT
Content-Length: 16

{
  "status": "up"
}
```
