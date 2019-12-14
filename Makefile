DATE := `date +%Y%m%d`

release-%: ## release version
	@if [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	bash deploy/release.sh $*; \
	else \
	echo "release failed"; exit 1; \
	fi
	@echo "release ok"

build-images-%: ## build docker images
	if [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t ranklier/ltk:$* -f ./cmd/Dockerfile .; \
	docker build -t ranklier/ltk-mysql-prisma-deploy:$* -f ./prisma/mysql-prisma/Dockerfile ./prisma/mysql-prisma; \
	else \
	echo "build-images failed"; exit 1; \
	fi
	@echo "build images ok"

push-images-%: ## push docker images
	@if [ "$* = "latest" ];then \
	docker push ranklier/ltk:latest; \
	docker push ranklier/ltk-mysql-prisma-deploy:latest; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push ranklier/ltk:$*; \
	docker push ranklier/ltk-mysql-prisma-deploy:$*; \
	fi
	@echo "push images ok"

export-prisma-file: ## export prisma.go from prisma-deploy container
	docker cp ltk-mysql-prisma-deploy:/pkg/prisma/mysql-prisma-client/prisma.go ./pkg/prisma/mysql-prisma-client
	@echo "export prisma file ok"

export-mysql-prisma-file:
	docker cp ltk-mysql-prisma-deploy:/pkg/prisma/mysql-prisma-client/prisma.go ./pkg/prisma/mysql-prisma-client
	@echo "export prisma file ok"

generate:
	cd api/proto && make
	cd pkg/apigateway && make

clean-%: ## clean release package
	@if [ ! -d deploy/ltk-$*-docker-compose-backup ];then \
	mkdir deploy/ltk-$*-docker-compose-backup; \
	else \
	echo "deploy/ltk-$*-docker-compose-backup is exist"; \
	fi
	@if [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	rm -rf deploy/ltk-$*-docker-compose-backup/${DATE}; \
	mkdir deploy/ltk-$*-docker-compose-backup/${DATE}; \
	mv deploy/ltk-$*-docker-compose/* deploy/ltk-$*-docker-compose/.[^.]* deploy/ltk-$*-docker-compose-backup/${DATE}/; \
	rm -rf deploy/ltk-$*-docker-compose; \
	rm -rf deploy/ltk-$*-docker-compose.tar.gz; \
	else \
	echo "clean failed"; exit 1; \
	fi
	@echo "clean ok"

dev-api:
	bee run -main=./cmd/api-gateway/main.go

dev-iam:
	bee run -main=./cmd/iam-engine/main.go

dev-process:
	bee run -main=./cmd/process-manager/main.go
	