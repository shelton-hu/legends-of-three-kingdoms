{{VERSION}}

deault: compose-put-global-config

.PHONY: compose-put-global-config
compose-put-global-config: compose-up-etcd compose-up-manager
	@test -s ./global_config.yaml || { echo "[global_config.yaml] not exit"; exit 1; }
	cat ./global_config.yaml | docker-compose exec -T ltk-etcd sh -c "export ETCDCTL_API=3 && etcdctl put /ltk/global_config"

.PHONY: compose-up-etcd
compose-up-etcd: 
	docker-compose up -d ltk-etcd
	@echo "compose-up-etcd ok"

.PHONY: compose-up-manager
compose-up-manager: compose-up-etcd
	docker-compose up -d ltk-db-mysql; sleep 15s
	docker-compose up -d ltk-mysql-prisma; sleep 10s
	docker-compose up -d ltk-mysql-prisma-deploy
	docker wait ltk-mysql-prisma-deploy
	docker-compose up -d ltk-iam-manager
	docker-compose up -d ltk-process-manager
	docker-compose up -d ltk-api-gateway
	@echo "compose-up-manager ok"

pull-manager:
	docker pull dnorange/ltk:$(VERSION)
	docker pull dnorange/ltk-mysql-prisma-deploy:$(VERSION)
	@echo "pull-manager ok"

compose-down:
	docker-compose down
	@echo "compose-down ok"

clean-images:
	docker rmi ranklier/ltk:$(VERSION) ranklier/ltk-mysql-prisma-deploy:$(VERSION)
	@echo "clean-images ok"

put-global-config:
	@test -s ./global_config.yaml || { echo "[global_config.yaml] not exit"; exit 1; }
	cat ./global_config.yaml | docker-compose exec -T ltk-etcd sh -c "export ETCDCTL_API=3 && etcdctl put /ltk/global_config"
