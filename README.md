# GO-PUBSUB

![main](https://github.com/GorginZ/go-pubsub/actions/workflows/lint-and-test.yaml/badge.svg?)


Inspired by this google learning exercise [video](https://www.youtube.com/watch?v=f5DOsB7Nlw0&list=PLIivdWyY5sqKwVLe4BLJ-vlh9r9zCdOse&index=3&t=317s) presented by GCP's wonderful Priyanka Vergadia.
But instead:
- use terraform for gcp resources
- go for applications
- testing
- tooling for local dev and testing

This project is simply a fun learning exercise for me while preparing for gcp architect exam.

Required:
- [docker](https://www.docker.com/)
- Access to a [GCP](https://console.cloud.google.com) account with billing enabled


#  Usage
```/Terraform``` has resources to create pubsub topics and subscriptions 

Everything is in docker to minimise requirements and to allow to focus on the exercise, which is simply to play with GCP pub-sub


## Set up

create a ```terraform/localvars.tfvars``` file.

```
billing_account = "<YOUR_BILLING_ACC_ID>"
```

This is required to register the project with a billing account.


### Authenticate with application credentials:

This authorises the gcp sdk, you will need at least the following permissions:

- project creator 

If you are working in your personal GCP account or org and working as your owner principal you should not lack any permissions. 

Auth with gcloud or use docker gcloud like so:

```docker compose run gcloud auth application-default login```
>use sudo if needed on linux

Follow prompts.

## initialise the terraform workspace

Locally:

```docker compose run terraform -chdir=terraform init```

Then tf plan:

```docker compose run terraform -chdir=terraform plan -var-file=localvars.tfvars```

Inspect the resources in the plan output.

Proceed to apply:

```docker compose run terraform -chdir=terraform apply -var-file=localvars.tfvars```

enter "yes"

Initially we will apply using a local backend. Optionally we can migrate this state to our bucket created in the main.tf.

We now have created:
- our pubsub project where we'll deploy the resources for this exercise
- pubsub topic
- pubsub subscription x2
- pubsub service account

---

### optional

### To use a remote backend


Create a file ```terraform/backend.tf```

Put the value of the bucket, this will hold our tf state and we will migrate the state by re initialising the tf workspace


```
 terraform {
   backend "gcs" {
     # bucket = "replaceme"
     prefix = "terraform/state"
   }
 }
```

then re initialise and migrate the state, this is simple, read about how to do this for a [GCS backend here](https://cloud.google.com/docs/terraform/resource-management/store-state) prompts .

### Cleanup GCP resources

```
docker compose run terraform -chdir=terraform destroy -var-file=localvars.tfvars
```


### Run the services

[Order Service](/go-pubsub/order/README.md)