# W3bstream

## Overview

W3bStream is a general framework for connecting data generated by devices and machines in the physical world to the blockchain world. In a nutshell, W3bStream uses the IoTeX blockchain to orchestrate a decentralized network of gateways (i.e., W3bStream nodes) that streams encrypted data from IoT devices and machines and generates proofs of real-world facts to different blockchains. An overview graphic of W3bstream is 


![image](https://user-images.githubusercontent.com/448293/196618039-365ab2b7-f50a-49c8-a02d-c28e48acafcb.png)

## Arch

![w3bstream](__doc__/modules_and_dataflow.png)

## 🚀 Why W3bstream?

💨 Accelerated Go-to-Market (GTM) Time: W3bstream streamlines the development process for building decentralized privacy-preserving IoT and machine applications (DePINs), resulting in faster GTM and lower development costs.

💪 Trustworthy Decentralized Architecture: W3bstream's decentralized architecture offers transparent application logic, instilling trust in users and eliminating the need for centralized computational oracles like Google Cloud and AWS.

🤝 Enhanced Composability: W3bstream's permissionless infrastructure can be freely composed with various devices and dApps, fostering collaboration, innovation, and enhanced interoperability.

🔒 Privacy Protection (Ownership): W3bstream supports Zero-Knowledge technologies, providing end-to-end protection for user data and ensuring that privacy is maintained throughout the entire process, unlike centralized computational oracles that may compromise users' data privacy.

## Run W3bstream with prebuilt docker images

### Run W3bstream node with W3bstream Studio
Check it out here [w3bstream-studio](https://github.com/machinefi/w3bstream-studio#run-w3bstream-node-with-prebuilt-docker-images).

### Run W3bstream node without W3bstream Studio

Make a path for w3bstream node. In the path, run the following command

```bash
curl https://raw.githubusercontent.com/machinefi/w3bstream/main/docker-compose.yaml > docker-compose.yaml
```

Edit the config in the `yaml` file if needed. Then run

```bash
docker-compose -p w3bstream -f ./docker-compose.yaml up -d
```

Your node should be up and running. 

Please note: the docker images are hosted at [GitHub Docker Registry](https://github.com/machinefi/w3bstream/pkgs/container/w3bstream)

## Getting started

### Start with W3bstream Studio
If you run W3bstream node with **W3bstream Studio**, You can use Metamask to log in to [localhost:3000](localhost:3000) and create a "Hello World" project.
You can follow the [doc](https://docs.w3bstream.com/get-started/deploying-an-applet)

### Start with admin user
1. Login with admin

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

export token for reuse.

```sh
export TOK=${token}
```

2. Create hello world project with default config

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

3. Create and deploy applet under project created previously

```sh
curl https://raw.githubusercontent.com/machinefi/w3bstream-wasm-golang-sdk/main/examples/wasms/log.wasm -o log.wasm
export WASMFILE=./log.wasm
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

4. Register publisher

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

5. Publish event through http

```sh
export TOPIC=${pub_topic} ## intact project name(required) -- you can get it from t_project.f_name
export PUBTOK=${publisher_token} ## created before(required)
export PAYLOAD=${payload} ## set your payload
http post :8889/srv-applet-mgr/v0/event/$TOPIC --raw=$PAYLOAD -A bearer -a $PUBTOK 
```

### Learn more
Please refer to [HOWTO.md](./HOWTO.md) for more details.

## Documentation

Please visit [https://docs.w3bstream.com/](https://docs.w3bstream.com/).

Interested in contributing to the doc? Please edit on [Github](https://github.com/machinefi/w3bstream-docs-gitbook)

## SDKs

### Client SDKs
- Javascript/Typesript: https://github.com/machinefi/w3bstream-client-js
- Python: https://github.com/machinefi/w3bstream-client-python
- Golang: https://github.com/machinefi/w3bstream-client-go
- Android: https://github.com/machinefi/w3bstream-android-sdk
- iOS: https://github.com/machinefi/w3bstream-ios-sdk
- ESP32: https://github.com/machinefi/w3bstream-client-esp32

### WASM
- Golang: https://github.com/machinefi/w3bstream-wasm-golang-sdk
- AssemblyScript: https://github.com/machinefi/w3bstream-wasm-ts-sdk
- Rust: https://github.com/machinefi/w3bstream-wasm-rust-sdk


## Examples

Learning how to get started with W3bstream? Here is a quick get-start example: https://github.com/machinefi/get-started

More code examples: https://github.com/machinefi/w3bstream-examples

Step-by-step tutorials can be found on dev portal: https://developers.iotex.io/

## Contribution Guide
The community welcomes everyone to contribute, you can find the good first [issue](https://github.com/machinefi/w3bstream/issues) in here if you are new to W3bstream.

## Community

- Developer portal: https://developers.iotex.io/
- Developer Discord (join #w3bstream channel): https://w3bstream.com/discord

## License
[Apache-2.0](LICENSE.md)
