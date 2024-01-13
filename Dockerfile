FROM golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ffapp .

CMD [ "/usr/src/app/ffapp" ]