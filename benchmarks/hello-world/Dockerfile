FROM k33g/simplism:0.0.8
COPY hello-world.wasm .
EXPOSE 8080
CMD ["/simplism", "listen", "hello-world.wasm", "handle", "--http-port", "8080", "--log-level", "info"]
