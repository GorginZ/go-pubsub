FROM python:3.9.0-alpine3.12 as ci
#gcloud needs python 
RUN wget https://releases.hashicorp.com/terraform/1.3.0/terraform_1.3.0_linux_amd64.zip
RUN unzip terraform_1.3.0_linux_amd64.zip && rm terraform_1.3.0_linux_amd64.zip
RUN mv terraform /usr/bin/terraform

RUN wget https://github.com/mikefarah/yq/releases/download/v4.28.1/yq_linux_amd64.tar.gz -O - |\
  tar xz && mv yq_linux_amd64 /usr/bin/yq

RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-404.0.0-linux-x86_64.tar.gz

RUN mkdir -p /usr/local/gcloud \
    && tar -C /usr/local/gcloud -xvf /google-cloud-cli-404.0.0-linux-x86_64.tar.gz \ 
    && /usr/local/gcloud/google-cloud-sdk/install.sh -q

ENV PATH $PATH:/usr/local/gcloud/google-cloud-sdk/bin
