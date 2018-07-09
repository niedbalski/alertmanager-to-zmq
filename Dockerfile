FROM golang:stretch as build
MAINTAINER  Jorge Niedbalski <jnr@metaklass.org>

WORKDIR /go/src/github.com/niedbalski/alertmanager-to-zmq
COPY . .

#RUN apt-get -yyq update && apt-get -yyq install libsodium18 libsodium-dev lzma lzma-dev libunwind8 libunwind8-dev
#RUN wget https://github.com/zeromq/libzmq/releases/download/v4.2.2/zeromq-4.2.2.tar.gz -O zmq.tar.gz && \
#   tar -xvzf zmq.tar.gz && cd zeromq* && ./configure && make install && ldconfig
##RUN  go get github.com/tools/godep && godep restore
#RUN CGO_LDFLAGS="$CGO_LDFLAGS -lstdc++ -lm -lsodium -lunwind" \
#  CGO_ENABLED=1 \
#  GOOS=linux \
#  go build -v -a --ldflags '-extldflags "-static" -v' -o alertmanager-to-zmq
RUN apt-get -yyq update && apt-get -yyq install libzmq3-dev
RUN go build -v -o alertmanager-to-zmq

FROM debian:stretch
WORKDIR /app
RUN apt-get -yyq update && apt-get -yyq install libzmq3-dev
COPY --from=build /go/src/github.com/niedbalski/alertmanager-to-zmq/alertmanager-to-zmq .

ENTRYPOINT ["./alertmanager-to-zmq"]
