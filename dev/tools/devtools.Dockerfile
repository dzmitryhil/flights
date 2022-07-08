FROM golang:1.17-alpine3.16

RUN apk add make gcc libc-dev git bash

RUN GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
RUN GO111MODULE=on go get github.com/vasi-stripe/gogroup/cmd/gogroup@v0.0.0-20200806161525-b5d7f67a97b5
RUN GO111MODULE=on go get mvdan.cc/gofumpt@v0.0.0-20200927160801-5bfeb2e70dd6

RUN dir=$(mktemp -d) && \
    git clone --depth 1 -b v0.27.0 https://github.com/go-swagger/go-swagger "$dir" && \
    cd "$dir" && \
    go install -ldflags "-X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=v0.27.0" ./cmd/swagger && \
    rm -rf "$dir"
RUN rm -rf /root/.cache/go-build/ /go/pkg/*

COPY entrypoint.sh /entrypoint.sh
RUN mkdir -p /root/.cache/ && \
    ln -s /cache/golangci-lint/ /root/.cache/golangci-lint && \
    ln -s /cache/go-build/ /root/.cache/go-build

WORKDIR /go/src/github.com/dzmitryhil/flights

ENTRYPOINT ["bash", "/entrypoint.sh"]
