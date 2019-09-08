FROM node:10.15

RUN apt-get update && \
    apt-get install python-dev python-pip python-setuptools -y && \
    easy_install --upgrade six && \
    pip install awscli && \
    npm i -g yarn@v1.3.2

ENTRYPOINT ["/bin/sh", "-c"]
