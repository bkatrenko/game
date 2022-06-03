# run build
FROM golang:1.18-alpine as build-env
ENV APP_NAME hockey-server

COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME

RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

# run stage
FROM alpine:3.14
ENV APP_NAME hockey-server

COPY --from=build-env /$APP_NAME .
CMD ./$APP_NAME
