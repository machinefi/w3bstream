# W3bstream

## Overview

W3bStream is a general framework for connecting data generated by devices and machines in the physical world to the blockchain world. In a nutshell, W3bStream uses the IoTeX blockchain to orchestrate a decentralized network of gateways (i.e., W3bStream nodes) that streams [encrypted] data from IoT devices and machines and generates proofs of real-world facts to different blockchains.

![image](https://user-images.githubusercontent.com/448293/196618039-365ab2b7-f50a-49c8-a02d-c28e48acafcb.png)


## Arch

![w3bstream](__doc__/modules_and_dataflow.png)

## Run with prebuilt docker

```
export WS_WORKING_DIR=$PWD/build_image
docker-compose -p w3bstream -f ./docker-compose.yaml up -d
```

`WS_WORKING_DIR` is the working directory for w3bstream node.

You will run with prebuilt docker image from recent stable versions.

## Run with docker

### init frontend

```bash
make init_frontend
```

### Update frontend to latest if needed

```bash
make update_frontend
```

### Build docker image

```bash
make build_image
```

### Run docker container

```bash
 make run_image
 ```

 ### drop docker image
 ```bash
 make drop_image
 ```

## Access W3bstream Studio

Visit http://localhost:3000 to get started.

The default admin password is `iotex.W3B.admin`

## Run with binary

### Dependencies:

- OS : macOS(11.0+) / Linux (tested on Ubuntu 16+)
- Go: golang (1.18+)
- Docker: to start a postgres
- Tinygo: to build wasm code
- make: run makefile
- GCC: 11.3.0
- protobuf: 3.12+
- Httpie: (optional) a simple curl command (used to interact with W3bstream node via cli)

### Init protocols and database

```sh
make run_depends # start postgres and mqtt
make migrate     # create or update schema
```

### Start a server

```sh
make run_server
```

keep the terminal alive, and open a new terminal for the other commands.

### Interact with W3bstream using CLI

Please refer to [HOWTO.md](./HOWTO.md) for more details.
