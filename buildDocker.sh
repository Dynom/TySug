#!/bin/sh
LATEST_TAG="$(git describe --tags "$(git rev-list --tags --max-count=1 4b825dc642cb6eb9a060e54bf8d69288fbee4904)")"
TAG_REF="$(git show-ref --hash --tags "${LATEST_TAG}")"

docker build -t "dynom/tysug:${LATEST_TAG}" \
	--build-arg VERSION="${LATEST_TAG}" \
	--build-arg GIT_REF="${TAG_REF}" \
	. &&
docker tag "dynom/tysug:${LATEST_TAG}" dynom/tysug:latest
