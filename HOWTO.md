# How to run a w3bstream node with docker

Suppose `$working_dir` is the directory you want to store your data.

## Install docker-compose

https://docker-docs.netlify.app/compose/install/

## Download docker-compose.yaml

```bash
cd $working_dir
curl https://raw.githubusercontent.com/machinefi/w3bstream/main/docker-compose.yaml > docker-compose.yaml
docker-compose up -d
```

You are all set.

## Customize settings

```bash
cd $working_dir
curl https://raw.githubusercontent.com/machinefi/w3bstream/main/.env.tmpl > .env
```

then modify the corresponding parameters in `.env`, and restart your docker
containers

```bash
docker-compose restart
```

# Run W3bstream node from code

If you are interested in diving into the code and run the node using a locally built docker, here is the steps of building the docker image from code.

### Build docker image from code

```bash
make build_backend_image
```

### Run server in docker containers

```bash
 make run_docker
 ```

### Stop server running in docker containers
 ```bash
 make stop_docker
 ```
### Delete docker resources
 ```bash
 make drop_docker
 ```

# How to interact with W3bstream Node Using CLI

### Login (fetch auth token)

command

```sh
# the default password is "iotex.W3B.admin"
echo '{"username":"admin","password":"iotex.W3B.admin"}' | http put :8888/srv-applet-mgr/v0/login 
```

output like

```json
{
  "accountID": "${account_id}",
  "expireAt": "2022-09-23T07:20:08.099601+08:00",
  "issuer": "srv-applet-mgr",
  "token": "${token}"
}
```

```sh
export TOK=${token}
```

### Login/Signup with wallet address

```shell
export MESSAGE=...   # siwe serailized message
export SIGNATURE=... # message signature
echo '{"message":"'$MESSAGE'","signature":"'$SIGNATURE'"}' | http put :8888/srv-applet-mgr/v0/login/eth
```

output like:

```json
{
  "accountID": "186912900253363206",
  "accountRole": "DEVELOPER",
  "expireAt": "2023-03-16T19:07:57.624481+08:00",
  "issuer": "iotex-w3bstream",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXlsb2FkIjoiMTg2OTEyOTAwMjUzMzYzMjA2IiwiaXNzIjoiaW90ZXgtdzNic3RyZWFtIiwiZXhwIjoxNjc4OTY0ODc3fQ.u7wLOBUeehHTURNY2L2d_F4u-dZ5sHnBBHZKujnpMRw"
}
```

### Get Account's Operator Address

```shell
http get :8888/srv-applet-mgr/v0/account/operatoraddr -A bearer -a $TOK

```

### Create your project with default config

command

```sh
export PROJECTNAME=${project_name}
echo '{"name":"'$PROJECTNAME'"}' | http post :8888/srv-applet-mgr/v0/project -A bearer -a $TOK
```

output like

```json
{
  "accountID": "11276794515805192",
  "channelState": true,
  "createdAt": "2023-05-03T05:39:17.835566714Z",
  "database": {
    "schemas": [
      {
        "schema": "public"
      }
    ]
  },
  "envs": {
    "env": null
  },
  "name": "demo",
  "projectID": "11276839333473280",
  "updatedAt": "2023-05-03T05:39:17.835567047Z"
}
```

### Create project with database and env vars configuration

```sh
export PROJECTDATABASE='{
  "schemas": [
    {
      "schema": "public",
      "tables": [
        {
          "name": "t_demo",
          "desc": "demo table",
          "cols": [
            {
              "name": "f_id",
              "constrains": {
                "datatype": "INT64",
                "autoincrement": true,
                "desc": "primary id"
              }
            },
            {
              "name": "f_name",
              "constrains": {
                "datatype": "TEXT",
                "length": 255,
                "desc": "name"
              }
            },
            {
              "name": "f_amount",
              "constrains": {
                "datatype": "FLOAT64",
                "desc": "amount"
              }
            },
            {
              "name": "f_income",
              "constrains": {
                "datatype": "DECIMAL",
                "length": 512,
                "decimal": 512,
                "default": "0",
                "desc": "income"
              }
            },
            {
              "name": "f_comment",
              "constrains": {
                "datatype": "TEXT",
                "default": "''",
                "null": true,
                "desc": "comment"
              }
            }
          ],
          "keys": [
            {
              "name": "primary",
              "isUnique": false,
              "columnNames": [
                "f_id"
              ]
            },
            {
              "name": "t_demo_ui_name",
              "isUnique": true,
              "columnNames": [
                "f_name"
              ]
            },
            {
              "name": "i_amount",
              "isUnique": false,
              "columnNames": [
                "f_amount"
              ]
            }
          ]
        }
      ]
    }
  ]
}'
export PROJECTENV='{
  "env": [
    ["envKey1", "envValue1"],
    ["envKey2", "envValue2"],
    ["envKey3", "envValue3"]
  ]
}'
echo '{"name":"'$PROJECTNAME'","database": '$PROJECTDATABASE',"envs":'$PROJECTENV'}' | http post :8888/srv-applet-mgr/v0/project -A bearer -a $TOK
```

You can define your own database model.

`schemas` defines database schema structure

`schemas[i].schema` defines schema name, default using `public` schema

`schemas[i].tables` defines table structures in schema.

`schemas[i].tables[i].name` defines table name(required)

`schemas[i].tables[i].desc` defines table description

`schemas[i].tables[i].cols` defines table columns

`schemas[i].tables[i].cols[i].name` defines column name

`schemas[i].tables[i].cols[i].constrains.datatype` defines column datatype

> column datatype

| column datatype | postgres datatype |
|:----------------|:------------------|
| "INT"           | integer           |
| "INT8"          | integer           |
| "INT16"         | integer           |
| "INT32"         | integer           |
| "INT64"         | bigint            |
| "UINT"          | integer           |
| "UINT8"         | integer           |
| "UINT16"        | integer           |
| "UINT32"        | integer           |
| "UINT64"        | bigint            |
| "FLOAT32"       | real              |
| "FLOAT64"       | double precision  |
| "TEXT"          | character varying |
| "BOOL"          | boolean           |
| "TIMESTAMP"     | bigint            |
| "DECIMAL"       | numeric           |
| "NUMERIC"       | numeric           |

> the timestamp stored in database using epoch timestamp (microsecond)

`schemas[i].tables[i].cols[i].constrains.datatype` defines column datatype

`schemas[i].tables[i].cols[i].constrains.length` defines `character varying`
length or `numeric` precision

`schemas[i].tables[i].cols[i].constrains.decimal` defines `numeric` scale

`schemas[i].tables[i].cols[i].constrains.default` defines column default value

`schemas[i].tables[i].cols[i].constrains.null` defines if column can be null

`schemas[i].tables[i].cols[i].constrains.autoincrement` defines if column is
autoincrement

> when column is autoincrement, integer defined as `serial` and bigint defined
> as `bigserial` in postgres datatype

`schemas[i].tables[i].cols[i].constrains.desc` column description

`schemas[i].tables[i].keys[i].name` defines index name

`schemas[i].tables[i].keys[i].isUnique` defines if index is unique

`schemas[i].tables[i].keys[i].columnNames` index related column names

> NOTE:
> if the key's name is `primary` or has suffix `pkey`, it defined as primary key
> of the table.
> the index name will be built by this pattern:
> 1. non-primary index: `tableName_[i|ui]_[columnName1]_[columnName2]_...`. if
     it is a unique index use `ui`, otherwise use `i` to split table name and
     index defines.
> 2. primary key: `tableName_primary`


output like:

```json
{
  "accountID": "11276794515805192",
  "channelState": true,
  "createdAt": "2023-05-03T06:44:05.132275513Z",
  "database": {
    "schemas": [
      {
        "schema": "public",
        "tables": [
          {
            "cols": [
              {
                "constrains": {
                  "autoincrement": true,
                  "datatype": "INT64",
                  "desc": "primary id"
                },
                "name": "f_id"
              },
              {
                "constrains": {
                  "datatype": "TEXT",
                  "desc": "name",
                  "length": 255
                },
                "name": "f_name"
              },
              {
                "constrains": {
                  "datatype": "FLOAT64",
                  "desc": "amount"
                },
                "name": "f_amount"
              },
              {
                "constrains": {
                  "datatype": "DECIMAL",
                  "decimal": 512,
                  "default": "0",
                  "desc": "income",
                  "length": 512
                },
                "name": "f_income"
              },
              {
                "constrains": {
                  "datatype": "TEXT",
                  "default": "''",
                  "desc": "comment",
                  "null": true
                },
                "name": "f_comment"
              }
            ],
            "desc": "demo table",
            "keys": [
              {
                "columnNames": [
                  "f_id"
                ],
                "name": "primary"
              },
              {
                "columnNames": [
                  "f_name"
                ],
                "isUnique": true,
                "name": "t_demo_ui_name"
              },
              {
                "columnNames": [
                  "f_amount"
                ],
                "name": "i_amount"
              }
            ],
            "name": "t_demo"
          }
        ]
      }
    ]
  },
  "envs": {
    "env": [
      [
        "envKey1",
        "envValue1"
      ],
      [
        "envKey2",
        "envValue2"
      ],
      [
        "envKey3",
        "envValue3"
      ]
    ]
  },
  "name": "demo2",
  "projectID": "11276843314064388",
  "updatedAt": "2023-05-03T06:44:05.132275805Z"
}
```

### Create or update project configurations after project created

```sh
echo $PROJECTENV | http post :8888/srv-applet-mgr/v0/project_config/x/$PROJECTNAME/PROJECT_ENV -A bearer -a $TOK
echo $PROJECTDATABASE | http post :8888/srv-applet-mgr/v0/project_config/x/$PROJECTNAME/PROJECT_DATABASE -A bearer -a $TOK
```

### Review your projects and project configurations

```shell
http get :8888/srv-applet-mgr/v0/project/x/$PROJECTNAME/data -A bearer -a $TOK ## fetch project by name
http get :8888/srv-applet-mgr/v0/project_config/x/demo/PROJECT_ENV -A bearer -a $TOK  # fetch project env configuration
http get :8888/srv-applet-mgr/v0/project_config/x/demo/PROJECT_DATABASE -A bearer -a $TOK  # fetch project database configuration
http get :8888/srv-applet-mgr/v0/project/datalist -A bearer -a $TOK # fetch project list you created
```

### Create and deploy applet under project created previously

```sh
export WASMFILE=build/wasms/log.wasm
export WASMNAME=log.wasm
export APPLETNAME=log
http --form post :8888/srv-applet-mgr/v0/applet/x/$PROJECTNAME file@$WASMFILE info='{"appletName":"'$APPLETNAME'","wasmName":"'$WASMNAME'"}' -A bearer -a $TOK 
```

output like

```json
{
  "appletID": "11276843999120385",
  "createdAt": "2023-05-03T06:55:14.131370253Z",
  "instance": {
    "appletID": "11276843999120385",
    "createdAt": "2023-05-03T06:55:14.146653045Z",
    "instanceID": "11276843999135746",
    "state": "STARTED",
    "updatedAt": "2023-05-03T06:55:14.146653128Z"
  },
  "name": "11276843999120386",
  "projectID": "11276843314064388",
  "resource": {
    "createdAt": "2023-05-03T06:55:14.112226878Z",
    "md5": "30b11f90b1d7453474496f5cc42f0869",
    "path": "30b11f90b1d7453474496f5cc42f0869",
    "resourceID": "11276843999092744",
    "updatedAt": "2023-05-03T06:55:14.112227086Z"
  },
  "resourceID": "11276843999092744",
  "updatedAt": "2023-05-03T06:55:14.131370336Z"
}
```

> you can create applet with deploy state and event routing strategy, if no
> strategy configuration, the default
> strategy with `DEFAULT` event type and `start` handler will be created

`info.appletName` defined the unique applet name under the project

`info.wasmName` the resource filename

`info.wasmMd5` the wasm file md5, if it is not empty, w3bstream node will check
md5 sum

`info.wasmCache` wasm cache config

`info.wasmCache.mode` cache mode, enumerated in `MEMORY` and `REDIS`, `MEMORY`
is default

`info.strategies` event routing strategies

`info.strategies[i].eventType` routing with eventType (user defined)

`info.strategies[i].handler` routing wasm handler name (wasm exported)

### create instance of the applet you created before

```sh
export APPLETID=11276843999120385 ## created before
http post :8888/srv-applet-mgr/v0/deploy/applet/$APPLETID -A bearer -a $TOK
## with cache config
export WASMCACHECONFIG='{"cache":{"mode": "REDIS"}}'
echo $WASMCACHECONFIG | http post :8888/srv-applet-mgr/v0/deploy/applet/$APPLETID -A bearer -a $TOK
```

output like

```json
{
  "appletID": "11276843999120385",
  "createdAt": "2023-05-03T07:13:28.409513718Z",
  "instanceID": "11276845119659014",
  "state": "CREATED",
  "updatedAt": "2023-05-03T07:13:28.409514176Z"
}
```

if the instance is already created, output like

```json
{
  "canBeTalk": true,
  "code": 409999008,
  "desc": "11276843999120385",
  "fields": null,
  "id": "",
  "key": "MultiInstanceDeployed",
  "msg": "Multi Instance Deployed",
  "sources": [
    "srv-applet-mgr@v1.1.0-rc3-6-gbf5cbc0"
  ]
}
```

### Control instance

```sh
export INSTANCEID=11276845119659014 ## created before
export DEPLOYCMD=START
http put :8888/srv-applet-mgr/v0/deploy/$INSTANCEID/$DEPLOYCMD -A bearer -a $TOK
```

deploy command enumerated in `START` and `HUNGUP`

`START` change the instance state to `STARTED`

`HUNGUP` change the instance state to `STOPPED`

### Update applet and redeploy instance

```sh
http --form put :8888/srv-applet-mgr/v0/applet/$APPLETID file@$WASMFILE info='{"appletName":"'$APPLETNAME'","wasmName":"'$WASMNAME'","start":true}' -A bearer -a $TOK
```

### Register publisher

```sh
export PUBNAME=mobile    # device name
export PUBKEY=mn20130503 # device unique identity, usually it is device's machine number or serial number
echo '{"name":"'$PUBNAME'", "key":"'$PUBKEY'"}' | http post :8888/srv-applet-mgr/v0/publisher/x/$PROJECTNAME -A bearer -a $TOK
```

output like

```sh
{
    "createdAt": "2023-05-03T16:13:16.343103+08:00",
    "key": "mn20130503",
    "name": "mobile",
    "projectID": "11276843314064388",
    "publisherID": "155392036869560322",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXlsb2FkIjoiMTU1MzkyMDM2ODY5NTYwMzIyIiwiaXNzIjoiaW90ZXgtdzNic3RyZWFtIn0.OHME3ij5MaJcvekctgYvosQ8DIo-K-guQbYPbQAdyYo",
    "updatedAt": "2023-05-03T16:13:16.343103+08:00"
}
```

> the `token` responded is used for validating publisher when publishing event.

### Review registered publisher

```sh
http get :8888/srv-applet-mgr/v0/publisher/x/$PROJECTNAME -A bearer -a $TOK
```

output like:

```json
{
  "data": [
    {
      "createdAt": "2023-05-03T16:13:16+08:00",
      "key": "mn20130503",
      "name": "mobile",
      "projectID": "11276843314064388",
      "publisherID": "155392036869560322",
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXlsb2FkIjoiMTU1MzkyMDM2ODY5NTYwMzIyIiwiaXNzIjoiaW90ZXgtdzNic3RyZWFtIn0.OHME3ij5MaJcvekctgYvosQ8DIo-K-guQbYPbQAdyYo",
      "updatedAt": "2023-05-03T16:13:16+08:00"
    }
  ],
  "total": 1
}
```

### Create strategy for applet created before

Create a strategy of handler in applet and eventType

```sh
export EVENTTYPE=mobile_geo
export HANDLER=handle_geo_data
echo '{"appletID":"'$APPLETID'", "eventType":"'$EVENTTYPE'", "handler":"'$HANDLER'"}' | http post :8888/srv-applet-mgr/v0/strategy/x/$PROJECTNAME -A bearer -a $TOK
```

output like:

```json
{
  "appletID": "11276843999120385",
  "createdAt": "2023-05-03T16:17:40.942225+08:00",
  "eventType": "mobile_geo",
  "handler": "handle_geo_data",
  "projectID": "11276843314064388",
  "strategyID": "155392037140510721",
  "updatedAt": "2023-05-03T16:17:40.942225+08:00"
}
```

### Review strategies under current project

```sh
http get :8888/srv-applet-mgr/v0/strategy/x/$PROJECTNAME/datalist -A bearer -a $TOK
```

output like:

```json
{
  "data": [
    {
      "appletID": "11276843999120385",
      "createdAt": "2023-05-03T16:17:40+08:00",
      "eventType": "mobile_geo",
      "handler": "handle_geo_data",
      "projectID": "11276843314064388",
      "strategyID": "155392037140510721",
      "updatedAt": "2023-05-03T16:17:40+08:00"
    },
    {
      "appletID": "11276843999120385",
      "createdAt": "2023-05-03T14:55:14+08:00",
      "eventType": "DEFAULT",
      "handler": "start",
      "projectID": "11276843314064388",
      "strategyID": "11276843999125505",
      "updatedAt": "2023-05-03T14:55:14+08:00"
    }
  ],
  "total": 2
}
```

### Publish event through http

```sh
export TOPIC=${pub_topic} ## intact project name(required)
export PUBTOK=${pub_token} ## created before(required)
export EVENTTYPE=mobile_geo # default means start handler
export EVENTID=`uuidgen` ## this id is used for tracing event(recommended)
export PAYLOAD=${payload} ## set your payload
export TIMESTAMP=`date +%s` ## event pub timestamp(recommended)
http post :8889/srv-applet-mgr/v0/event/$TOPIC\?eventType=$EVENTTYPE\&eventID=$EVENTID\&timestamp=$TIMESTAMP --raw=$PAYLOAD -A bearer -a $PUBTOK 
```

> note event handler service using 8889 for default

output like

```json
{
  "channel": "aid_11276794515805192_demo2",
  "eventID": "3d5d76d6-24be-4e47-9f44-cac2b4855e1a_w3b",
  "publisherID": "155392036869560322",
  "results": [
    {
      "appletName": "log",
      "code": -1,
      "error": "instance not running",
      "handler": "handle_geo_data",
      "instanceID": "11276845119659014",
      "returnValue": null
    }
  ]
}
```

### Publish event through mqtt (use `pub_client` CLI)

```sh
./pub_client -topic $TOPIC -token $PUBTOK -data $PAYLOAD
```

server log like

```json
{
  "@lv": "info",
  "@prj": "srv-applet-mgr",
  "@ts": "20221017-092252.877+08:00",
  "msg": "sub handled",
  "payload": {
    "payload": "..."
  }
}
```

the `pub_client` sends event message through mqtt broker using protobuf
encoding. the event defined as follows

| field name       | protobuf filed    | protobuf seq | datatype | requirement | comment                                                                            |
|:-----------------|:------------------|:-------------|:---------|:------------|:-----------------------------------------------------------------------------------|
| Header           | header            | 1            | object   | required    |                                                                                    |
| Header.EventType | header.event_type | header.1     | string   | recommended | event type(user defined) for event routing, according to strategies created before |
| Header.PubId     | header.pub_id     | header.2     | string   | recommended | publisher id, usually it is the device machine number, you register before         |
| Header.Token     | header.token      | header.3     | string   | required    | publisher token. it contains the publisher id (or DID)                             |
| Header.PubTime   | header.pub_time   | header.4     | int64    | recommended | message timestamp when published, use unix epoch timestamp (in UTC)                |
| Header.EventId   | header.event_id   | header.5     | string   | recommended | event id is the unique identity of this message related with the publisher         |
| Payload          | payload           | 2            | bytes    | -           | message payload                                                                    |

### Data cleanup

Be careful.
It will delete anything in the project, contains applet, publisher, strategy
etc.

```sh
## delete project, all configurations in the project, contains applet, publisher,
## strategy and instances will be deleted, the database model will be kept.
http delete :8888/srv-applet-mgr/v0/project/x/$PROJECTNAME -A bearer -a $TOK
## delete applet, all instances, strategy and the configurations will be deleted
http delete :8888/srv-applet-mgr/v0/applet/data/$APPLETID -A bearer -a $TOK
## delete instance, the instance will be released from host memory and the configurations
## will be deleted.
http delete :8888/srv-applet-mgr/v0/deploy/data/$INSTANCEID -A bearer -a $TOK
```

### Post blockchain contract event log monitor

```sh
echo '{"eventType": "DEFAULT", "chainID": 4690, "contractAddress": "${contractAddress}","blockStart": ${blockStart},"blockEnd": ${blockEnd},"topic0":"${topic0}"}' | http :8888/srv-applet-mgr/v0/monitor/x/$PROJECTNAME/contract_log -A bearer -a $TOK
```

output like

```json
{
  "blockCurrent": 16737070,
  "blockEnd": 16740080,
  "blockStart": 16737070,
  "chainID": 4690,
  "contractAddress": "${contractAddress}",
  "contractlogID": "2162022028435556",
  "createdAt": "2022-10-19T21:21:30.220198+08:00",
  "eventType": "ANY",
  "projectName": "${projectName}",
  "topic0": "${topic0}",
  "updatedAt": "2022-10-19T21:21:30.220198+08:00"
}
```

delete it

```sh
export CONTRACTLOGID=${contractlogID}
http delete :8888/srv-applet-mgr/v0/monitor/x/$PROJECTNAME/contract_log/$CONTRACTLOGID -A bearer -a $TOK
```

### Post blockchain transaction monitor

```sh
echo '{"eventType": "DEFAULT", "chainID": 4690, "txAddress": "${txAddress}"}' | http :8888/srv-applet-mgr/v0/monitor/x/$PROJECTNAME/chain_tx -A bearer -a $TOK
```

output like

```json
{
  "chainID": 4690,
  "chaintxID": "2724127039316068",
  "createdAt": "2022-10-21T10:35:06.498594+08:00",
  "eventType": "ANY",
  "projectName": "testproject",
  "txAddress": "${txAddress}",
  "updatedAt": "2022-10-21T10:35:06.498594+08:00"
}
```

delete it

```sh
export CHAINTXID=${chaintxID}
http delete :8888/srv-applet-mgr/v0/monitor/x/$PROJECTNAME/chain_tx/$CHAINTXID -A bearer -a $TOK
```

### Post blockchain height monitor

```sh
echo '{"eventType": "DEFAULT", "chainID": 4690, "height": ${height}}' | http :8888/srv-applet-mgr/v0/monitor/x/$PROJECTNAME/chain_height -A bearer -a $TOK
```

output like

```json
{
  "chainHeightID": "2727219570933860",
  "chainID": 4690,
  "createdAt": "2022-10-21T10:47:23.815552+08:00",
  "eventType": "ANY",
  "height": 16910805,
  "projectName": "testproject",
  "updatedAt": "2022-10-21T10:47:23.815553+08:00"
}
```

delete it

```sh
export CHAINHEIGHTID=${chainHeightID}
http delete :8888/srv-applet-mgr/v0/monitor/x/$PROJECTNAME/chain_height/$CHAINHEIGHTID -A bearer -a $TOK
```

