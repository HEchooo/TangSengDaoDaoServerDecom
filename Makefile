build:
	docker build -t tangsengdaodaoserver .
push:
	docker tag tangsengdaodaoserver 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserver:latest
	docker push 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserver:latest
deploy:
	docker build -t tangsengdaodaoserver .
	docker tag tangsengdaodaoserver 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserver:latest
	docker push 851725583589.dkr.ecr.us-east-1.amazonaws.com/tangsengdaodaoserver:latest
deploy-v1.5:
	docker build -t tangsengdaodaoserver .
	docker tag tangsengdaodaoserver registry.cn-shanghai.aliyuncs.com/wukongim/tangsengdaodaoserver:v1.5
	docker push registry.cn-shanghai.aliyuncs.com/wukongim/tangsengdaodaoserver:v1.5
run-dev:
	docker-compose build;docker-compose up -d
stop-dev:
	docker-compose stop
env-test:
	docker-compose -f ./testenv/docker-compose.yaml up -d 