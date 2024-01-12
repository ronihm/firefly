FROM golang:1.21
WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download
COPY . ./
RUN go build -o /builtapp
CMD [ "/builtapp" ]