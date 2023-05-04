FROM golang:latest AS builder
WORKDIR /build
COPY . .
RUN make

FROM scratch
COPY --from=builder /build/haniho /app
EXPOSE 8000
CMD ["/app", "--server"]
