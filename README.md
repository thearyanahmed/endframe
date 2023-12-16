# Endframe

This is a framework for building go based applications. Current it supports HTTP REST apis using [Chi](https://github.com/go-chi/chi).

Everything else is bound with interface. Database to handlers to services. Configurable from the outside.
This also comes with docker, docker compose. The default docker file also keeps watch for files changes so it can
handle the automatic code rebuild.

Checkout the `reframe` repository that has a working demo for a location based proximity app.

## Running the project

Run the following commands to run the project.

- `git clone git@github.com:thearyanahmed/endframe.git`
- `cd endframe`
- `cp .env.example .env`
- Update proper values with `.env`. For this demo `.env` values have be filled in `.env.example`
- `make start` will start the necessary containers.

## Running tests

While running the container, run `make test` to run the tests. Or to run outside of container run `go test -v ./...`

```txt
┌──────────┐    ┌──────────┐    ┌───────────┐    ┌────────────────────────────────────────┐
│          │    │          │    │           │    │           Request Serializer           │
│  Server  ├────►  Router  ├────►  Handler  ├────►────────────────────┬───────────────────┼───┐
│          │    │          │    │           │    │       Serializer   │     Validator     │   │
└──────────┘    └──────────┘    └───────────┘    └────────────────────┴───────────────────┘   │
                                                                                              │
     ┌────────────────────────────────────────────────────────────────────────────────────────┘
     │
┌────▼────┐                                       ┌─────────────────┐
│ Service ├─────────────────────────────────────► │     Response    │
│         │                                       └─────────────────┘
└────↑────┘ ┌──────────────────────┐                 ┌─────────────────────┐
     │      │                      │                 │                     │
     └──────►  Other Service       │←───────│──────► │ Repository (Redis ) │
     │      │                      │        │        │                     │
     │      └──────────────────────┘        │        └─────────────────────┘
     │                                      │
     │ ─────────────────────────────────────│
```

