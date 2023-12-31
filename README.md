# GO-PUBSUB

WIP

Inspired by this but want to use terraform and go done while preparing for gcp architect exam.

https://github.com/googleapis/python-pubsub


todo:

docker tf version


#  Usage

## auth:

Auth with gcloud or use docker gcloud like so:

```docker compose run gcloud auth application-default login```

Follow prompts.

## initialise the terraform workspace

Locally:

```docker compose run terraform -chdir=terraform init```

Then tf plan:

```docker compose run terraform -chdir=terraform plan```

Inspect the resources in the plan output.

Initially we will apply using a local backend. Optionally we can migrate this state to our bucket created in the main.tf.

---

### optional

### To use a remote backend


Create a file terraform/backend.tf

Put the value of the bucket, this will hold our tf state and we will migrate the state by re initialising the tf workspace

```
 terraform {
   backend "gcs" {
     # bucket = "replaceme"
     prefix = "terraform/state"
   }
 }
```
