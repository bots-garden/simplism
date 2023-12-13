FROM k33g/simplism:0.0.7
COPY hello.wasm .
EXPOSE 8000
CMD ["/simplism", "listen", "hello.wasm", "handle", "--http-port", "8000", "--log-level", "info"]
