# If the first argument is "run"...
# WIP...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif


THIS_FILE := $(lastword $(MAKEFILE_LIST))

git_username="Charles Watkins"
git_email="chris@watkinslabs.com"
 
.DEFAULT: help
.PHONY: build test install-all install-user

help:
	@ echo " make build          | build it"
	@ echo " make test           | test it"
	@ echo ""


build:
	@go build


test: build
	@./term21  ../ttygif-assets/cast/ls2.cast 

# install for all users
install-all:
	@mkdir -p /etc/term21/
	@cp assets/themes/ /etc/term21/ -R
	@cp assets/fonts/ /etc/term21/ -R
	@cp term21 /usr/bin/

# install for local user
install-user:
	@mkdir -p ~/.config/term21/
	@cp assets/themes/ ~/.config/term21/ -R
	@cp assets/fonts/ ~/.config/term21/ -R
	@cp term21 /usr/bin/

