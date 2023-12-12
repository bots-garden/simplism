Only with flock mode and spawn mode ?
and if discovery mode is activated

- [ ] Yes
- [ ] No

perhaps use the notifiy system (the same for the killing notification)


in `server.go`
- add a `/function` handler

```go
	if wasmArgs.ServiceDiscovery == true {
		fmt.Println("🤖 this service is a service discovery")
		http.HandleFunc("/discovery", discoveryHandler(wasmArgs))
	}
```