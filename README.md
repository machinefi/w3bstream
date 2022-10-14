# w3bstream

## Arch

![w3bstream](__doc__/modules_and_dataflow.png)

## Features

1. wasm applet management
2. wasm runtime instance deployment
3. interact with wasm (a word count demo)

## Run with docker

### build docker image

```bash
make build_image
```

### Run docker container

```bash
 docker-compose -f ./docker-compose.yaml up -d
 ```

### Access Admin Panel

Visit http://localhost:3000 to get started.

The default admin password is `iotex.W3B.admin`

## Run with binary
### Dependencies:

- os : macOS(11.0+)
- docker: to start a postgres
- httpie: a simple curl command
- tinygo: to build wasm code

### Init protocols and database

```sh
make run_depends # start postgres and mqtt
make migrate     # create or update schema
```

### start a server
```sh
make run_server
```
keep the terminal alive, and open a new terminal for the other commands.

### login (fetch auth token)

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

### create your project

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
### build demo wasm scripts

```sh
make wasm_demo ## build to `examples` use to deploy wasm applet
```

### create and deploy applet


upload wasm script

> use `examples/word_count/word_count.wasm` or `examples/log/log.wasm`

```sh
## set env vars
export PROJECTID=${project_id}
export PROJECTNAME=${project_name}
export WASMFILE=exampls/log/log.wasm
http --form post :8888/srv-applet-mgr/v0/applet/$PROJECTID file@$WASMFILE info='{"appletName":"log","strategies":[{"eventType":"ANY","handler":"start"}]}' -A bearer -a $TOK
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

### publish event to server

```sh
echo '{"header":{},"payload":"xxx yyy zzz"}' | http post :8888/srv-applet-mgr/v0/event/$PROJECTNAME
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
