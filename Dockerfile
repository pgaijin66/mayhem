FROM golang:1.24-alpine AS build_base
ARG LDFLAGS
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "${LDFLAGS}" \
    -a -installsuffix cgo \
    -o ./bin/phailure \
    main.go

FROM alpine:3.15

RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

RUN addgroup -g 1000 phailure && \
    adduser -D -s /bin/sh -u 1000 -G phailure phailure

WORKDIR /app

COPY --from=build_base /src/bin/phailure /app/phailure

RUN chown -R phailure:phailure /app

USER phailure

RUN chmod +x phailure

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/_chaos/health || exit 1

ENV phailure_PORT=8080
ENV phailure_LOG_LEVEL=info

ENTRYPOINT ["./phailure"]

CMD ["--port", "8080"]