
build:
	docker build -t tangsengdaodaoserverdecom .
push:
	docker tag tangsengdaodaoserverdecom 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
	docker push 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
deploy:
	docker build -t tangsengdaodaoserverdecom .
	docker tag tangsengdaodaoserverdecom 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserverdecom:latest
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
