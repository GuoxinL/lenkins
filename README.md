# Lenkins

Lenkins是一个轻量级部署工具。它可以很方便的通过一个yaml配置文件自动执行脚本，部署应用程序，远程执行命令;它支持git插件，sh插件(本地执行命令)，CMD插件(远程执行命令)，SCP插件(上传或下载)等。

## 安装指南

1. 使用go命令安装lenkins

```shell
~ go install github.com/GuoxinL/lenkins@latest
go: downloading github.com/GuoxinL/lenkins v0.0.0-20230617174852-c395440cf0a0
```

2. 运行Lenkins!

```shell
~ lenkins -h
Lenkins is a lightweight deployment tool. Lenkins can automatically execute scripts, deploy applications, and remotely execute commands through a configuration file;   it supports git plug-ins, sh plug-ins (local execution commands), cmd plug-ins (remote execution commands), scp plugins (upload or download) etc.

Usage:
lenkins [flags]

Flags:
-c, --conf string   Deployment configuration file (required)
-h, --help          help for lenkins

```

如果运行失败

```shell
zsh: command not found: xxx

```

请在环境变量中配置GOPATH

```shell
echo '
export PATH=$PATH:$GOPATH
' >> /path/to/profile
```

## Hello world !!!

```shell
lenkins -c https://gitee.com/guoxinliu/lenkins/raw/master/example/primary.yaml
```