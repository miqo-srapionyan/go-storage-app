# go-storage-app
## General info
HI there 👋, this is storage app written in Go, that importing big CSV file into MySQL database, 
and exposing and API endpoint to request data.

## Technologies
Project is created with:
* Go
* MySQL
* Redis
* Docker

## Setup
To run this project, make sure you have installed Docker, and these ports are not busy
* 3306 - for MySQL
* 6379 - for Redis
* 1321 - for localhost:1321/promotions/:id

install it locally using:

```
$ clone repo
$ cd go-storage-app.git
$ docker-compose up -d --build
```

This will setup mysql database, and import default csv file(promotions.csv) into it.

## API
Simply request localhost:1321/promotions/:id

## Additional info, answer of questions
As application receiving CSV file per 30 min, and clients can read data frequently, i assume that application is read heavy.

### The steps.

#### Setup database:
First of all, to insert CSV we must choose database, it can be Apache Cassandra, Apache HBase (they are designed to handle massive amounts of data and provide high scalability) but as i do not have experience with them i would choose MySQL. We must structure DB, create table add necessary columns.

#### Read and parse CSV:
We can use a streaming CSV parser: Instead of loading the entire CSV file into memory, we should process the file line by line. This way, we can avoid loading the entire file into memory.
we can Parallelize the Insertion: parallelizing the insertion process by dividing the data into multiple batches and inserting them concurrently using multiple threads or processes.

#### Endpoint:
We can use Go web framework to expose HTTP endpoint. We receive ID and we can query MySQL database with that id. It will take O(logn) (big O notation) time to receive an item as MySQL by default uses B-tree for indexes.

#### Read heavy - performance:
We should use appropriate indexing on columns that are frequently queried (in this case primary key ID is enough), also we should consider partitioning. At peak periods we should scale horizontally, by having replicas. Also we can use table sharding, to have smaller tables. Here I used LRU(Least Recently Used) technique, which is caching the views in redis, then if storage is full, it is starting removing least viewed items. 

#### Handling Immutable Files:
For each new file, we can create a new table, and import the data into the corresponding table. We can clean up or archive the previous tables as needed.

#### Deployment, Scaling, and Monitoring:
We can deploy our application to cloud servers, we can use Docker for easier deployment and scaling. Set up auto-scaling mechanisms based on the load using orchestration tools like Kubernetes. Monitor the application's performance, resource utilization, and MySQL metrics using monitoring tools like Prometheus and Grafana.

I would add a unit tests with database and redis mocking for this, unfortunately i'm running out of time.

##### Thank You
