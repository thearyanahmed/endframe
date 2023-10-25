# NdFrame

The name comes from "and-(another)-framework", a framework for serving restful apis over http using go. Everything exposed via interfaces.

## Running the project
Run the following commands to run the project.

- `git clone git@github.com:thearyanahmed/nordsec.git`
- `cd nordsec`
- `cp .env.example .env`
- Update proper values with `.env`. For this demo `.env` values have be filled in `.env.example`
- `make start` will start the necessary containers. 

## Running tests
While running the container, run `make test` to run the tests. Or to run outside of container run `go test -v ./...` 


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
