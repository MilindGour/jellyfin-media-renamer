pipeline {
    agent any

    stages {
        stage("Test") {
            steps {
                echo 'Testing project...'
                make test
              }
          }
        stage("Clean") {
            steps {
                echo 'Testing project...'
                make clean
              }
          }
        stage("Build") {
            steps {
                echo 'Building project...'
                make build
              }
          }
        stage("Deploy") {
            steps {
                echo 'Deploying project...'
                make deploy
              }
          }
      }
  }
