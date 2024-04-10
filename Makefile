USER_BINARY=userApp
BROKER_BINARY=brokerApp


## up_build: stops docker-compose if running, build all projects and start docker-compose
# up_build: build_user build_broker
# 	@echo "stops all container if running"
# 	docker-compose down
# 	@echo "building (when required) and start docker images"
# 	docker-compose up --build
# 	@echo "docker images built and started"

# ## down: stop docker-compose
# down:
# 	@echo "stoping docker-compose"
# 	docker-compose down
# 	@echo "done"

# ## build_user: build the user binary as linux executable file
# build_user:
# 	@echo "building user binary"
# 	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd
# 	@echo "done"

# build_broker:
# 	@echo "building user binary"
# 	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd
# 	@echo "done"

user_repo_test:
	@echo "running test for user repository"
	cd internal/user/repositories && go test . --coverprofile=cover.out
	@echo "finished running all test"

user_service_test:
	@echo "running test for user service"
	cd internal/user/services && go test . --coverprofile=cover.out
	@echo "finished running all test"

user_controller_test:
	@echo "running test for user controller"
	cd internal/user/controllers && go test . --coverprofile=cover.out
	@echo "finished running all test"

show_user_repo:
	cd internal/user/repositories && go tool cover -html=cover.out

show_user_service:
	cd internal/user/services && go tool cover -html=cover.out

show_user_controller:
	cd internal/user/controllers && go tool cover -html=cover.out
