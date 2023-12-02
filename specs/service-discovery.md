
- one simplism process store the data about the other simplism processes
- the other simplism processes will send data to the discovery process

- use flags
- don't forget to use the flags with the flock mode and config mode


```bash
# discovery manager
simplism listen simple.wasm handle --http-port 8080 --discovery true
# if --discovery is true -> add an handler

# discoverable services
simplism listen hello.wasm handle --http-port 8081 --discovery-endpoint http://localhost:8080/discovery
simplism listen hey.wasm handle --http-port 8082 --discovery-endpoint http://localhost:8080/discovery
# if --discovery-endpoint not empty -> go routine to ping the discovery manager
```

## Remark

Theorically, the discovery manager could be write with another program like Nodejs

<!--TODO: check if a wasm plug-in can call the endpoints of simplism -->


<!--TODO: remove the discovery mode of the flock mode -->