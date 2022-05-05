FROM tinygo/tinygo:0.23.0 AS builder

RUN apt-get install make -y

COPY . .

RUN make build-all

FROM alpine

COPY --from=builder ./receiver.wasm ./sender.wasm ./

