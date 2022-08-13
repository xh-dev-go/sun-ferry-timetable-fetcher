def project_version
def project_branchName

def project_docker_name='sun-ferry-timetable-fetcher'
def setStatus(status){
    def m = '{"state": "'+status+'","context": "continuous-integration/jenkins", "description": "Jenkins", "target_url": "https://jks.xh-network.xyz/job/webclip2-server/'+env.BUILD_NUMBER+'/console"}'
    m = m.replaceAll("\"", "\\\\\"")
    msg = "curl -i -s \"https://api.GitHub.com/repos/$GITHUB_CREDENTIALS_USR/webclip2-server/statuses/$GIT_COMMIT\" -H \"Authorization: token $GITHUB_CREDENTIALS_PSW\" -H \"Content-Type: application/json\" -X POST -d \"$m\""
    //echo "$msg"
    sh msg
}

pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS=credentials('dockerhub-id')
        GITHUB_CREDENTIALS=credentials('github-id')
        DEPLOY_CREDENTIALS=credentials('ssh-deployment')
    }

    stages {
        stage('Init') {
            environment {
                C_VERSION = sh(returnStdout: true, script: 'cat version.txt').trim()
                branchName= sh (returnStdout: true, script: 'echo $GIT_BRANCH').trim()
                commitId= sh (returnStdout: true, script: 'echo $GIT_COMMIT').trim()
            }
            steps {
                script {
                    sh 'printenv'
                    setStatus("pending")
                    project_version = sh(returnStdout: true, script: 'echo $C_VERSION')
                    project_branchName = env.branchName
                }
            }
        }
        stage('Build') {
            environment {
                branchName= sh (returnStdout: true, script: 'echo $GIT_BRANCH').trim()
                commitId= sh (returnStdout: true, script: 'echo $GIT_COMMIT').trim()
                projectName= sh(returnStdout: true, script: "echo $project_docker_name").trim()
                uiBranchId= sh(returnStdout: true, script: "pushd web && echo $(git log --pretty=format:'%h' -n 1) && popd").trim()
            }
            when { expression { return env.GIT_BRANCH == 'origin/master'}}
            steps {
                sh 'printenv'
                sh 'docker build --build-arg branchName=$GIT_BRANCH --build-arg commitId=$GIT_COMMIT -t $projectName:latest .'
                echo 'build complete'
            }
        }
        stage('Publish') {
            environment {
                projectName= sh(returnStdout: true, script: "echo $project_docker_name").trim()
            }
            when { expression { return env.GIT_BRANCH == 'origin/master'}}
            steps {
                sh 'echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin'
                sh 'docker tag $projectName:latest $DOCKERHUB_CREDENTIALS_USR/$projectName:latest'
                sh 'docker push $DOCKERHUB_CREDENTIALS_USR/$projectName:latest'
                sh 'docker tag $DOCKERHUB_CREDENTIALS_USR/$projectName:latest $DOCKERHUB_CREDENTIALS_USR/$projectName:'+project_version
                sh 'docker push $DOCKERHUB_CREDENTIALS_USR/$projectName:'+project_version
            }
        }
//         stage('Deploy') {
//             when { expression { return env.GIT_BRANCH == 'origin/master'}}
//             steps {
//                 withCredentials([string(credentialsId: 'deployment-host', variable: 'host')]) {
//                     sshagent(credentials: ['ssh-deployment']){
//                         sh "ssh -i /ssh/id_rsa $DEPLOY_CREDENTIALS_USR@$host /data/$projectName/refresh.sh"
//                     }
//                 }
//             }
//         }
    }
    post {
        always {
            sh 'docker logout'
        }
        success {
            setStatus("success")
        }
        failure {
            setStatus("failure")
        }
    }
}