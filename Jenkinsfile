pipeline {
	agent any

	environment {
		PRIVATEKEY=credentials('cd1436f4-0f72-48e6-b7df-43fd0f57990c')
		GOROOT="/usr/local/src/go"
		PATH="$GOROOT/bin:$PATH"
	}

	stages {
		stage('Build') {
			steps {
			echo 'Building..'
			sh 'go mod download'
			sh 'GOARCH=arm GOARM=7 go build -o api main.go'
			}
		}

		stage('Staging') {
			steps {
				script {
					// Copy binary to EC2 using scp
					sh '''
					scp -i $PRIVATEKEY api \
					admin@172.31.28.221:~/go/bin
					'''

					// SSH into ec2 and setup systemd service
					sh '''
					ssh -i $PRIVATEKEY admin@172.31.28.221 \
					"sudo mv ~/go/bin/api /usr/local/bin"
					ssh -i $PRIVATEKEY admin@172.31.28.221 \
					"sudo chmod +x /usr/local/bin/api"
					'''
					// restart the service
					sh '''
					ssh -i $PRIVATEKEY admin@172.31.28.221 \
					"sudo systemctl daemon-reload"
					ssh -i $PRIVATEKEY admin@172.31.28.221 \
					"sudo systemctl restart myapi.service"
					'''
				}
			}
		}
	}
}
