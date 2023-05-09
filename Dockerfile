FROM alpine:3.6

RUN echo -e "https://mirrors.aliyun.com/alpine/v3.8/community/\nhttps://mirrors.aliyun.com/alpine/v3.8/main/" > /etc/apk/repositories

ENV TZ Asia/Shanghai

RUN apk add --update --no-cache \
        tzdata \
        && rm -f /etc/localtime \
        && ln -s /usr/share/zoneinfo/$TZ /etc/localtime \
        && echo $TZ > /etc/timezone

RUN apk add --update --no-cache \
    ca-certificates \
    && rm -rf /var/cache/apk/*

COPY shorturl /bin/shorturl
COPY docker-entrypoint /usr/local/bin/docker-entrypoint

WORKDIR /usr/local/var/shorturl

EXPOSE 11050 11051
VOLUME ["/opt", "/usr/local/etc/shorturl"]

ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
