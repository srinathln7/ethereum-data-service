# Paths and URLs for documentation
DOC_PATH=${PWD}/docs/html
DOC_URL=http://localhost:6060/pkg/ethereum-data-service/?m=all
CLIENT_URL=http://localhost:6060/pkg/ethereum-data-service/internal/client/?m=all
CONFIG_URL=http://localhost:6060/pkg/ethereum-data-service/internal/config/?m=all
MODEL_URL=http://localhost:6060/pkg/ethereum-data-service/internal/model/?m=all
SERVICES_URL=http://localhost:6060/pkg/ethereum-data-service/internal/services/?m=all
STORAGE_URL=http://localhost:6060/pkg/ethereum-data-service/internal/storage/?m=all

# URLs for internal services
PUB_URL=http://localhost:6060/pkg/ethereum-data-service/internal/services/pub/?m=all
SUB_URL=http://localhost:6060/pkg/ethereum-data-service/internal/services/sub/?m=all
BOOTSTRAP_URL=http://localhost:6060/pkg/ethereum-data-service/internal/services/bootstrapper/?m=all

# Build documentation
.PHONY: docs
docs:
	rm -rf ${DOC_PATH}
	mkdir -p ${DOC_PATH}

	# Build the docs for the main package
	godoc -url ${DOC_URL} > ${DOC_PATH}/index.html

	# Build the docs for internal packages
	godoc -url ${CLIENT_URL} > ${DOC_PATH}/client.html
	godoc -url ${CONFIG_URL} > ${DOC_PATH}/config.html
	godoc -url ${MODEL_URL} > ${DOC_PATH}/model.html
	godoc -url ${STORAGE_URL} > ${DOC_PATH}/storage.html

	# Build the docs for internal services
	godoc -url ${PUB_URL} > ${DOC_PATH}/pub.html
	godoc -url ${SUB_URL} > ${DOC_PATH}/sub.html
	godoc -url ${BOOTSTRAP_URL} > ${DOC_PATH}/bootstrap.html

.PHONY: buildup

buildup:
	docker-compose -f docker-compose.yml up --build

.PHONY: builddown

builddown:
	docker-compose -f docker-compose.yml down
