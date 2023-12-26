FROM ubuntu:22.04

LABEL maintainer="@k33g_org"

ARG GO_ARCH=${GO_ARCH}
ARG GO_VERSION=${GO_VERSION}
ARG TINYGO_ARCH=${TINYGO_ARCH}
ARG TINYGO_VERSION=${TINYGO_VERSION}
ARG EXTISM_ARCH=${EXTISM_ARCH}
ARG EXTISM_VERSION=${EXTISM_VERSION}

ARG NODE_DISTRO=${NODE_DISTRO}
ARG NODE_VERSION=${NODE_VERSION}

ARG USER_NAME=${USER_NAME}

ARG DEBIAN_FRONTEND=noninteractive

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8
ENV LC_COLLATE=C
ENV LC_CTYPE=en_US.UTF-8

# Update the system and install necessary tools
RUN <<EOF
apt-get update 
apt-get install -y curl wget git build-essential xz-utils bat exa software-properties-common htop openssh-server pkg-config libssl-dev sudo
ln -s /usr/bin/batcat /usr/bin/bat
apt-get -y install hey
EOF

# ------------------------------------
# Install Go
# ------------------------------------
RUN <<EOF
GO_ARCH="${GO_ARCH}"
GO_VERSION="${GO_VERSION}"

wget https://golang.org/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
tar -xvf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
mv go /usr/local
EOF

# ------------------------------------
# Set Environment Variables for Go
# ------------------------------------
#ENV GOROOT=/usr/local/go
#ENV GOPATH=$HOME/go
#ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/home/${USER_NAME}/go"
ENV GOROOT="/usr/local/go"

RUN <<EOF
go version
go install -v golang.org/x/tools/gopls@latest
go install -v github.com/ramya-rao-a/go-outline@latest
go install -v github.com/stamblerre/gocode@v1.0.0
EOF

# ------------------------------------
# Install TinyGo
# ------------------------------------
RUN <<EOF
TINYGO_ARCH="${TINYGO_ARCH}"
TINYGO_VERSION="${TINYGO_VERSION}"

wget https://github.com/tinygo-org/tinygo/releases/download/v${TINYGO_VERSION}/tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb
dpkg -i tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb
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

tar -xf extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz -C /usr/bin
rm extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz

extism --version
EOF


# ------------------------------------
# Install NodeJS
# ------------------------------------
RUN <<EOF
NODE_VERSION="${NODE_VERSION}"
NODE_DISTRO="${NODE_DISTRO}"
wget https://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-${NODE_DISTRO}.tar.xz
mkdir -p /usr/local/lib/nodejs
tar -xJvf node-$NODE_VERSION-$NODE_DISTRO.tar.xz -C /usr/local/lib/nodejs
rm node-$NODE_VERSION-$NODE_DISTRO.tar.xz
EOF

ENV VERSION="${NODE_VERSION}"
ENV DISTRO="${NODE_DISTRO}"
ENV PATH=/usr/local/lib/nodejs/node-$VERSION-$DISTRO/bin:$PATH

#RUN <<EOF
# -----------------------
# Install Extism JS PDK
# -----------------------
#export TAG="v0.5.0"
#export ARCH="aarch64"
#export ARCH="x86_64"
#export  OS="linux"
#curl -L -O "https://github.com/extism/js-pdk/releases/download/$TAG/extism-js-$ARCH-$OS-$TAG.gz"
#gunzip extism-js*.gz
#chmod +x extism-js-*
#mv extism-js-* /usr/local/bin/extism-js
#EOF

# ------------------------------------
# Install Spin
# ------------------------------------
RUN <<EOF
mkdir spin-framework
cd spin-framework
curl -fsSL https://developer.fermyon.com/downloads/install.sh | bash
mv spin /usr/local/bin/
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
mv wws /usr/bin
rm install.sh
EOF

# Create new regular user `${USER_NAME}` and disable password and gecos for later
# --gecos explained well here: https://askubuntu.com/a/1195288/635348
RUN adduser --disabled-password --gecos '' ${USER_NAME}

#  Add new user `${USER_NAME}` to sudo group
RUN adduser ${USER_NAME} sudo

# Ensure sudo group users are not asked for a password when using 
# sudo command by ammending sudoers file
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Set the working directory
WORKDIR /home/${USER_NAME}

# Set the user as the owner of the working directory
RUN chown -R ${USER_NAME}:${USER_NAME} /home/${USER_NAME}

# Switch to the regular user
USER ${USER_NAME}

# ------------------------------------
# Install Rust + Wasm Toolchain
# ------------------------------------
RUN <<EOF
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
export RUSTUP_HOME=~/.rustup
export CARGO_HOME=~/.cargo
export PATH=$PATH:$CARGO_HOME/bin
curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh 
rustup target add wasm32-wasi
EOF


# ------------------------------------
# Install GoReleaser
# ------------------------------------
RUN <<EOF
go install github.com/goreleaser/goreleaser@latest
EOF

# ------------------------------------
# Install OhMyBash
# ------------------------------------
RUN <<EOF
bash -c "$(curl -fsSL https://raw.githubusercontent.com/ohmybash/oh-my-bash/master/tools/install.sh)"
EOF

# ------------------------------------
# Install K8S tools
# ------------------------------------
RUN <<EOF
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

chmod +x kubectl
mkdir -p ~/.local/bin
mv ./kubectl ~/.local/bin/kubectl

curl -sS https://webi.sh/k9s | sh
EOF

# Command to run when starting the container
CMD ["/bin/bash"]