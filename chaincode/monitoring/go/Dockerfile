# This image is a microservice in golang for the Degree chaincode
FROM golang:1.18.3-alpine AS build

COPY ./ /go/src/github.com/basic
WORKDIR /go/src/github.com/basic

# Build application
RUN go build -o chaincode -v .

# Production ready image
# Pass the binary to the prod image
FROM alpine:3.16 as prod

COPY --from=build /go/src/github.com/basic/chaincode /app/chaincode

USER 1000

WORKDIR /app
CMD ./chaincode