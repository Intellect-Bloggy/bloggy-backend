# Build
FROM golang:alpine AS build

WORKDIR /go/src/bloggy-backend/

COPY go.mod .
COPY go.sum .

RUN apk add make
COPY Makefile .

RUN make deps
COPY . .
RUN make docs
RUN go build -o bloggy_backend cmd/app/main.go

# Deploy
FROM alpine:latest
COPY --from=build /go/src/bloggy-backend/bloggy_backend .

ENTRYPOINT [ "./bloggy_backend" ]
