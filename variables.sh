#!/bin/bash

export COMPILE_DATE = $(shell date "+%F-%T")
export COMMIT_HASH = $(shell git rev-parse --short HEAD)