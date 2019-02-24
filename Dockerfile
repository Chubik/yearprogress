FROM golang:alpine as builder
# RUN apk update && apk add git && go get gopkg.in/natefinch/lumberjack.v2

WORKDIR /yprog
ENV CGO_ENABLED=0
COPY . .
RUN go build -o yprog .

FROM  alpine as runner

WORKDIR /yprog
COPY --from=builder yprog/ ./

CMD ["/yprog/yprog"]
EXPOSE 8085
