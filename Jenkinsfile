pipeline {
  agent any
  stages {
    stage('checkout') {
      steps {
        tool(name: 'go10.3', type: 'go')
      }
    }
    stage('install') {
      steps {
        sh 'go version'
      }
    }
  }
}