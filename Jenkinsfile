pipeline {
  agent any
  stages {
    stage('checkout') {
      steps {
        sh '''#!/usr/bin/bash

WORKROOT=$(pwd)
cd ${WORKROOT}

# unzip go environment
go_env="go1.10.3.linux-amd64.tar.gz"
wget -c https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
tar -zxf $go_env
if [ $? -ne 0 ];
then
    echo "fail in extract go"
    exit 1
fi
echo "OK for extract go"
rm -rf $go_env

# prepare PATH, GOROOT and GOPATH
export PATH=$(pwd)/go/bin:$PATH
export GOROOT=$(pwd)/go
export GOPATH=$(pwd)

# build
go version
if [ $? -ne 0 ];
then
    echo "fail to go build"
    exit 1
fi
echo "OK for go build"'''
      }
    }
    stage('install') {
      steps {
        sleep(unit: 'SECONDS', time: 5)
      }
    }
  }
}