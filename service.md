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

1. targets\*: for building all assets(binaries). in this entry, it will traverse all directories in `root/cmd/`, if `Makefile` exists, run `make target`
2. images: for building all images(docker). in this entry, it will traverse all directories in `root/cmd/`, if `Dockerfile` exists, run `make image`
3. test\*: project level testing entry.

Entries in root/cmd/Makefile:

1. target\*: building binary
2. image: build docker image

## Env Format

```
[SERVICENAME]__[GROUPED_CONFIG]_[CONFIG_NAME] = [CONFIG_VALUE]
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

