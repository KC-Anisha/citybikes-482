# # FROM golang:onbuild
# # EXPOSE 8080

# FROM golang:latest

# # expose default port
# EXPOSE 8000

# # set environment path
# ENV PATH /go/bin:$PATH

# # cd into the api code directory
# WORKDIR /go/src/github.com/KC-Anisha/citybikes-482

# # create ssh directory
# RUN mkdir ~/.ssh
# RUN touch ~/.ssh/known_hosts
# RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

# # allow private repo pull
# RUN git config --global url."https://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:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# # copy the local package files to the container's workspace
# ADD . /go/src/github.com/KC-Anisha/citybikes-482

# # install the program
# RUN go install github.com/KC-Anisha/citybikes-482

# # start application
# CMD ["citybikes-482"] 

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# RUN go mod tidy

COPY *.go ./

RUN go build -o /482

EXPOSE 8080

CMD [ "/482" ]