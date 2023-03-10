FROM golang:1.20.2 AS builder

COPY . build/

WORKDIR build/

RUN CGO_ENABLED=0 go build -o app . && ls -lah app

FROM scratch 

COPY --from=builder --chown=1000:1000 /go/build/app /app 

USER 1000:1000

EXPOSE 8080
EXPOSE 6969

CMD [ "/app" ]
