# Lenkins

Lenkins is a lightweight component tool. It can easily execute scripts, component applications, and execute commands remotely through a yaml configuration file; it supports git plug-ins, sh plug-ins (local execution commands), CMD plug-ins (remote execution commands), SCP plug-ins (upload or download) etc.

## Installation guide

1. Install Lenkins using the go command

```shell
~ go install github.com/GuoxinL/lenkins@latest
go: downloading github.com/GuoxinL/lenkins v0.0.0-20230617174852-c395440cf0a0
```

2. Run Lenkins!

```shell
~ lenkins -h
Lenkins is a lightweight deployment tool. Lenkins can automatically execute scripts, deploy applications, and remotely execute commands through a configuration file; it supports git plug-ins, sh plug-ins (local execution commands), cmd plug-ins (remote execution commands), scp plugins (upload or download) etc.

Usage:
  lenkins [flags]

Flags:
  -c, --conf string   Deployment configuration file (required)
  -h, --help          help for lenkins

```

If execution fails.  

```shell
zsh: command not found: xxx

```

Please configure GOPATH in the environment variables.  

```shell
echo '
export PATH=$PATH:$GOPATH
' >> /path/to/profile
```

## Hello world !!!

```shell
lenkins -c https://raw.githubusercontent.com/GuoxinL/lenkins/master/example/primary.yaml
```