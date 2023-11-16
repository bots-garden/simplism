#FROM gitpod/workspace-full
FROM gitpod/workspace-base


LABEL maintainer="@k33g_org"

ARG GO_ARCH=amd64
ARG GO_VERSION=1.21.3
ARG TINYGO_ARCH=amd64
ARG TINYGO_VERSION=0.30.0
ARG EXTISM_ARCH=amd64
ARG EXTISM_VERSION=0.3.0

ARG NODE_VERSION=v21.1.0
ARG NODE_DISTRO=linux-x64

USER root

ARG DEBIAN_FRONTEND=noninteractive

# Update the system and install necessary tools
RUN <<EOF
apt-get update 
apt-get install -y curl wget git build-essential xz-utils bat exa software-properties-common htop openssh-server
ln -s /usr/bin/batcat /usr/bin/bat
apt-get -y install hey
EOF

USER gitpod

# ------------------------------------
# Install Go
# ------------------------------------
RUN <<EOF
GO_ARCH="${GO_ARCH}"
GO_VERSION="${GO_VERSION}"

wget https://golang.org/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
tar -xvf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
rm go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
sudo mv go /usr/local
EOF

# ------------------------------------
# Set Environment Variables for Go
# ------------------------------------
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN <<EOF
go version
go install -v golang.org/x/tools/gopls@latest
go install -v github.com/ramya-rao-a/go-outline@latest
go install -v github.com/stamblerre/gocode@v1.0.0
go install -v github.com/mgechev/revive@v1.3.2
EOF

# ------------------------------------
# Install TinyGo
# ------------------------------------
RUN <<EOF
TINYGO_ARCH="${TINYGO_ARCH}"
TINYGO_VERSION="${TINYGO_VERSION}"

wget https://github.com/tinygo-org/tinygo/releases/download/v${TINYGO_VERSION}/tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb
sudo dpkg -i tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb
rm tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb
EOF

# ------------------------------------
# Install Wasmtime, Wazero, Wasmer CLI
# ------------------------------------
RUN <<EOF
curl https://wasmtime.dev/install.sh -sSf | bash

curl https://wazero.io/install.sh | sh
sudo mv ./bin/wazero /usr/bin/wazero

curl -sSf https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash
EOF

# ------------------------------------
# Install Extism CLI
# ------------------------------------
RUN <<EOF
EXTISM_ARCH="${EXTISM_ARCH}"
EXTISM_VERSION="${EXTISM_VERSION}"

wget https://github.com/extism/cli/releases/download/v${EXTISM_VERSION}/extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz

sudo tar -xf extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz -C /usr/bin
rm extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz

extism --version
EOF

# ------------------------------------
# Install Rust + Wasm Toolchain
# ------------------------------------
RUN <<EOF
sudo apt install -y pkg-config libssl-dev
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
export RUSTUP_HOME=~/.rustup
export CARGO_HOME=~/.cargo
export PATH=$PATH:$CARGO_HOME/bin
curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh 
rustup target add wasm32-wasi
EOF

# ------------------------------------
# Install NodeJS
# ------------------------------------
RUN <<EOF
NODE_VERSION="${NODE_VERSION}"
NODE_DISTRO="${NODE_DISTRO}"
wget https://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-${NODE_DISTRO}.tar.xz
sudo mkdir -p /usr/local/lib/nodejs
sudo tar -xJvf node-$NODE_VERSION-$NODE_DISTRO.tar.xz -C /usr/local/lib/nodejs
rm node-$NODE_VERSION-$NODE_DISTRO.tar.xz
EOF

ENV VERSION="${NODE_VERSION}"
ENV DISTRO="${NODE_DISTRO}"
ENV PATH=/usr/local/lib/nodejs/node-$VERSION-$DISTRO/bin:$PATH

RUN <<EOF
# -----------------------
# Install Extism JS PDK
# -----------------------
#curl -O https://raw.githubusercontent.com/extism/js-pdk/main/install.sh
#sh install.sh
export TAG="v0.5.0"
#export ARCH="aarch64"
export ARCH="x86_64"
export  OS="linux"
curl -L -O "https://github.com/extism/js-pdk/releases/download/$TAG/extism-js-$ARCH-$OS-$TAG.gz"
gunzip extism-js*.gz
chmod +x extism-js-*
sudo mv extism-js-* /usr/local/bin/extism-js
EOF

# ------------------------------------
# Install Spin
# ------------------------------------
RUN <<EOF
mkdir spin-framework
cd spin-framework
curl -fsSL https://developer.fermyon.com/downloads/install.sh | bash
sudo mv spin /usr/local/bin/
cd ..
rm -rf spin-framework
spin templates install --git https://github.com/fermyon/spin --upgrade

spin templates list
spin plugins update
EOF

# ------------------------------------
# Install Wasm Workers Server
# ------------------------------------
RUN <<EOF
curl -fsSL https://workers.wasmlabs.dev/install -o install.sh
chmod +x install.sh
./install.sh --local
sudo mv wws /usr/bin
rm install.sh
EOF

# ------------------------------------
# Install GoReleaser
# ------------------------------------
RUN <<EOF
curl -sfL https://goreleaser.com/static/run | bash
EOF


# Command to run when starting the container
CMD ["/bin/bash"]
