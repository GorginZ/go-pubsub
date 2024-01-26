# Order app

Run locally

>note this example uses service account key credentials, this is really bad, delete your key after you play.  Long lived credentials are both a security concern and a operational hassle and are only used for convenience in this exercise. 

```
export PROJECT_ID=<prj-id>
export TOPIC_ID=<topic_id>
export AUTH_JSON=<path to service account key>
```

in /pubsub/order

```
go run main.go
```


Pretend you're some other service sending an order through

```
➜  curl -X POST -H "Content-Type: application/json" -d '{"Email": "georgialeng@.com", "Product": "car", "Amount": 99}'  http://localhost:8080/order;
{"message":"order created"}%
```