## Overview

Application can test following things (with config names in brackets):
    
* HTTP(s) server can be reached from client's network (`http_test`)
* NTP server can be reached from client's network (`ntp_test`)
* TCP connection can be established from client's network (`tcp_out`)
* TCP connection can be established to client's machine from internet (`tcp_in`)
* UDP packets can reach client's machine from internet (`udp_in`)

## Build

To build application just run `make`

Executables will be placed in `./bin` directory.

## Run

### Server application

```
./bin/server <http_listen_port>
```

### Client application

```
./bin/client --config <test_configuration_file.json>
```

To save test results in json file, run following
```
./bin/client --config <test_configuration_file.json> --output_result --output_path=<output_result.json>
```

#### Client configuration file example

In order for tests `tcp_in` and `udp_in` to work, server application should be available and its address should be specified in `host` field. Also it should be possible to reach server app via HTTP.

```json
{
    "tests":[
        {
            "name":"ntp_test",
            "params":{
                "ntp_server":"2.ch.pool.ntp.org"
            }
        },
        {
            "name":"http_test",
            "params":{
                "host":"https://www.scion-architecture.net",
                "method":"GET",
                "timeout":10
            }
        },
        {
            "name":"tcp_out",
            "params":{
                "host":"www.zvv.ch",
                "port":"80",
                "request":"GET / HTTP/1.1\n\n",
                "compare_response":true,
                "timeout":10,
                "expected_response":"HTTP/1.1"
            }
        },
        {
            "name":"tcp_in",
            "params":{
                "host":"http://localhost:8080/tcp-test",
                "timeout":5,
                "my_port":"9999"
            }
        },
        {
            "name":"udp_in",
            "params":{
                "host":"http://localhost:8080/udp-test",
                "timeout":5,
                "my_port":"50000"
            }
        }
    ]
}
``` 