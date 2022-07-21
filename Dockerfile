FROM golang:1.18-bullseye AS build

RUN useradd -u 11801 -m radio

WORKDIR /app

COPY . ./
RUN go mod download


RUN go build -ldflags="-s -w" -o /radiohelsinkitospotify


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /radiohelsinkitospotify /radiohelsinkitospotify
COPY /static/index.html /static/index.html

USER 11801

EXPOSE 8080

ENTRYPOINT ["/radiohelsinkitospotify"]
