pipeline {
  agent any
  stages {
    stage('checkout') {
      steps {
        sleep 6
        sh '''cd ./workers/service-task-manager

go install'''
        script {
          withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin:${HOME}/go/bin"]) {
            sh "go version"
          }
        }

      }
    }
    stage('install') {
      steps {
        sleep(unit: 'SECONDS', time: 5)
      }
    }
  }
}