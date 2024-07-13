FROM golang:1.22.5-alpine3.19 AS builder

RUN mkdir /src /src/build
WORKDIR /src

RUN apk add make
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN make build_api
RUN make build_db_cmd

FROM alpine:3.16 AS runner

WORKDIR /src

COPY --from=builder /src/build/api ./
COPY --from=builder /src/build/db ./

COPY .env ./
