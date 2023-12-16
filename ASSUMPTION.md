# Assumptions

- We will group cities by proximity in regions.
- Depending on the geohash key length, we will need different number of geohash keys to cover the region. By default, we are using `5` with the length of 4.9km x 4.9km.
- We expect ~100K active drivers in one region.

The amount of QPS and bandwidth we can expect are:

We expect ~100K location updates/writes per second
Each location message consists of

- RideUuid 16 bytes
- PassengerUuid 16 bytes
- ClientUuid 16 bytes
- TripUuid 16 bytes
- Lat, Lon 16 bytes
- Timestamp 8 bytes
- Other metadata 32 bytes

So we get `100,000 * 120 = 120,00,000` or `12MB/s` (almost) bandwidth upload.
For peak our, we can expect to double it and assume `12 * 2 = 24MB/s`

## Non-functional

- Before starting a ride, we assume that the ride goes to the trip origin and starts the trip from there. If the ride is in a bit far away position, the travel point to trip origin is not recorded in the simulation. But endpoint to do so does exist to do so.
