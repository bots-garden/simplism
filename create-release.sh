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
- update of ./scratch/main.go
- update of README.md
- update of generators/docker.image.txt
- update of generators/simplism.version.txt
- update of k8S/.env
- update of k8S/README.md
- update the kubernetes manifest (docker image)
- update tag in 02-wasm-registry-with-two-pods.md
- update tag in 01-wasm-function-in-a-pod.md
COMMENT

cd ./scratch
tinygo build -scheduler=none --no-debug \
-o ../server/embedded/scratch.wasm \
-target wasi main.go

cd ../

echo "$TAG $MESSAGE"
echo -n $MESSAGE > cmds/version.txt

find . -name '.DS_Store' -type f -delete

git add .
git commit -m "📦 ${MESSAGE}"

git tag -a ${TAG} -m "${MESSAGE}"
git push origin ${TAG}

#goreleaser release --snapshot --clean
goreleaser release --clean

gh release upload ${TAG} ./k8s/manifests/deploy-wasm-from-remote.yaml
gh release upload ${TAG} ./k8s/manifests/deploy-wasm-from-volume.yaml
gh release upload ${TAG} ./k8s/manifests/deploy-wasm-registry.yaml
gh release upload ${TAG} ./k8s/manifests/deploy-wasm-from-registry.yaml
gh release upload ${TAG} ./k8s/manifests/wasm-registry-volume.yaml
gh release upload ${TAG} ./k8s/manifests/wasm-files-volume.yaml


echo "👋 Create the 🐳 image manually with dockerize-simplism.sh"
