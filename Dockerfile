FROM golang:latest as builder

WORKDIR /yprog
ENV CGO_ENABLED=0
COPY . .
RUN go build -o yprog .

FROM  alpine as runner

WORKDIR /yprog
COPY --from=builder yprog/ ./

CMD ["/yprog/yprog"]
EXPOSE 8085
