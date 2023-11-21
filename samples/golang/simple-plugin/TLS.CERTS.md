# Use of TLS Certificates

```bash
sudo chmod 777 simplism.bots.garden.*
```

To test it locally, add `0.0.0.0 simplism.bots.garden` to the `hosts` file of the client machine`. Then serve the wasm service like this:

```bash
go run ../../../main.go listen \
simple.wasm say_hello --http-port 443 \
  --cert-file simplism.bots.garden.crt \
  --key-file simplism.bots.garden.key
```

And, finally, query the service like this:

```bash
curl https://simplism.bots.garden
```














