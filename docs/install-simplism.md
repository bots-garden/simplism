# Install Simplism

## The simplest way

```bash
SIMPLISM_DISTRO="Linux_arm64" # 👀 https://github.com/bots-garden/simplism/releases
VERSION="0.1.1"
wget https://github.com/bots-garden/simplism/releases/download/v${VERSION}/simplism_${SIMPLISM_DISTRO}.tar.gz -O simplism.tar.gz 
tar -xf simplism.tar.gz -C /usr/bin
rm simplism.tar.gz
simplism version
```

## The Docker way

```bash
docker run --rm k33g/simplism:0.1.1 /simplism version
```
