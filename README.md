## Useful links
- [https://redis.io/commands/geosearch/](https://redis.io/commands/geosearch/)
- [https://github.com/redis/redis/issues/4417#issuecomment-696256741](https://github.com/redis/redis/issues/4417#issuecomment-696256741)
- [https://www.memurai.com/blog/geospatial-queries-in-redis](https://www.memurai.com/blog/geospatial-queries-in-redis)

## My Thoughts 
**Using Kafka**

For the server side app, I want to use kafka. As it is append only log, it is create a series of events out of the box. So, we will be able to plot the location(s) of scooters in an orderly sequence.

**Namespacing/Partioning**

Having proper namespcaing / partioning based on location, our loadbalancers can improve significantly in performance. Imagine this theory, suppose we have a 10x10 KM area, where 1x1KM is Name Area1 to AreaN. 

If we could determine the current location of the mobile client (eg: Area5), and in our data store, we namespace/ prefix the records in `Partion Area5` : Append only log all the scooters that are currently in that area. So scooterID, lat, long etc. 

That way our load balancer / service can simply route the request to that specific partion, giving us performance boost and less resoruce is used. The trade off would be simpy to add logical elements in the code to make sure we push to the right partion and keep in mind a single trip can be in multiple partions (not at the same time) when the rider crosses **AreaX** to **AreaY**. 

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

