FROM golang:latest AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM scratch
COPY --from=builder /build/app /app
EXPOSE 8000
CMD ["/app", "--server"]
