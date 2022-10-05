pipeline {

    agent { label 'master' }

    tools {
        maven 'mvn362'
    }

    options {
        timestamps()
    }

    stages {
        stage ('Clone sources') {
            steps {
                checkout scm
            }
        }

        stage('Build') {
            steps {
                withCredentials([file(credentialsId: 'livescore-dev-project-container-registry', variable: 'GC_KEY')]) {
                    script {                      
                        env.COMPONENT       = 'config-api'
                        env.VERSION         = sh(script:'cat VERSION', returnStdout: true).trim()
                        env.BRANCH          = sh(script:'echo $GIT_BRANCH | sed -e "s/.*\\///g" -e "s/-/_/g" -e "s/main/master/g"', returnStdout: true).trim()
                        env.FULL_VERSION    = env.VERSION + '-' + env.BRANCH + '-' + BUILD_NUMBER
                        env.PROJECT_ID      = sh(script:'cat "$GC_KEY" | jq -r ".project_id"', returnStdout: true).trim()
                        
                        echo "\n\nfull ver tag: $FULL_VERSION\n\n"
                        sh 'sudo su'
                        echo "\n\nactivate-service-account: $PROJECT_ID\n\n"
                        sh '''
                            gcloud auth activate-service-account --key-file="${GC_KEY}" --project=$PROJECT_ID
                            gcloud auth configure-docker
                        '''
                        echo "\n\ndocker build ls-$COMPONENT:$FULL_VERSION\n\n"
                        sh '''
                            #!/bin/sh -e
                            DOCKER_BUILDKIT=1 docker build . \
                                -t ls-$COMPONENT:$FULL_VERSION
                        '''
                        echo "\n\ndocker push: ls-$COMPONENT:$FULL_VERSION to: $PROJECT_ID\n\n"
                        sh '''
                            echo "ls-$COMPONENT:$FULL_VERSION"
                            docker tag ls-$COMPONENT:$FULL_VERSION eu.gcr.io/$PROJECT_ID/ls-$COMPONENT:$FULL_VERSION
                            docker push eu.gcr.io/$PROJECT_ID/ls-$COMPONENT:$FULL_VERSION
                        '''
                        echo "\n\ndocker rmi images\n\n"
                        sh '''
                            docker rmi ls-$COMPONENT:$FULL_VERSION -f
                            docker rmi eu.gcr.io/$PROJECT_ID/ls-$COMPONENT:$FULL_VERSION -f
                            #docker images | grep \"none\" | awk \'{print $3}\' | xargs docker rmi -f
                            #a=`docker images | grep "ls-$COMPONENT" | awk \'{print $3}\'`; if [ `echo $a | wc -l` -gt 1 ]; then echo xargs docker rmi -f; fi
                        '''
                    }
                }
                buildName FULL_VERSION //readMavenPom().getVersion()
            }
        }
    }
}