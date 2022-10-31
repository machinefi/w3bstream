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

### Create your project

command

```sh
echo '{"name":"${project_name}"}' | http :8888/srv-applet-mgr/v0/project -A bearer -a $TOK
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

### Build demo wasm scripts

```sh
make wasm_demo ## build to `examples` use to deploy wasm applet
```

### Create and deploy applet

upload wasm script

> use `examples/word_count/word_count.wasm` or `examples/log/log.wasm`

```sh
## set env vars
export PROJECTID=${project_id}
export PROJECTNAME=${project_name}
export WASMFILE=examples/log/log.wasm
http --form post :8888/srv-applet-mgr/v0/applet/$PROJECTID file@$WASMFILE info='{"appletName":"log","strategies":[{"eventType":"DEFAULT","handler":"start"}]}' -A bearer -a $TOK
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

start applet

```sh
export INSTANCEID=${instance_id}
http put :8888/srv-applet-mgr/v0/deploy/$INSTANCEID/START -A bearer -a $TOK
```

### Register publisher

```sh
export PUBNAME=${publisher_name}
export PUBKEY=${publisher_unique_key} # global unique
echo '{"name":"'$PUBNAME'", "key":"'$PUBKEY'"}' | http post :8888/srv-applet-mgr/v0/publisher/$PROJECTID -A bearer -a $TOK
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
echo '{"header":{"event_type":"'$EVENTTYPE'","pub_id":"'$PUBKEY'","pub_time":'`date +%s`',"token":"'$PUBTOKEN'"},"payload":"xxx yyy zzz"}' | http post :8888/srv-applet-mgr/v0/event/$PROJECTNAME
```

output like

```json
[
  {
    "instanceID": "${instance_id}",
    "resultCode": 0
  }
]
```

that means some instance handled this event successfully

### Delete project

Be careful.
It will delete anything in the project, contains applet, publisher, strategy etc.
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
cd build && ./pub_client -c '{"header":{"event_type":'$EVENTTYPE',"pub_id":"'$PUBKEY'","pub_time":'`date +%s`',"token":"'$PUBTOKEN'"},"payload":"xxx yyy zzz"}' -t $PROJECTNAME
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
echo '{"contractLog":{"eventType": "DEFAULT", "chainID": 4690, "contractAddress": "${contractAddress}","blockStart": ${blockStart},"blockEnd": ${blockEnd},"topic0":"${topic0}"}}' | http :8888/srv-applet-mgr/v0/monitor/$PROJECTID -A bearer -a $TOK
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
echo '{"chainTx":{"eventType": "DEFAULT", "chainID": 4690, "txAddress": "${txAddress}"}}' | http :8888/srv-applet-mgr/v0/monitor/$PROJECTID -A bearer -a $TOK
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
echo '{"chainHeight":{"eventType": "DEFAULT", "chainID": 4690, "height": ${height}}}' | http :8888/srv-applet-mgr/v0/monitor/$PROJECTID -A bearer -a $TOK
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
