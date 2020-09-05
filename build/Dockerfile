FROM golang:1.14
WORKDIR /app
COPY . /app/
RUN go mod download && go build ./cmd/service
RUN ls -al /app
CMD ["/app/service"]