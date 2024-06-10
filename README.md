# Golang Gateway Service Assignment

### Running in your marchine

Pull the repository and go to the main folder. 
Rename `.env.example` with `.env`

For the assignement I used docker in order to deploy and test the application. 
To run the application, please run this command in the root folder: 
```sh
docker compose up -d
```
Once done, visit: http://localhost:8080/swagger/index.html

### Users

I made sure to inject some users for test, so please login into one of the accounts : 

You can login throw a POST request http://localhost:8080/login
```json
{
  "username": "username_basic"
  "password": "password",
}
```
The response will generate a token, inject it in your Authorization Header.

This user is subscribe to a Basic plan which has 5 credits.
To test the credits limit, i have created 2 external services: 

- Geolocation service
- Bearer service

Each request for an external service will deduct a credit from the user connected.
Each request for an external service is cached in redis
To run tests, please run the following command in the root folder:
```sh
go test ./tests/*
```

External service structure:

```
│   service-name
│   ├── authenticator.go
│   ├── errors.go
│   ├── init.go
│   └── service.go
```
Once the service is added, you need to add in in `external/services.yaml`.

The `init.go` file is reserved to the depency injection 
```
    service := NewGeolocationService()
	external.RegisterService("service-name", service)
```

Make sure that your service name is the same as configured in the `services.yaml`

Each service has 3 main methods interfaced :
`Do` : Needed to call the external service endpoint
`Authenticate` : Needed to authenticate the service
`GetCacheDuration` : Get the cache ms parameter configured in `external/services.yaml`.

In order to add enable the external service package, go to `main.go` and add your import as the following: 

`_ "assignment/external/service-name"`

Once an external service added, please run the following command:
```
docker compose up -d --build
```
