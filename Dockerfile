FROM golang:1.12

WORKDIR /tests-parser-dir

COPY . .

RUN go build

ENTRYPOINT ["/tests-parser-dir/testsParser"]
