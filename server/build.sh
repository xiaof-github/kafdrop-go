#!/bin/sh

git pull

BUILD_DATE=`date +%Y-%m-%d:%H:%M:%S`
echo `date +%Y-%m-%d:%H:%M:%S` > version

base_url=docker.io
docker_url=xiaof-github/kafdrop-go:latest


docker run --rm -v `pwd`:/go/src/github.com/xiaof-github/kafdrop-go -v `pwd`:/go -i golang:stretch \
    /go/src/github.com/xiaof-github/kafdrop-go/install.sh

docker build -t ${docker_url} .

echo "docker login -u  "
echo docker login -u xiaof-github ${base_url}
echo docker push ${base_url}/${docker_url}

echo "build & push finish ."
echo "##############  $BUILD_DATE  ##############"
