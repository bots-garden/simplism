FROM k33g/simplism:0.1.1
COPY hello.wasm .
EXPOSE 8080
CMD ["/simplism", "listen", "hello.wasm", "handle", "--http-port", "8080", "--log-level", "info"]
