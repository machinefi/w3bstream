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
echo '{"username":"admin","password":"${password}"}' | http put :8888/srv-applet-mgr/v0/login
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

### Create project with database schema for wasm db storage

```sh
export PROJECTSCHEMA='{
  "tables": [
    {
      "name": "t_demo_table",
      "desc": "demo",
      "cols": [
        {
          "name": "f_autoinc_id",
          "constrains": {
            "datatype": "INT64",
            "length": 40,
            "autoincrement": true,
            "null": true,
            "desc": "datatype: bigserial"
          }
        },
        {
          "name": "f_text",
          "constrains": {
            "datatype": "TEXT",
            "default": "",
            "length": 128,
            "desc": "datatype: varchar(128)"
          }
        },
        {
          "name": "f_double_precious",
          "constrains": {
            "datatype": "FLOAT64",
            "default": "0"
          }
        },
        {
          "name": "f_decimal_with_precision_and_scale",
          "constrains": {
            "datatype": "DECIMAL",
            "length": 128,
            "decimal": 512,
            "default": "0"
          }
        },
        {
          "name": "f_numeric_with_precision_and_scale",
          "constrains": {
            "datatype": "NUMERIC",
            "length": 512,
            "decimal": 128,
            "default": "0"
          }
        },
        {
          "name": "f_decimal_default",
          "constrains": {
            "datatype": "DECIMAL"
          }
        },
        {
          "name": "f_numeric_default",
          "constrains": {
            "datatype": "NUMERIC"
          }
        },
        {
          "name": "f_timestamp",
          "constrains": {
            "datatype": "TIMESTAMP",
            "default": ""
          }
        }
      ],
      "keys": [
        {
          "name": "primary",
          "isUnique": true,
          "columnNames": [
            "f_autoinc_id"
          ]
        }
      ]
    }
  ]
}'
echo $PROJECTSCHEMA | http post :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_SCHEMA -A bearer -a $TOK
```

the exactly sql of creating is 

```sql
CREATE TABLE IF NOT EXISTS t_demo_table (
    f_autoinc_id bigserial,
    f_text character varying(128) NOT NULL,
    f_double_precious double precision NOT NULL DEFAULT '0'::double precision,
    f_decimal_with_precision_and_scale decimal(128,512) NOT NULL DEFAULT '0'::decimal,
    f_numeric_with_precision_and_scale numeric(512,128) NOT NULL DEFAULT '0'::numeric,
    f_decimal_default decimal NOT NULL,
    f_numeric_default numeric NOT NULL,
    f_timestamp bigint NOT NULL,
    f_created_at bigint NOT NULL DEFAULT '0'::bigint,
    f_updated_at bigint NOT NULL DEFAULT '0'::bigint,
    PRIMARY KEY (f_autoinc_id)
);
```

### Create or update project env vars

```sh
export PROJECTENV='[["key1","value1"],["key2","value2"],["key3","value3"]]'
echo '{"env":'$PROJECTENV'}' | http post :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_ENV -A bearer -a $TOK
```

> the database for wasm storage is configured by w3bstream server and the name
> of schema is name of project.

### Create project mqtt broker
```sh
export KEY=`cat $KEYFILE | base64 -w 0`
export CRT=`cat $CRTFILE | base64 -w 0`
export CA=`cat $CAILE | base64 -w 0`
export BROKER='{"scheme":"mqtts","host":"127.0.0.1","port":8883,"username":"applet_management","password":"PaSsW0rD","topics":["/device/#","/backend/#"],"tls":{"key":"'$KEY'","crt":"'$CRT'","ca":"'$CA'"}}'
echo $BROKER | http post :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_MQTT_BROKER -A bearer -a $TOK
```

### Review your project config

```shell
http get :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_SCHEMA -A bearer -a $TOK
http get :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_ENV -A bearer -a $TOK
http get :8888/srv-applet-mgr/v0/project_config/$PROJECTNAME/PROJECT_MQTT_BROKER -A bearer -a $TOK
```

### Create and deploy applet

upload wasm script


```sh
## set env vars
export WASMFILE=${wasm_path}
http --form post :8888/srv-applet-mgr/v0/applet/$PROJECTNAME file@$WASMFILE info='{"appletName":"log","wasmName":"log.wasm","strategies":[{"eventType":"DEFAULT","handler":"start"}]}' -A bearer -a $TOK
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
```

output like

```json
{
  "instanceID": "${instance_id}",
  "instanceState": "CREATED"
}
```

deploy applet with cache and chain client config

```sh
echo '{"cache":{"mode": "MEMORY"}}' | http post :8888/srv-applet-mgr/v0/deploy/applet/$APPLETID -A bearer -a $TOK
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
