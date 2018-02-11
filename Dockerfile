############################################
# version : godcong/oauth2server:latest
# desc :
############################################
FROM golang:latest

MAINTAINER godcong (jumbycc@163.com)

#RUN mkdir -p /home/oauth2server
RUN mkdir -p /home/config

WORKDIR /go/src/github.com/godcong/oauth2server


#ADD static /home/oauth2server/static
#ADD templates /home/oauth2server/templates

#ADD docker/server /home/oauth2server/server
#ADD docker/cmd /home/oauth2server/cmd
#ADD docker/config.toml /home/oauth2server/config.toml
ADD . /go/src/github.com/godcong/oauth2server

RUN go build ./main/server
RUN go build ./main/cmd

EXPOSE 8080

# CMD ["./cmd","/home/config.toml"]
CMD ["./server","-c","/home/config/config.toml"]