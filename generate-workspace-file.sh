#!/bin/bash

cat > ./simplism-project.code-workspace <<- EOM
{
  "folders": [
    {
      "name": "benchmarks",
      "path": "$(pwd)/benchmarks"
    },
    {
      "name": "rust-examples",
      "path": "$(pwd)/samples/rustlang"
    },
    {
      "name": "go-examples",
      "path": "$(pwd)/samples/golang"
    },
    {
      "name": "simplism",
      "path": "$(pwd)"
    },

  ],
  "settings": {
    "workbench.iconTheme": "material-icon-theme",
    "workbench.colorTheme": "Cobalt2",
    "terminal.integrated.fontSize": 15,
    "files.autoSave": "afterDelay",
    "files.autoSaveDelay": 1000
  }
}
EOM
