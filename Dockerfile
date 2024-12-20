# OS
FROM rockylinux:8

WORKDIR /tmp

# Install Dependencies
RUN dnf update -y;
RUN dnf install wget make gcc rpm-build rpmdevtools redhat-rpm-config perl perl-IPC-Cmd perl-Test-Simple git python3 -y;

# Install Google Cloud CLI
RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-437.0.0-linux-x86_64.tar.gz
RUN tar -xf google-cloud-cli-437.0.0-linux-x86_64.tar.gz
RUN CLOUDSDK_CORE_DISABLE_PROMPTS=1 ./google-cloud-sdk/install.sh
ENV PATH=$PATH:/tmp/google-cloud-sdk/bin

# Install Golang
RUN wget https://dl.google.com/go/go1.21.6.linux-amd64.tar.gz; \
    tar -xvf go1.21.6.linux-amd64.tar.gz; \
    mv go /usr/local

# Set golang env variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN go install golang.org/x/tools/gopls@v0.16.0
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

# Setup RPM build environment
RUN rpmdev-setuptree

# Default powerline10k theme, no plugins installed
RUN dnf install -y zsh

# Change workspace
WORKDIR /ticrypt-file-copy
