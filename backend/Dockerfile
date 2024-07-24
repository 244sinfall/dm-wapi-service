FROM golang:latest as build

ENV CGO_ENABLED=0

WORKDIR /build
COPY . .
RUN go build  -o app .
 
FROM alpine:latest as run
WORKDIR /app
COPY --from=build /build/app .
COPY . .
EXPOSE 80

ENTRYPOINT ["/app/app"]