#!/bin/bash
set -o allexport; source .github.env; set +o allexport
set -o allexport; source .release.env; set +o allexport
# You need to create a .github.env file with these variables
# - GITHUB_TOKEN

: <<'COMMENT'
Todo:
- update of .release.env:
  - TAG
  - MESSAGE
  - IMAGE_TAG
  - GITPOD_IMAGE_TAG (only if needed)
- update of README.md
- update of generators/docker.image.txt
- update of generators/simplism.version.txt
- update of k8S/.env
- update of k8S/README.md
COMMENT

echo "$TAG $MESSAGE"
echo -n $MESSAGE > cmds/version.txt

find . -name '.DS_Store' -type f -delete

git add .
git commit -m "📦 ${MESSAGE}"

git tag -a ${TAG} -m "${MESSAGE}"
git push origin ${TAG}

#goreleaser release --snapshot --clean
goreleaser release --clean

echo "👋 Create the 🐳 image manually with dockerize-simplism.sh"
