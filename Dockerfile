FROM golang:1.19-alpine AS builder

WORKDIR /build

ADD go.mod .

ENV GO111MODULE=on

COPY . .

# download Go modules and dependencies
RUN go mod download

RUN go build -o crm-service-go .

FROM alpine

WORKDIR /build

COPY --from=builder /build/crm-service-go /build/crm-service-go
COPY .env /build
COPY gcredentials.json /build

EXPOSE 3000

CMD ["./crm-service-go"]