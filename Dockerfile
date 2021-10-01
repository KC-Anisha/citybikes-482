FROM golang:1.17-alpine

WORKDIR $GOPATH/src/github.com/KC-Anisha/citybikes-482

# RUN apk update && apk add git

# RUN git config --global url."https://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:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# RUN git clone https://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:x-oauth-basic@github.com/KC-Anisha/citybikes-482.git

# ADD . /go/src/github.com/KC-Anisha/citybikes-482

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# RUN go mod tidy

COPY *.go ./

ENV LOGGLY_TOKEN=3978ab6c-18d0-4709-8d65-38a8b73f88a3

RUN go build -o /482

# RUN go install github.com/KC-Anisha/citybikes-482

EXPOSE 8080

CMD [ "/482" ]