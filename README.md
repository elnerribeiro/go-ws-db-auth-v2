# go-ws-db-auth-v2

This is a version 2 Proof of Concept of using Golang in a somewhat standard commercial program.
After so many years, Golang have evolved a lot - so let's check how that changes the code.
Here we have some REST web services with Authentication/Authorization using JWT, database connection with Postgres, a persistence layer that uses structs to abstract table models and some examples of sync/async calls.

UPDATED THINGS:
    
    1. All running with docker compose
    2. Better logging, using zerolog - JSON format
    3. Using context for logging and DB connections
    4. Postgres connection done using pgx
    5. Connection pool
    6. Updated libraries, updated Go version


Build and run (with docker):
    
    docker build -t gopoc .
    docker compose up -d

Stop (with docker):

    docker compose down

Build and run (no docker):

	(Linux)
	go build -o main
	./main

	(Windows)
	go build -o main.exe
	.\main.exe

Default URL:

	http://localhost:8000

Postman Collection:

	https://www.getpostman.com/collections/06704a4c68b44e63502e

Calls:

	/api/login (POST)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Basic SldUcGFzc3dvcmQxMjNA
			Body:
				{
					"email":"user@user.com","password":"DDAF35A193617ABACC417349AE20413112E6FA4E89A97EA20A9EEEE64B55D39A2192992A274FC1A836BA3C23A3FEEBBD454D4423643CE80E2A9AC94FA54CA49F"
				}
            Or Body:
                {
                    "email":"admin@admin.com","password":"3C9909AFEC25354D551DAE21590BB26E38D53F2173B8D3DC3EEE4C047E7AB1C1EB8B85103E3BE7BA613B31BB5C9C36214DC9F14A42FD7A2FDB84856BCA5C44C2"
                }
		Response:
			{
				"account": {
					"id": 1,
					"email": "usuario@usuario.com",
					"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsIlJvbGUiOiJ1c2VyIiwiZXhwIjoxNTg1MTI0ODIwfQ.d5GsE9WDWlbRxQRtfuAO-G0SFnYV1ZhAs-m5rb1t--E",
					"role": "user"
				},
				"message": "Logged In",
				"status": true
			}

	/api/validate (GET)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "success",
				"role": "user",
				"status": true,
				"userId": 1
			}
	
	/api/users (POST) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "Unauthorized user",
				"status": false
			}
			ou
			{
				"data": [
					{
						"id": 1,
						"email": "user@user.com",
						"role": "user"
					},
					{
						"id": 2,
						"email": "admin@admin.com",
						"role": "admin"
					}
				],
				"message": "success",
				"status": true
			}

	/api/user/{id} (GET) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 2,
					"email": "admin@admin.com",
					"role": "admin"
				},
				"message": "success",
				"status": true
			}

	/api/user (PUT) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
			Body:
				{
					"email":"elner.ribeiro@gmail.comx",
					"id": 9, //optional, only used for updates
					"role": "admin",
					"password":"3C9909AFEC25354D551DAE21590BB26E38D53F2173B8D3DC3EEE4C047E7AB1C1EB8B85103E3BE7BA613B31BB5C9C36214DC9F14A42FD7A2FDB84856BCA5C44C2"
				}
		Response:
			{
				"data": {
					"id": 9,
					"email": "elner.ribeiro@gmail.comx",
					"role": "admin"
				},
				"message": "success",
				"status": true
			}

	/api/user/{id} (DELETE) - Only role admin
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "success",
				"status": true
			}

	/api/insert (DELETE)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"message": "success",
				"status": true
			}

	/api/insert/{id} (GET)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 6,
					"type": "sync",
					"quantity": 10,
					"status": "Finished",
					"list": [
						{
							"id": 120107,
							"id_ins_id": 6,
							"pos": 1
						},
						...
					]
				},
				"message": "success",
				"status": true
			}

	/api/insert/sync/{quantity} (PUT)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 6,
					"type": "sync",
					"quantity": 10,
					"status": "Finished",
					"tstampinit": 1585175306,
        			"tstampend": 1585175306,
					"list": [
						{
							"id": 120107,
							"id_ins_id": 6,
							"pos": 1
						},
						{
							"id": 120108,
							"id_ins_id": 6,
							"pos": 2
						},
						{
							"id": 120109,
							"id_ins_id": 6,
							"pos": 3
						},
						{
							"id": 120110,
							"id_ins_id": 6,
							"pos": 4
						},
						{
							"id": 120111,
							"id_ins_id": 6,
							"pos": 5
						},
						{
							"id": 120112,
							"id_ins_id": 6,
							"pos": 6
						},
						{
							"id": 120113,
							"id_ins_id": 6,
							"pos": 7
						},
						{
							"id": 120114,
							"id_ins_id": 6,
							"pos": 8
						},
						{
							"id": 120115,
							"id_ins_id": 6,
							"pos": 9
						}
					]
				},
				"message": "success",
				"status": true
			}

	/api/insert/async/{quantity} (PUT)
		Request:
			Headers:
				Content-Type: application/json
				Authorization: Bearer {{token}}
		Response:
			{
				"data": {
					"id": 5,
					"type": "async",
					"quantity": 20000,
					"status": "Running",
					"tstampinit": 1585174744
				},
				"message": "success",
				"status": true
			}
