# service env manual

## Purpose

The purpose of this document is to help SRE and co-developers quickly track service configuration changes.

## Project build SOP

### Building structure

```
└── root ## repo root
    ├── Makefile
    ├── build
    │   ├── target1
    │   ├── target2
    │   ├── ...
    │   └── targetN
    └── cmd
        ├── target1
        │   ├── Makefile
        │   └── Dockerfile
        ├── target2
        │   ├── Makefile
        │   └── Dockerfile
        ├── ...
        └── targetN
            ├── Makefile
            └── Dockerfile
```

> flag * means the make entry is required

Entries in root/Makefile:

1. targets\*: for building all assets(binaries). in this entry, it will traverse all directories in `root/cmd/`,
   if `Makefile` exists, run `make target`
2. images: for building all images(docker). in this entry, it will traverse all directories in `root/cmd/`,
   if `Dockerfile` exists, run `make image`
3. test\*: project level testing entry.

Entries in root/cmd/Makefile:

1. target\*: building binary
2. image: build docker image

## Env Format

```
[SERVICE_NAME]__[CONFIG_GROUP]_[CONFIG_ELEMENT] = [CONFIG_VALUE]
```

eg:

```yaml
SRV_APPLET_MGR__Logger_Format: JSON
```

the configuration above means use `JSON` log format under `SRV_APPLET_MGR`

## Config Description

### RobotNotifier

```yaml
SRV_APPLET_MGR__RobotNotifier_Env: ""     ## service env. eg dev-staging, prod 
SRV_APPLET_MGR__RobotNotifier_Secret: ""  ## lark group secret, default ''
SRV_APPLET_MGR__RobotNotifier_URL: ""     ## required: lark group webhook url, 
SRV_APPLET_MGR__RobotNotifier_Vendor: ""  ## robot vendor. eg Lark, DingTalk WeWork
```

### WasmDBConfig

```yaml
SRV_APPLET_MGR__WasmDBConfig_Endpoint: ""           ## wasm database endpoint, default ''
SRV_APPLET_MGR__WasmDBConfig_MaxConnection: "2"     ## wasm database max connection for each wasm instance, default 2
SRV_APPLET_MGR__WasmDBConfig_ConnMaxLifetime: "20s" ## wasm database max connection lifetime default 20 seconds
SRV_APPLET_MGR__WasmDBConfig_PoolSize: "2"          ## wasm database connection pool size default 2
```

### Logger

```yaml
SRV_APPLET_MGR__NewLogger_Format: "JSON"            ## log format default `JSON`, use `JSON` or `TEXT`
SRV_APPLET_MGR__NewLogger_Level: "info"             ## enum in `error`, `warn`, `debug`, `info`, default `debug` suggested `info`
SRV_APPLET_MGR__NewLogger_Output: "ALWAYS"          ## default `ALWAYS`, enums in `ALWAYS` `ON_FAILURE` and `NEVER`, output to trace collector
SRV_APPLET_MGR__NewLogger_Service: "srv-applet-mgr" ## service name
SRV_APPLET_MGR__NewLogger_Version: "unknown"        ## service version
```

### Tracer

```yaml
SRV_APPLET_MGR__Tracer_DebugMode: "true"                      ## if enable tracer debug mode
SRV_APPLET_MGR__Tracer_GrpcEndpoint: "http://127.0.0.1:4317"  ## GRPC collector endpoint, default use GRPC collector
SRV_APPLET_MGR__Tracer_HttpEndpoint: "http://127.0.0.1:4218"  ## HTTP collector endpoint 
SRV_APPLET_MGR__Tracer_InstanceID: "xxx"                      ## unique instance id to identify service
SRV_APPLET_MGR__Tracer_ServiceName: "srv-applet-mgr"          ## service name
SRV_APPLET_MGR__Tracer_ServiceVersion:                        ## service version
SRV_APPLET_MGR__Tracer_TLS_Ca: ""                             ## endpoint TLS configurations, use `value` or `file path`
SRV_APPLET_MGR__Tracer_TLS_CaPath: ""
SRV_APPLET_MGR__Tracer_TLS_Crt: ""
SRV_APPLET_MGR__Tracer_TLS_CrtPath: ""
SRV_APPLET_MGR__Tracer_TLS_Key: ""
SRV_APPLET_MGR__Tracer_TLS_KeyPath: ""
```

### Task Manager MQ 

```yaml
SRV_APPLET_MGR__Mq_Channel: ""              ## channel name, if empty use env `PRJ_NAME`
SRV_APPLET_MGR__Mq_Limit: "1024"            ## queue limit default 1024
SRV_APPLET_MGR__Mq_PushQueueTimeout: "1s"   ## push timeout default 1s
SRV_APPLET_MGR__Mq_Store: "MEM"             ## support mem only now
SRV_APPLET_MGR__Mq_WorkerCount: "256"       ## task worker count, default 256
```