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

### Create your project without schema

command

```sh
export PROJECTNAME=${project_name}
echo '{"name":"'$PROJECTNAME'"}' | http :8888/srv-applet-mgr/v0/project -A bearer -a $TOK
```

output like

```json
{
  "accountID": "${account_id}",
  "createdAt": "2022-10-14T12:50:26.890393+08:00",
  "name": "${project_name}",
  "projectID": "${project_id}",
  "updatedAt": "2022-10-14T12:50:26.890407+08:00"
}
```

### Create project database schema for wasm db storage

```sh
export PROJECTSCHEMA='{
  "tables": [
    {
      "name": "tbl",
      "desc": "test table",
      "cols": [
        {
          "name": "f_username",
          "constrains": {
            "datatype": "TEXT",
            "length": 255,
            "desc": "user name"
          }
        },
        {
          "name": "f_gender",
          "constrains": {
            "datatype": "UINT8",
            "length": 255,
            "default": "0",
            "desc": "user name"
          }
        }
      ],
      "keys": [
        {
          "name": "ui_username",
          "isUnique": true,
          "columnNames": [
            "f_username"
          ]
        }
      ],
      "withSoftDeletion": true,
      "withPrimaryKey": true
    }
  ]
}'
echo $PROJECTSCHEMA | http post :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_SCHEMA -A bearer -a $TOK
```

### Create or update project env vars

```sh
export PROJECTENV='[["key1","value1"],["key2","value2"],["key3","value3"]]'
echo '{"env":'$PROJECTENV'}' | http post :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_ENV -A bearer -a $TOK
```

> the database for wasm storage is configured by w3bstream server and the name
> of schema is name of project.

### Create project with project env vars and schema

```sh
echo '{"name":"'$PROJECTNAME'","envs":'$PROJECTENV',"schema":'$PROJECTSCHEMA'}' | http post :8888/srv-applet-mgr/v0/project -A bearer -a $TOK
```

output like 
```json
{
    "envs": [
        [
            "key1",
            "value1"
        ],
        [
            "key2",
            "value2"
        ],
        [
            "key3",
            "value3"
        ]
    ],
    "project": {
        "accountID": "186913331796320263",
        "createdAt": "2023-03-27T14:52:20.037217+08:00",
        "name": "demo2",
        "projectID": "186913955765090305",
        "updatedAt": "2023-03-27T14:52:20.037217+08:00"
    },
    "schema": {
        "name": "wasm_project__demo2"
    }
}
```

### Review your project config

```shell
http get :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_SCHEMA -A bearer -a $TOK
http get :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_ENV -A bearer -a $TOK
```

### Create and deploy applet

upload wasm script


```sh
## set env vars
export WASMFILE=${wasm_path}
export APPLETNAME=${applet_name}
export WASMNAME=${wasm_name}
http --form post :8888/srv-applet-mgr/v0/applet/$PROJECTNAME file@$WASMFILE info='{"appletName":"'$APPLETNAME'","wasmName":"'$WASMNAME'","strategies":[{"eventType":"DEFAULT","handler":"start"}]}' -A bearer -a $TOK
```

output like

```json
{
  "appletID": "${apple_id}",
  "createdAt": "2022-10-14T12:53:10.590926+08:00",
  "name": "${applet_name}",
  "projectID": "${project_id}",
  "updatedAt": "2022-10-14T12:53:10.590926+08:00"
}
```

deploy applet

```sh
export APPLETID=${applet_id}
http post :8888/srv-applet-mgr/v0/deploy/applet/$APPLETID -A bearer -a $TOK

# in-memory cache is used in default, to use redis cache:
# echo '{"cache":{"mode": "REDIS"}}' | http post :8888/srv-applet-mgr/v0/deploy/applet/$APPLETID -A bearer -a $TOK
```

output like

```json
{
  "instanceID": "${instance_id}",
  "instanceState": "CREATED"
}
```

start applet

```sh
export INSTANCEID=${instance_id}
http put :8888/srv-applet-mgr/v0/deploy/$INSTANCEID/START -A bearer -a $TOK
```


### Register publisher

```sh
export PUBNAME=${publisher_name}
export PUBKEY=${publisher_unique_key} # global unique
echo '{"name":"'$PUBNAME'", "key":"'$PUBKEY'"}' | http post :8888/srv-applet-mgr/v0/publisher/$PROJECTNAME -A bearer -a $TOK
```

output like

```sh
{
    "createdAt": "2022-10-16T12:28:49.628716+08:00",
    "key": "${publisher_unique_key}",
    "name": "${publisher_name}",
    "projectID": "935772081365103",
    "publisherID": "${pub_id}",
    "token": "${pub_token}",
    "updatedAt": "2022-10-16T12:28:49.628716+08:00"
}
```

### Config Strategy

Create a strategy of handler in applet and eventType

```sh
export EVENTTYPE=${event_type}
export HANDLER=${applet_handler}
echo '{"strategies":[{"appletID":"'$APPLETID'", "eventType":"'$EVENTTYPE'", "handler":"'$HANDLER'"}]}' | http post :8888/srv-applet-mgr/v0/strategy/$PROJECTNAME -A bearer -a $TOK
```

get strategy info in the applet

```sh
http -v get :8888/srv-applet-mgr/v0/strategy/$PROJECTNAME appletID==$APPLETID -A bearer -a $TOK
```

### Publish event to server by http

```sh
export PUBTOKEN=${pub_token}
export EVENTTYPE=DEFAULT # default means start handler
export EVENTID=`uuidgen`
export PAYLOAD=${payload} # set your payload
echo '{"events":[{"header":{"event_id":"'$EVENTID'","event_type":"'$EVENTTYPE'","pub_id":"'$PUBKEY'","pub_time":'`date +%s`',"token":"'$PUBTOKEN'"},"payload":"'`echo $PAYLOAD | base64 -w 0`'"}]}' | http post :8888/srv-applet-mgr/v0/event/$PROJECTNAME
```

output like

```json
[
  {
    "eventID": "78C77DA7-8CE3-4E78-B970-95B685B02409",
    "projectName": "test",
    "wasmResults": [
      {
        "code": 0,
        "errMsg": "",
        "instanceID": "2612094299059956738"
      }
    ]
  }
]
```

that means some instance handled this event successfully

### Delete project

Be careful.
It will delete anything in the project, contains applet, publisher, strategy
etc.

```sh
http delete :8888/srv-applet-mgr/v0/project/$PROJECTNAME -A bearer -a $TOK
```

### Publish event to server through MQTT

- make publishing client

```sh
make build_pub_client
```

- try to publish a message

* event json message

```json
{
  "header": {
    "event_type": '$EVENTTYPE',
    "pub_id": "'$PUBKEY'",
    "pub_time": '`date +%s`',
    "token": "'$PUBTOKEN'"
  },
  "payload": "xxx yyy zzz"
}
```

* event_type: 0x7FFFFFFF any type
* pub_id: the unique publisher id assiged when publisher registering
* token: empty if dont have
* pub_time: timestamp when message published

```sh
# -c means published content
# -t means mqtt topic, the target project name created before
export PAYLOAD=${payload}
cd build/pub_client && ./pub_client -c '{"header":{"event_type":"'$EVENTTYPE'","pub_id":"'$PUBKEY'","pub_time":'`date +%s`',"token":"'$PUBTOKEN'"},"payload":"'`echo $PAYLOAD | base64 -w 0`'"}' -t $PROJECTNAME
```

server log like

```json
{
  "@lv": "info",
  "@prj": "srv-applet-mgr",
  "@ts": "20221017-092252.877+08:00",
  "msg": "sub handled",
  "payload": {
    "payload": "xxx yyy zzz"
  }
}
```

### Post blockchain contract event log monitor

```sh
echo '{"eventType": "DEFAULT", "chainID": 4690, "contractAddress": "${contractAddress}","blockStart": ${blockStart},"blockEnd": ${blockEnd},"topic0":"${topic0}"}' | http :8888/srv-applet-mgr/v0/monitor/contract_log/$PROJECTNAME -A bearer -a $TOK
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

### Post blockchain transaction monitor

```sh
echo '{"eventType": "DEFAULT", "chainID": 4690, "txAddress": "${txAddress}"}' | http :8888/srv-applet-mgr/v0/monitor/chain_tx/$PROJECTNAME -A bearer -a $TOK
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

### Post blockchain height monitor

```sh
echo '{"eventType": "DEFAULT", "chainID": 4690, "height": ${height}}' | http :8888/srv-applet-mgr/v0/monitor/chain_height/$PROJECTNAME -A bearer -a $TOK
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

### remove instance

```shell
export INSTANCEID=${instance_id}
http put :8888/srv-applet-mgr/v0/deploy/$INSTANCEID/REMOVE -A bearer -a $TOK 
```

### remove applet

> the instance will be stopped and removed

```shell
export APPLETID=${applet_id}
http delete :8888/srv-applet-mgr/v0/applet/$APPLETID -A bearer -a $TOK
```

### remove project

> the applets and the related instances included in this project will be stopped and removed

```shell
export PROJECTNAME=${project_name}
http delete :8888/srv-applet-mgr/v0/project/$PROJECTNAME -A bearer -a $TOK
```

### eth sigin/signup

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
