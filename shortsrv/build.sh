#!/bin/bash
cd src
go build -o ../bin/shortsrv *.go 
cd ../deploy
rm -rf repo
mkdir repo
cp ../bin/shortsrv repo/shortsrv
cp ../conf/wshort.conf repo/wshort.conf
docker image rm shortsrv:latest
docker build -t shortsrv:latest . 

REG=ccr.ccs.tencentyun.com

if docker system info >/dev/null 2>&1; then
  if grep -q "$REG" ~/.docker/config.json 2>/dev/null; then
    echo "本机保存过 $REG 的登录信息"
  else
    echo "本机没有 $REG 的登录信息，开始登录"
    docker login $REG --username=100035760184
  fi
else
  echo "Docker 没启动"
  exit 1
fi

docker tag shortsrv:latest ccr.ccs.tencentyun.com/hhlrepo/shortsrv:latest
docker push ccr.ccs.tencentyun.com/hhlrepo/shortsrv:latest


