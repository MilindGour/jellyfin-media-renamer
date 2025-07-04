pipeline {
    agent any

    stages {
        stage("Test") {
            steps {
                echo 'Testing project...'
                sh 'make test'
              }
          }
        stage("Clean") {
            steps {
                echo 'Testing project...'
                sh 'make clean'
              }
          }
        stage("Build") {
            steps {
                echo 'Building project...'
                sh 'make build'
              }
          }
        stage("Deploy") {
            steps {
                echo 'Deploying project...'
                sh 'make deploy'
              }
          }
      }
  }
