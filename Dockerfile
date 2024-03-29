FROM python:3.13.0a2-alpine3.19 as ci-infra
#gcloud needs python 
RUN wget https://releases.hashicorp.com/terraform/1.6.6/terraform_1.6.6_linux_amd64.zip
RUN unzip terraform_1.6.6_linux_amd64.zip && rm terraform_1.6.6_linux_amd64.zip
RUN mv terraform /usr/bin/terraform

RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-458.0.1-linux-x86_64.tar.gz

RUN mkdir -p /usr/local/gcloud \
    && tar -C /usr/local/gcloud -xvf /google-cloud-cli-458.0.1-linux-x86_64.tar.gz \ 
    && /usr/local/gcloud/google-cloud-sdk/install.sh -q

ENV PATH $PATH:/usr/local/gcloud/google-cloud-sdk/bin

FROM golang:1.21.6-alpine3.18 as go-ci
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.49.0
