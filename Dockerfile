# docker build --tag=kubernetes-example-gateway .
# docker run -it -p 8080:5000 kubernetes-example-gateway

FROM golang:1.10.3 as builder
WORKDIR  /go/src/bitbucket.org/revenuemonster
ADD . /go/src/bitbucket.org/revenuemonster
# Install Dependencies
RUN go get github.com/labstack/echo
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/go-redis/redis
RUN make

FROM alpine
WORKDIR  /app
# install ca cert if you want expose the app directly using load balancer
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/bitbucket.org/revenuemonster/kubernetes-example-gateway /app
# Container Environment
# It will be overwrite by deployment env if same key exist
ENV SYSTEM_NAME 'Revenue Monster Kubernetes Session'
ENTRYPOINT ./kubernetes-example-gateway