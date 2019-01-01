FROM golang:1.11

RUN apt-get update -y && \
    apt-get install python-dev python-pip -y && \
    easy_install --upgrade six && \
    pip install awscli

ENTRYPOINT ["/bin/sh", "-l", "-c"]

