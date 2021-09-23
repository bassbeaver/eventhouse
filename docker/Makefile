.ONESHELL: ;             # recipes execute in same shell
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

ARGS = $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
build:
	./scripts/build-image.sh $(ARGS)


.PHONY: run
run:
	./scripts/run.sh $(ARGS)


.PHONY: debug
debug:
	./scripts/debug.sh $(ARGS)


.PHONY: stop
stop:
	docker-compose stop
