#!/bin/bash

#docker entrypint
version="$(cat VERSION)"
IRIS_VERSION=$version iris "$@"