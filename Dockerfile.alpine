FROM golang:1.11.2-stretch AS build-env
ADD . /src
RUN cd /src && go build

FROM alpine:3.7

RUN apk add --no-cache bind-tools python3 conntrack-tools tcpdump iproute2
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY entrypoint.sh /
COPY entrypoint.py /

COPY --from=build-env /src/dns-utils /

ENTRYPOINT ["/entrypoint.sh"]
CMD ["google.com"]
