# Install Simplism

## The simplest way

```bash
SIMPLISM_DISTRO="Linux_arm64" # 👀 https://github.com/bots-garden/simplism/releases
VERSION="0.0.8"
wget https://github.com/bots-garden/simplism/releases/download/v${VERSION}/simplism_${SIMPLISM_DISTRO}.tar.gz -O simplism.tar.gz 
tar -xf simplism.tar.gz -C /usr/bin
rm simplism.tar.gz
simplism version
```

## The Docker way

```bash
docker run --rm k33g/simplism:0.0.7 /simplism version
```
