pipeline {
  agent any
  stages {
    stage('checkout') {
      steps {
        tool(name: 'go10.3', type: 'go')
        sh 'ls'
      }
    }
    stage('install') {
      steps {
        script {
          def root = tool name: 'go10.3', type: 'go'
          withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"


            stage 'preTest'
            sh 'go version'
          }
        }

      }
    }
  }
}