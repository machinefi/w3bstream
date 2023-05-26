# service env manual

## Purpose

The purpose of this document is to help SRE and co-developers quickly track service configuration changes.

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

