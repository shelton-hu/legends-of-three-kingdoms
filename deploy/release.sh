#!/bin/bash

mkdir deploy/ltk-$1-docker-compose
cp -r deploy/docker-compose/. deploy/ltk-$1-docker-compose
echo `./deploy/version.sh ltk-$1` >> deploy/ltk-$1-docker-compose/.env
sed -e 's/ /\'$'\n/g' deploy/ltk-$1-docker-compose/.env >> deploy/ltk-$1-docker-compose/tmp
rm -rf deploy/ltk-$1-docker-compose/.env
mv deploy/ltk-$1-docker-compose/tmp deploy/ltk-$1-docker-compose/.env
sed -e 's/{{}}/\'$'\n/g' deploy/ltk-$1-docker-compose/.env >> deploy/ltk-$1-docker-compose/tmp
rm -rf deploy/ltk-$1-docker-compose/.env
mv deploy/ltk-$1-docker-compose/tmp deploy/ltk-$1-docker-compose/.env
sed -e "s/{{VERSION}}/VERSION := $1/g" deploy/ltk-$1-docker-compose/Makefile >> deploy/ltk-$1-docker-compose/tmp
rm -rf deploy/ltk-$1-docker-compose/Makefile
mv deploy/ltk-$1-docker-compose/tmp deploy/ltk-$1-docker-compose/Makefile
cp -r deploy/config/global_config.init.yaml deploy/ltk-$1-docker-compose/global_config.yaml
cd deploy/ && tar -czvf ltk-$1-docker-compose.tar.gz ltk-$1-docker-compose