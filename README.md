## Useful links
- [https://redis.io/commands/geosearch/](https://redis.io/commands/geosearch/)
- [https://github.com/redis/redis/issues/4417#issuecomment-696256741](https://github.com/redis/redis/issues/4417#issuecomment-696256741)
- [https://www.memurai.com/blog/geospatial-queries-in-redis](https://www.memurai.com/blog/geospatial-queries-in-redis)

## My Thoughts 
We want to use redis with geohash. Given an 

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


## Area Coordinates
```
┌──────────────────────────┐  
│(x4,y4)            (x3,y3)│  
│                          │
│(x1,y1)            (x2,y2)│    
└──────────────────────────┘    
```

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
