#!/usr/bin/env bash
#  vim:ts=4:sts=4:sw=4:et
#
#  Author: Hari Sekhon
#  Date: 2016-07-25 18:14:49 +0100 (Mon, 25 Jul 2016)
#
#  https://github.com/harisekhon/bash-tools
#
#  License: see accompanying Hari Sekhon LICENSE file
#
#  If you're using my code you're welcome to connect with me on LinkedIn and optionally send me feedback to help improve or steer this or other code I publish
#
#  https://www.linkedin.com/in/harisekhon
#

set -euo pipefail
[ -n "${DEBUG:-}" ] && set -x
srcdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$srcdir/.."

. bash-tools/utils.sh

section "Dockerfiles checks"

tests/check_repo_names.sh

tests/check_docker-compose_images.sh

tests/pytools_checks.sh

echo

echo "Checking post build hook scripts separately as they're not inferred by .sh extension"
bash-tools/check_shell_syntax.sh */hooks/post_build

bash-tools/all.sh

echo

tests/projects_without_docker-compose_yet.sh

tests/projects_without_README_yet.sh
