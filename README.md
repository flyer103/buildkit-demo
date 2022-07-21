# Overview

Some demos for exploring [buildkitd](https://github.com/moby/buildkit).

# Demo

## Dependencies

- Golang: https://go.dev/dl/
- Docker: https://docs.docker.com/
- Image: https://hub.docker.com/r/moby/buildkit
- CLI:
  - buildctl: https://github.com/moby/buildkit
  - jq: https://stedolan.github.io/jq/

## hardcode-image

The demo hardcodes [LLB](https://github.com/moby/buildkit#exploring-llb) to understand how Frontends converts image source to LLB.

Get json formatted LLB:

```shell
$ go run cmd/hardcode-image/main.go | buildctl debug dump-llb | jq .
```

Build:

```shell
$ docker run -d --name buildkitd --privileged moby/buildkit:v0.10.3
$ export BUILDKIT_HOST=docker-container://buildkitd
$ go run cmd/hardcode-image/main.go | buildctl build --local context=. --output type=tar,dest=out.tar
```

You can unpack `out.tar` and find `README.md` and `built.txt` files which are operated in the program.

You can also modify the output type to build the type that `buildkitd` supports: [Output](https://github.com/moby/buildkit#output).

## hardcode-scratch

The demo hardcodes [LLB](https://github.com/moby/buildkit#exploring-llb) to understand how Frontends converts scratch source to LLB.

Get json formatted LLB:

```shell
$ go run cmd/hardcode-scratch/main.go | buildctl debug dump-llb | jq .
```

Build:

```shell
$ docker run -d --name buildkitd --privileged moby/buildkit:v0.10.3
$ export BUILDKIT_HOST=docker-container://buildkitd
$ go run cmd/hardcode-scratch/main.go | buildctl build --local context=. --output type=tar,dest=out.tar
```

You can unpack `out.tar` and find copied files which are operated in the program.

You can also get a feel for buildkitd's parallel builds and caching by looking at the execution.


Note:
- I don't understand the internals yet -_-

## frontend-mockerfile

This demo uses other people's examples to see how buildkitd frontend is implemented: [detail](./cmd/frontend-mockerfile/README.md).

Get json formatted LLB:

```shell
$ go run cmd/frontend-mockerfile/main.go -filename cmd/frontend-mockerfile/Mockerfile.yaml -graph | buildctl debug dump-llb | jq .
```

Build:

```shell
# You should change image registry in the following script.
$ ./build/build-frontend-mockerfile.sh

# PUSH IMAGE THAT IS BUILT ABOVE.
$ docker push <IMAGE/BUILT/ABOVE>

# Image Example: `ccr.ccs.tencentyun.com/flyer103/frontend-mocker:latest`
$ buildctl build --frontend=gateway.v0 --opt source=<IMAGE/BUILT/ABOVE> --local context=. --local dockerfile=./yaml/frontend-mockerfile/ --output type=oci,dest=output.tar
```
