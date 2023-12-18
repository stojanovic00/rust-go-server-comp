FROM rust:latest as builder

WORKDIR /usr/src/app

COPY . .

RUN cargo build --release

FROM debian:stable-slim

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/target/release/rust_server .

EXPOSE 8001

#For graceful shutdown
STOPSIGNAL SIGINT

CMD ["./rust_server"]
