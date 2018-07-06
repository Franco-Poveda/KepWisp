pipeline {
  agent any
  stages {
    stage('checkout') {
      steps {
        tool(name: 'go10.3', type: 'go')
        sh '''#!/usr/bin/bash

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