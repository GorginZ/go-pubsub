# Order service

Run locally

>note this example uses service account key credentials, this is really bad, delete your key after you play.  Long lived credentials are both a security concern and a operational hassle and are only used for convenience in this exercise. 


manual step: create a service account key for the service account created by the terraform for this project and download it

```
export PROJECT_ID=<prj-id>
export TOPIC_ID=<topic_id>
export AUTH_JSON=<path to service account key>
```

in /pubsub/order

```
go get ./...
```

start order service
```
go run main.go
```


Pretend you're some other service sending an order through

```
➜  curl -X POST -H "Content-Type: application/json" -d '{"Email": "georgialeng@.com", "Product": "car", "Amount": 99}'  http://localhost:8080/order;
{"message":"order created"}%
```

will be able to see this traffic hitting the topic in console.