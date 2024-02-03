#!/bin/bash

# ------------------------------------
# Install Code Server Extensions
# ------------------------------------
code-server --install-extension wesbos.theme-cobalt2
code-server --install-extension PKief.material-icon-theme
code-server --install-extension PKief.material-product-icons
code-server --install-extension golang.go
code-server --install-extension rust-lang.rust-analyzer
code-server --install-extension aaron-bond.better-comments
code-server --install-extension GitHub.github-vscode-theme
code-server --install-extension huytd.github-light-monochrome
code-server --install-extension pomdtr.excalidraw-editor

echo "🌍 open: http://0.0.0.0:${CODER_HTTP_PORT}"
echo "🔐 if you activated https mode:" 
echo "🌍  open: https://${LOCAL_DOMAIN}:${CODER_HTTP_PORT}"

