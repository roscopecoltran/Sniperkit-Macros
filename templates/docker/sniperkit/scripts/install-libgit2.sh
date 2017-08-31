#!/bin/sh
set -x
set -e

# Set temp environment vars
export LIBGIT2_VCS_REPO_URL=https://github.com/libgit2/libgit2.git
export LIBGIT2_VCS_BRANCH=v0.26.0
export LIBGIT2_VCS_CLONE_DEPTH=1
export LIBGIT2_VCS_CLONE_PATH=/tmp/libgit2

# Compile & Install libgit2 (v0.23)
git clone -b ${LIBGIT2_VCS_BRANCH} --depth ${LIBGIT2_CLONE_DEPTH} -- ${LIBGIT2_VCS_REPO_URL} ${LIBGIT2_VCS_CLONE_PATH}

mkdir -p ${LIBGIT2_VCS_CLONE_PATH}/build
cd ${LIBGIT2_VCS_CLONE_PATH}/build
cmake .. -DBUILD_CLAR=off
cmake --build . --target install

# Cleanup
rm -r ${LIBGIT2_VCS_CLONE_PATH}