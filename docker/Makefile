.ONESHELL: ;             # recipes execute in same shell
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

ARGS = $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
build:
	./scripts/build-image.sh $(ARGS)



.PHONY: build-run
build-run:
	./scripts/run.sh $(ARGS) 1

.PHONY: run
run:
	./scripts/run.sh $(ARGS)



.PHONY: build-debug
build-debug:
	./scripts/debug.sh $(ARGS) 1

.PHONY: debug
debug:
	./scripts/debug.sh $(ARGS)



.PHONY: stop
stop:
	docker-compose stop
