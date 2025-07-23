
build:
	docker build -t tangsengdaodaoserverdecom .
push:
	docker tag tangsengdaodaoserverdecom 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
	aws ecr get-login-password --region us-east-1 | sudo docker login --username AWS --password-stdin 851725583589.dkr.ecr.us-east-1.amazonaws.com
	docker push 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
deploy:
	docker build -t tangsengdaodaoserverdecom .
	docker tag tangsengdaodaoserverdecom 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
	aws ecr get-login-password --region us-east-1 | sudo docker login --username AWS --password-stdin 851725583589.dkr.ecr.us-east-1.amazonaws.com
	docker push 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
deploy-v1.5:
	docker build -t tangsengdaodaoserverdecom .
	docker tag tangsengdaodaoserverdecom registry.cn-shanghai.aliyuncs.com/wukongim/tangsengdaodaoserverdecom:v1.5
	docker push registry.cn-shanghai.aliyuncs.com/wukongim/tangsengdaodaoserverdecom:v1.5
run-dev:
	docker-compose build;docker-compose up -d
stop-dev:
	docker-compose stop
env-test:
	docker-compose -f ./testenv/docker-compose.yaml up -d

# 部署相关命令
deploy-test:
	./deploy.sh test
deploy-prod:
	./deploy.sh prod
deploy-test-build:
	./deploy.sh test -b
deploy-prod-build:
	./deploy.sh prod -b
logs-test:
	./deploy.sh test -l
logs-prod:
	./deploy.sh prod -l
status-test:
	./deploy.sh test -s
status-prod:
	./deploy.sh prod -s
stop-test:
	./deploy.sh test -d
stop-prod:
	./deploy.sh prod -d
backup:
	./deploy.sh test --backup
clean:
	./deploy.sh --clean 
