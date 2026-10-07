# STEP2: BUILD BACKEND
FROM golang:1.26-trixie AS builder_backend
WORKDIR /home/filestash/
COPY . .
RUN apt-get update > /dev/null && \
    apt-get install -y curl make > /dev/null 2>&1 && \
    apt-get install -y libjpeg-dev libtiff-dev libpng-dev libwebp-dev libraw-dev libheif-dev libgif-dev libvips-dev > /dev/null 2>&1 && \
    apt-get install -y libavcodec-dev libavdevice-dev libavfilter-dev libavformat-dev libswresample-dev libswscale-dev libavutil-dev > /dev/null 2>&1 && \
    make init && \
    make build

# STEP3: BUILD PROD IMAGE
FROM debian:stable-slim
MAINTAINER mickael@kerjean.me
WORKDIR /app/
COPY --from=builder_backend /home/filestash/dist/ .
RUN apt-get update > /dev/null && \
    apt-get -y upgrade > /dev/null && \
    apt-get install -y --no-install-recommends ca-certificates curl ffmpeg libbrotli1 poppler-utils && \
    useradd filestash && \
    mkdir -p /app/data/state/ && \
    chown -R filestash:filestash /app/ && \
    chmod 730 /app/filestash && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/*

USER filestash
CMD ["/app/filestash"]
EXPOSE 8334
