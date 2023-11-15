# Simplism: a tiny HTTP server for Extism

![image](imgs/simplism-small-logo.jpeg)

## How is Simplism developed?

Simplism is developed in Go with **[Wazero](https://wazero.io/)**[^1] as the Wasm runtime and **[Extism](https://extism.org/)**[^2], which offers a Wazero-based Go SDK and a Wasm plugin system.

### Ready to use environments

- [🍊 Open it with Gitpod](https://gitpod.io/#https://github.com/bots-garden/simplism)
- [🐳 Open it with Docker Dev Environment (**✋ arm version**)](https://open.docker.com/dashboard/dev-envs?url=https://github.com/bots-garden/simplism/tree/main)
  - Prerequisites:
    - https://docs.docker.com/desktop/dev-environments/create-dev-env/#prerequisites
    - [Visual Studio Code Remote Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

[^1]: Wazero is a project from **[Tetrate](https://tetrate.io/)**
[^2]: Extism is a project from **[Dylibso](https://dylibso.com/)**