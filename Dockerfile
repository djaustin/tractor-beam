FROM golang:alpine as builder
WORKDIR /build 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o tractor-beam .
FROM scratch
COPY --from=builder /build/tractor-beam /app/
WORKDIR /app
ENV PATH=/app/:$PATH
CMD ["tractor-beam"]