pipeline {
  agent any
  stages {
    stage('checkout') {
      steps {
        sleep 6
        sh '''cd ./workers
ls'''
      }
    }
    stage('install') {
      steps {
        sleep(unit: 'SECONDS', time: 5)
      }
    }
  }
}