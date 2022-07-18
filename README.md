# Overview

Some demos for exploring [buildkitd](https://github.com/moby/buildkit).

# Demo

## hardcode

The demo hardcodes [LLB](https://github.com/moby/buildkit#exploring-llb) to understand how Frontends converts source to LLB.

Get json formatted LLB:

```shell
$ go run cmd/hardcode/main.go | buildctl debug dump-llb | jq .
```

Build:

```shell
$ docker run -d --name buildkitd --privileged moby/buildkit:v0.10.3
$ export BUILDKIT_HOST=docker-container://buildkitd
$ go run cmd/hardcode/main.go | buildctl build --local context=. --output type=tar,dest=out.tar
```

You can unpack `out.tar` and find `README.md` and `built.txt` files which are operated in the program.
You can also modify the output type to build the type that `buildkitd` supports: [Output](https://github.com/moby/buildkit#output).
