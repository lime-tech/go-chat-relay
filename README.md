go-chat-relay
--------------------------------------------------------------

## About
...

## Clone
We assume that all your golang root located at `~/go` directory.
```bash
git clone git-server:... go/src/go-chat-relay
```

## How to build
Again, we assume that all your golang root located at `~/go` directory.
If your golang root is in different place, just setup `GOPATH`, `PATH` variables manualy before doing anything specified in this chapter.
This project have a script to setup environment and bootstrap the dependencies, all you need to do is just source it like this and make:
```bash
source go-chat-relay/bootstrap
make -C go-chat-relay
```

## Packaging information
All packaging things should be placed in a subdirectory of `packaging/`.
At this time there are:
- RPM spec

## Daemonizing
All daemonizing things should be placed in a subdirectory of `startup/`.
At this time there are:
- SystemD support(unit+tmpfiles.d)
