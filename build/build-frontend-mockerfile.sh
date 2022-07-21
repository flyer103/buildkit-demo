#!/bin/sh

IMAGE=ccr.ccs.tencentyun.com/flyer103/frontend-mocker

docker build -t ${IMAGE} -f dockerfiles/frontend-mockerfile/Dockerfile .
