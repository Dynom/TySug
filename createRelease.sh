#!/usr/bin/env bash

set -o pipefail -o nounset -o errexit -o errtrace

# Find our .git directory
ROOT_DIR="$(pwd)"
while [ ! -d "${ROOT_DIR}/.git" ]; do

    ROOT_DIR="$(dirname "${ROOT_DIR}")"
    if [[ "x${ROOT_DIR}" == "x/" ]]; then
        echo "Cannot find .git directory, I use that as reference for the commands."
        exit 1
    fi
done

# Determine our project name
NAME="$(basename "$(pwd)")"

# Checking if we have any tags to start with, the cid is Git's magical initial repo hash
TAGS=$(git rev-list --tags --count 4b825dc642cb6eb9a060e54bf8d69288fbee4904)
if [[ "${TAGS}" -eq 0 ]];
then
	echo "No tags detected for ${ROOT_DIR}, please create a tag first!"
	exit 1;
fi

# Figuring out what tag's we're on
LATEST_TAG=$(git describe --tags "$(git rev-list --tags --max-count=1 4b825dc642cb6eb9a060e54bf8d69288fbee4904)")
PREV_TAG=$(git tag --sort version:refname | tail -2 | head -1 || true)

if [[ "x${LATEST_TAG}" == "x" && "x${PREV_TAG}" == "x" ]];
then
    echo "No tag has been found?"
    exit 1
fi
echo "Previous tag is: ${PREV_TAG}"
echo "Building a release for tag: ${LATEST_TAG}"

# Falling back to the first commit, if we only have one tag
if [[ "x${PREV_TAG}" == "x${LATEST_TAG}" ]];
then
    PREV_TAG=$(git rev-list --max-parents=0 HEAD)
fi

# Dependencies
go get github.com/c4milo/github-release
go get github.com/mitchellh/gox

# Cleanup
rm -rf build dist && mkdir -p build dist

# Build
gox -ldflags "-s -w -X main.Version=${LATEST_TAG}" \
    -osarch="darwin/amd64" \
    -osarch="linux/amd64" \
    -osarch="windows/amd64" \
    -rebuild \
    -output "build/{{.Dir}}-${LATEST_TAG}-{{.OS}}-{{.Arch}}/${NAME}" \
	./cmd/web

# Archive
HERE="$(pwd)"
BUILD_DIR="${HERE}/build"
for DIR in "${BUILD_DIR}"/*;
do
    BASE="$(basename "${DIR}")"
    OUT_DIR="${HERE}/dist"
    OUT_FILE_NAME="${BASE}.tar.gz"
    OUT_FILE="${OUT_DIR}/${OUT_FILE_NAME}"
    cd "${DIR}" && \
        tar -czf "${OUT_FILE}" ./* && \
    cd "${OUT_DIR}" && \
        shasum -a 512 "${OUT_FILE_NAME}" > "${OUT_FILE}".sha512
done
cd "${HERE}"

# Building the changelog
DIFF_REF="${PREV_TAG}..${LATEST_TAG}"
CHANGELOG="$(printf '# %s\n%s' 'Changelog' "$(git log "${DIFF_REF}" --oneline --no-merges --reverse)")"

echo "Building the changelog based on these two ref's: '${DIFF_REF}'"
ghr -owner "${GITHUB_USERNAME:-Dynom}" \
    -repository "${GITHUB_REPOSITORY:${NAME}}" \
    -commitish "$(git rev-parse HEAD)" \
    -delete \
    -body "${CHANGELOG}" \
    "${LATEST_TAG}" \
    ./dist/


