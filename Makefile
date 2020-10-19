SHELL := /usr/bin/env bash -o pipefail
GOPKG ?= github.com/MrEhbr/gofsm
DOCKER_IMAGE ?=	mrehbr/gofsm
GOBINS ?= .
GO_APP ?= gofsm

include rules.mk
