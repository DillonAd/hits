FROM golang:1.15-alpine AS build

RUN apk update
RUN apk add git

RUN mkdir /app
COPY ./src /app
WORKDIR /app

RUN go env -w GO111MODULE=on
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/hits .

FROM scratch AS deploy

COPY --from=build /app/bin/hits ./hits
EXPOSE 8080

CMD ["./hits"]