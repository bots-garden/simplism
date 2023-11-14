
## string to json

    "github.com/tidwall/gjson"
## create json string

    "github.com/tidwall/sjson"




```bash
Usage:
  extism call [flags]

Flags:
      --allow-host stringArray   Allow access to an HTTP host, if no hosts are listed then all requests will fail. Globs may be used for wildcards
      --allow-path stringArray   Allow a path to be accessed from inside the Wasm sandbox, a path can be either a plain path or a map from HOST_PATH:GUEST_PATH
      --config stringArray       Set config values, should be in KEY=VALUE format
  -h, --help                     help for call
  -i, --input string             Input data
      --log-level string         Set log level: trace, debug, warn, info, error
      --loop int                 Number of times to call the function (default 1)
  -m, --manifest                 When set the input file will be parsed as a JSON encoded Extism manifest instead of a WASM file
      --memory-max int           Maximum number of pages to allocate
      --set-config config        Create config object using JSON, this will be merged with any config arguments
      --stdin                    Read input from stdin
      --timeout uint             Timeout in milliseconds
      --wasi                     Enable WASI

Global Flags:
      --github-token string   Github access token, can also be set using the $GITHUB_TOKEN env variable
  -q, --quiet                 Enable additional logging
  -v, --verbose               Enable additional logging
```