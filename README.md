# grpc-sample

### High performance distributed architecture to ingest massive amounts of data using GRPC, Nats and Mongodb

![image](doc/arch.png)

## J4F

The J4F (Just for fun) projects are elaborate to demonstrate concepts, show benchmarks, and provoke doubts about well-established certainties

#

## Run Server

```sh
docker-compose up -d
```

![image](doc/ps.png)


## Run client

```sh
make client
```

## Samples
```sh
grpc-client --messages 100000 --concurrency 50 --servers 10
```

Sending 100k messages under 1 second using 50 routines

![image](doc/100000.png)

```sh
grpc-client --messages 1000000 --concurrency 50 --servers 10
```
Sending 1M messages under 8 seconds using 50 routines

![image](doc/1000000.png)

```sh
grpc-client --messages 10000000 --concurrency 50 --servers 10
```
Sending 10M messages under 75 seconds using 100 routines

![image](doc/10000000.png)