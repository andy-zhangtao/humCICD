# HICD
My CI/CD framework

## Agents

- **hicd** 接受Github发送的事件通知
- **gitAgent** 接受HICD转发的GitHub通知消息,解析工程的配置信息
- **buildAgent** 接受gitAgent解析后的工程配置信息,根据配置信息选择相对应的语言处理模块
- **goAgent** 构建Golang工程的Agent
- **echoAgent** 接受语言构建模块返回的日志信息,并通过邮件提醒用户

## Images
> 当前所有镜像的tag均为latest

- vikings/hicd
- vikings/gitagent
- vikings/buildagent
- vikings/goagent
- vikings/echoagent

## How to run hicd

- **hicd**
```
docker run \
        --name hicd \
        --log-driver=json-file \
        --net host \
        -e HICD_NSQD_ENDPOINT=127.0.0.1:4150 \
        vikings/hicd
```

- **buildAgent**
```
docker run \
        --name buildagent \
        --log-driver=json-file \
        --net host \
        -e HICD_NSQD_ENDPOINT=127.0.0.1:4150 \
        -v /var/run/docker.sock:/var/run/docker.sock \
        vikings/buildagent
```

- **gitagent**
```
docker run \
        --name gitagent \
        --log-driver=json-file \
        --net host \
        -e HICD_NSQD_ENDPOINT=127.0.0.1:4150 \
        vikings/gitagent
```