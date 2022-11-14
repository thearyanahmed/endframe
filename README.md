## Running the project
Run the following commands to run the project.

- `git clone git@github.com:thearyanahmed/nordsec.git`
- `cd nordsec`
- `cp .env.example .env`
- Update proper values with `.env`. For this demo `.env` values have be filled in `.env.example`
- `make start` will start the necessary containers. 

## Running tests
While running the container, run `make test` to run the tests. Or to run outside of container run `go test -v ./...` 

## Run simulation
Running `make simulate` would start a simulation. The simulation steps are as follows

1. It would spawn some riders with in the latitude bound left bottom 52.50000 , 13.46000 and top right 52.70000 ,13.54000 range. 
2. After that, it'll replicate some (5) clients, select a random ride and sets off to a different destination.
3. The start ride endpoint will return a `route` property, which is an array of `points (lat, lon)`. 
4. On the client side, for every route point it will ping it's current location to `/api/v1/trip/notify/location` and sleep for 3 seconds. Replicating traveling time.
5. After it reaches the last destination, it will finally end it's trip by sending a request to `/api/v1/trip/end`. On the server side, the ride will enter a cooldown mode.

## Project structure
The codebase lives in `/pkg` directory. The entry point to the project is as of now are REST http endpoints. It will take the http request and send it to downstream application service.
The router will call the handlers, which handles the business logic using `ride service`. The `ride service` is where our main business logic lives. The ride service calls
`location service` as it requires. `location service` a datastore (`redis` in this scenario) via a repository.

The reason separating our logic to ride service is because of independence of responsibilities. Our handlers just validates the request and passes down the `ride service`. 
This gives us the ability to add on more ways to port core service. If we want to use graphql, or grpc or some other method of providing our service, we just need to
add a handler/equivalent layer and not think about our business logic as it lives by itself. 

The `location service` and `ride service` are separate. Because the `location service` only keeps the rides data. `ride service` will handle location data and other data and logic. 
As for this demonstration, we are currently using redis, to keep things simple. Thus `location service` has the only repository.

`Logging` has been left out mostly. We would also `trace (tracing and spanning`) the activities, successful scenarios and unsuccessful scenarios as well.

```
┌──────────┐    ┌──────────┐    ┌───────────┐    ┌────────────────────────────────────────┐
│          │    │          │    │           │    │           Request Serializer           │
│  Server  ├────►  Router  ├────►  Handler  ├────►────────────────────┬───────────────────┼───┐
│          │    │          │    │           │    │       Serializer   │     Validator     │   │
└──────────┘    └──────────┘    └───────────┘    └────────────────────┴───────────────────┘   │
                                                                                              │
     ┌────────────────────────────────────────────────────────────────────────────────────────┘
     │
┌────▼────┐                                       ┌─────────────────┐
│ Ride    ├─────────────────────────────────────► │     Response    │
│ Service │                                       └─────────────────┘    
└────↑────┘ ┌──────────────────────┐                ┌─────────────────────┐             
     │      │                      │                │                     │
     └──────►  Location Service    │←────────────►  │ Repository (Redis ) │
            │                      │                │                     │
            └──────────────────────┘                └─────────────────────┘
```

## Location service

Our location service keeps track of the rides. It uses `redis`'s `ZADD`. It is a sorted set, with score values. Score values are unix timestamps of the event. The keys are 
`geohash`. We use lat,lon to calculate the geohash. The records are as follows

```
geohash-01 timestamp-01 [event0, event1, ...eventN]
geohash-01 timestamp-02 [event0, event1, ...eventN]
geohash-01 timestamp-0N [event0, event1, ...eventN]
geohash-0N timestamp-0N [event0, event1, ...eventN]
```

**Retrieving nearby rides** 

To retrieve nearby rides, we need the coordinate pair. Where we will take left bottom point and right top point. We will use `x1y1` and `x3y3`. Then, based on that, we take the 
center of the bounding box and get the nearby neighbours (area). And from the sorted set `ZADD` set values, we run a `ZRANGE` on those keys. Each geolocation can have multiple records
of the same ride. So we filter them out and sort them based on score, which are the timestamp. 

It would be something like this,
1. GET EVENTS WITHIN GEOHASH [geohash1,geohash2,...,geohashN]
2. FLAT THAT LIST RETRIEVED, TAKE ONLY UNIQUE BASED ON RIDE_UUID
3. SORT THEM BY TIMESTAMP

We will have our rides in near by location. Then we apply the cooldown filter and unique status filter to get the ride with the appropriate ride status.

**Cooldown**

One requirement was make the ride unavailable for an amount of time after one ride has finished. We are using a simple redis key value with expiration. The first thought was 
to trigger some queue for a later time. Or even a goroutine. But both are depended on application layer, in case our queue or goroutine fails to trigger/ server is unavailable, 
the ride will not come out of cooldown status. 

As we needed some automation way to update the cooldown data value, `redis`'s expiration will help us in that and also it syncs with the datastore. So we have a single source of truth.

**Area Coordinates Bounding Box**

```
┌──────────────────────────┐  
│(x4,y4)            (x3,y3)│  
│                          │
│(x1,y1)            (x2,y2)│    
└──────────────────────────┘    
```

**Geohash Tile Sizes**
```bash
Length	Tile Size
1	5,009.4km x 4,992.6km
2	1,252.3km x 624.1km
3	156.5km x 156km
4	39.1km x 19.5km
5	4.9km x 4.9km # | Our default value 
6	1.2km x 609.4m 
7	152.9m x 152.4m
8	38.2m x 19m
9	4.8m x 4.8m
10	1.2m x 59.5cm
11	14.9cm x 14.9cm
12	3.7cm x 1.9cm
```

## API Endpoints

**health check**
```bash
  curl --location --request GET 'http://localhost:8080/api/v1/health-check'
```

**Update ride's location**
This is when the scooter is not on a trip. Think when a scooter comes online. 
```bash
curl --location --request POST 'http://localhost:8080/api/v1/ride/activate' \
--header 'Authorization: rider-key' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'ride_uuid=91122e1b-a36b-4d39-8267-95fe6c5eeb0f' \
--data-urlencode 'latitude=52.323551' \
--data-urlencode 'longitude=13.47453'
```

**Get near-by rides**
```bash
curl --location --request GET 'http://localhost:8080/api/v1/rides/near-by?x1=52.3251&y1=13.453&x3=52.3361&y3=13.475' \
--header 'Authorization: client-key'
```

**Start a trip**
```bash
curl --location --request POST 'http://localhost:8080/api/v1/trip/start' \
--header 'Authorization: client-key' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'origin_latitude=52.332425' \
--data-urlencode 'origin_longitude=13.533946' \
--data-urlencode 'ride_uuid=efb2eeb1-dfbb-4692-bdcd-3148a93eddab' \
--data-urlencode 'client_uuid=91122e1b-a36b-4d39-8267-95fe6c5eeb1f' \
--data-urlencode 'destination_latitude=53.885487' \
--data-urlencode 'destination_longitude=13.380731'
```

**Notify location while on trip**
```bash
curl --location --request POST 'http://localhost:8080/api/v1/trip/notify/location' \
--header 'Authorization: client-key' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'latitude=52.00001' \
--data-urlencode 'longitude=13.00001' \
--data-urlencode 'ride_uuid=1cc846eb-9f14-4c2b-8879-e796d618afb5' \
--data-urlencode 'client_uuid=8757fc35-45f9-49aa-b0b0-9b482496d79d' \
--data-urlencode 'passenger_uuid=91122e1b-a36b-4d39-8267-95fe6c5eeb1f' \
--data-urlencode 'trip_uuid=74f3efaf-3363-47f3-94be-eecc1ed1fd89'
```

**End trip**
```bash
curl --location --request POST 'http://localhost:8080/api/v1/trip/end' \
--header 'Authorization: client-key' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'latitude=52.331425' \
--data-urlencode 'longitude=13.531946' \
--data-urlencode 'ride_uuid=efb2eeb1-dfbb-4692-bdcd-3148a93eddab' \
--data-urlencode 'client_uuid=91122e1b-a36b-4d39-8267-95fe6c5eeb0f' \
--data-urlencode 'passenger_uuid=91122e1b-a36b-4d39-8267-95fe6c5eeb1f' \
--data-urlencode 'trip_uuid=1a7c24e1-b9c2-4df4-913e-6fcd05ad125e'
```

## Scalability
Our system is write-heavy. A fast key-value store was chose for that reason. While mongodb has built-in functions to query data between two coordinates, and mongo can handle 
this type of application but, we want to provide the best possible performance. While checking for database benchmarks, I saw that redis' speed was superior to mongo's.

Also, for this demonstration, the data schema was pretty simple, we don't have complex queries. We have one limitation of fetching multiple ZRANGE values. Redis as of now,
does not support retrieving multiple ZRANGE data in one go, but we are fetching the neighbours, 9 in our scenarios. So we spawn 9 goroutines + 1 more for the original geohash
to fetch them concurrently.

With go or similar multithreading supported languages (eg RUST) we can use green-threads which are not as heavy as OS threads, giving as a good performance optimization
out of the box.

Redis also supports clustering. While I did not add any clustering mechanism in this demo, but we can easily add cluster nodes to redis and expose a load balancer beforehand and as well.
We can also make the load balancer a bit more-faster by using a layer 4 load balancer as we don't need to check into the route path/or any use of layer 7 load balancing here (in this demo). 

Secondly, we kept our source of truth as 1. Everything is in 1 data store. So we can easily filter out data and update into one place very easily.

Thirdly, as it is write-heavy with little read operations. We separate read and write replicas. Giving us a more performance boost.

Environments similar to lambda can also be used to deploy this application as itself is stateless.

Another optimization can be achieved by geohash keys. For example, requests from Berlin area should be routed to a datastore that holds the geohash keys to berlin and near by areas.
We can have a static map, giving us `O(1)` lookup. 

Container orchestration with a load balancer should give us a good scalability.

## Note
Some things have been left out. 
- logging
- persistent data store like mongo/mysql or similar
- transactions when multiple write is taking place 
- retry mechanisms
- in multiple places, I've used `fmt.Sprintf()` to format float64 to string, this should have been extracted to a dedicated method, so we have single source of truth
- the docker image is for development purpose only
- `index` can be added to gain a performance boost



## Task
A company called Scootin' Aboot will deploy electric scooters in Ottawa and
Montreal. Design and implement a backend service that exposes a REST-like
API intended for scooter event collecting and reporting to mobile clients.

1. The scooters report an event when a trip begins, report an event when the
trip ends, and send in periodic updates on their location. After beginning a
trip, the scooter is considered occupied. After a trip ends the scooter
becomes free for use. A location update must contain the time, and
geographical coordinates.

2. Mobile clients can query scooter locations and statuses in any rectangular
location (e.g. two pair of coordinates), and filter them by status. While there
will be no actual mobile clients, implement child process that would start
with main process and spawn three fake clients using API randomly (finding
scooters, travelling for 10-15 seconds whilst updating location every 3
seconds, and resting for 2-5 seconds before starting next trip). Client
movement in straight line will be good enough.

3. Both scooters and mobile client users can be identified by an UUID.

4. For the sake of simplicity, both mobile client apps and scooters can
authenticate with the server using a static API key (i.e. no individual
credentials necessary but will most probably be introduced as the project
develops further).
