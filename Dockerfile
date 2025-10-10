FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod .
COPY . .
RUN go build -o /bin/dman ./cmd/server


FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /bin/dman /bin/dman
EXPOSE 8080
ENTRYPOINT ["/bin/dman"]
