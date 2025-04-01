IMAGE_NAME = nicholasjackson/lamport
TAG = v0.1.1

build_commander:
	docker build -t $(IMAGE_NAME)-commander:$(TAG) -f commander_server.Dockerfile ../

build_general:
	docker build -t $(IMAGE_NAME)-general:$(TAG) -f general_server.Dockerfile ../

build_ui:
	docker build -t $(IMAGE_NAME)-ui:$(TAG) -f ui_server.Dockerfile ../

push_commander:
	docker push $(IMAGE_NAME)-commander:$(TAG)

push_general:
	docker push $(IMAGE_NAME)-general:$(TAG)

push_ui:
	docker push $(IMAGE_NAME)-ui:$(TAG)

push_all: push_commander push_general push_ui