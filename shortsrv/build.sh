#!/bin/bash
cd deploy
rm -rf repo
mkdir repo
cp ../bin/main repo/shortsrv
cp ../conf/wshort.conf repo/wshort.conf
docker image rm shortsrv:1.0
docker build -t shortsrv:1.0 . 