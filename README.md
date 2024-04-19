# Drugs API

-----

# **Stack**

- Router: [Chi 🚀](https://github.com/go-chi/chi)
- Logger: [Zap ⚡](https://github.com/uber-go/zap)
- Mocks: [gomock 💀](https://github.com/uber-go/mock)
- Asserts: [testify 💀](https://github.com/stretchr/testify)
- DI: [fx 🤖](https://github.com/uber-go/fx)
- Deploy: [Docker 🐳](https://www.docker.com)
- Database: [PostgreSQL](https://www.postgresql.org/)

## **Deploy con Docker**
- El archivo `scripts/build-container.sh` contiene las instrucciones para construir el contenedor.
- El archivo `scripts/run-container.sh` para ejecutar el contenedor de la imagen construida con `scripts/build-container.sh`.
- Por default el puerto es `8080`.
- Para cambiar el puesto, proporcione la variable de entorno `PORT`.

## **Variables de entorno**

El archivo dev.env contiene las variables que se pueden configurar en el archivo `.env`

```sh
# App
APP_SECRET=B@nk4I
# Database
DATABASE_DRIVER=pgx
DATABASE_URL=postgres://postgres:123456@192.168.100.47/ionix
DATABASE_MAX_OPEN_CONNECTIONS=25
DATABASE_MAX_IDDLE_CONNECTIONS=25
DATABASE_MAX_IDDLE_TIME=15m
# HTTP
HTTP_SERVER_IDLE_TIMEOUT=60s
PORT=8080
HTTP_SERVER_READ_TIMEOUT=1s
HTTP_SERVER_WRITE_TIMEOUT=2s
# JWTF
JWT_PRIVATE_KEY=RacconCity
TOKEN_TTL=300
# Context
CONTEXT_TIMEOUT=10

# Postgres
POSTGRES_DBNAME=ionix
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
```

````shell
# Crear una red
docker network create pg-network

# Levantar la DB
docker run -d \
    -e POSTGRES_USER="root" \
    -e POSTGRES_PASSWORD="root" \
    -e POSTGRES_DB="ionix" \
    -p 5432:5432 \
    --network=pg-network \
    --name pg-database \
    postgres:latest
    
# Ejecutar la migracion
task dbUp

# Ejecutar el servicio, en el puerto 3000
docker run -it \
    --network=pg-network \
    -p 3000:3000 \
    --env-file .env \
    api_drugs:v1 
````

## **Deploy con Docker Compose**
Para desplegar el api, genere el archivo `docker-compose.yml` y dentro pegue lo siguiente:

```yaml
version: "3.9"
services:
  database:
    image: postgres:latest
    container_name: database
    hostname: database
    ports:
      - 5432:5432
    volumes:
      - "./postgres_data:/var/lib/postgresql/data:rw"
      - "./init-db.sh:/docker-entrypoint-initdb.d/init-user-db.sh"
    env_file:
      - ./.env
    networks:
      - drugs
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready -U ${POSTGRES_USER}" ]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    image: kiramishima/api_drugs:v1
    build:
      context: .
      dockerfile: Dockerfile
    container_name: api
    ports:
      - 8080:8080
    links:
      - database
    env_file:
      - ./.env
    networks:
      - drugs

  migrate:
    image: migrate/migrate:latest
    profiles: [ "tools" ]
    volumes:
      - ./migrations:/migrations
    command: ["-path=/migrations", "-database=postgres://postgres:postgres@database:5432/ionix?sslmode=disable", "up", "3"]
    depends_on:
      database:
        condition: service_healthy
    restart: on-failure
    networks:
      - drugs

networks:
  drugs:
    driver: bridge
```

Para levantar los servicios ejecute en el siguiente orden en una terminal 

1. Levantar todos los servicios

`docker compose up -d`

2. Ejecutar la migracion:

```shell
docker compose up migrate -d
```

3. Levantar el servicio del API

````shell
docker compose up api -d
````

# Deploy en local

- Instalar [golang](https://golang.org/dl)
- Instalar [PostgreSQL](https://www.postgresql.org/)
  - Ejecutar la migración con [migrate](https://github.com/golang-migrate/migrate)
    - Crear la base de datos `ionix`
    - Si tiene ya instalado task, ejecutar el comando `task dbUp`.
- Instalar [Task CLI](https://taskfile.dev/) para ejecutar las tareas del taskfile.
    - Ejecuta el comando `task run` para levantar el servicio. Default port is 8080

---
## Endpoints

### **AUTH**
#### Endpoint: Auth/sign-in

* Path: `/v1/auth/sign-in`
* Method: `POST`
* Payload: `{email: string|email|required, password: string|required}`
* Respuesta: JSON Response.

Descripción:

Autenticar al usuario

Ejemplo respuesta con estatus 200:

```sh
curl localhost:8080/v1/auth/sign-in -d '{"email": "giny@mail.com", "password": "12356"}'
```

```json
{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM0ODc1MjksIm5iZiI6MTcxMzQ4NzIyOSwiaWF0IjoxNzEzNDg3MjI5LCJqdGkiOiIxIn0.6DDctPomEGkVlOtW6QQWVZbAez2HifBjcAOXtbwBmw8"}
```

Ejemplo respuesta con estatus 400:

```sh
curl localhost:8080/v1/auth/sign-in -d '{"email": "giny@mail.com", "password": ""}'
```
```json
{"error":"error message"}
```

#### Endpoint: Auth/sign-up

* Path: `/v1/auth/sign-in`
* Method: `POST`
* Payload: `{email: string|email|required, password: string|required, name: string}`
* Respuesta: JSON Response.

Descripción:

Registrar nuevo usuario

Ejemplo respuesta con estatus 200:

```sh
curl -X POST localhost:8080/v1/auth/sign-up -d '{"email": "giny@mail.com", "password": "12356", "name": "Gina"}'
```

```json
{"message":"Registro exitoso."}
```

Ejemplo respuesta con estatus 400:

```sh
curl localhost:8080/v1/auth/sign-in -d '{"email": "giny@mail.com", "password": ""}'
```
```json
{"error":"error message"}
```


### **Drugs**
#### Endpoint: /v1/drugs

* Path: `/v1/drugs`
* Method: `GET`
* Auth: **JWT Token**
* Respuesta: JSON Response.

Descripción:

Obtiene el listado de drugs

Ejemplo respuesta con estatus 200:

```sh
curl localhost:8080/v1/drugs -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM0ODc5NTMsIm5iZiI6MTcxMzQ4NzY1MywiaWF0IjoxNzEzNDg3NjUzLCJqdGkiOiIxIn0.0kjZyvSmswM36lUsZdUunTSFClNj8y8NqawK24bj_Qc"
```

```json
{
"data":[
    {
      "id":2,
      "name":"Cafiaspirina",
      "approved":true,
      "min_dose":1,
      "max_dose":4,
      "available_at":"2024-05-15T12:00:00Z"
    }
  ]
}
```

Ejemplo respuesta con estatus 400:

```sh
curl localhost:8080/v1/drugs -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM0ODc5NTMsIm5iZiI6MTcxMzQ4NzY1MywiaWF0IjoxNzEzNDg3NjUzLCJqdGkiOiIxIn0.0kjZyvSmswM36lUsZdUunTSFClNj8y8NqawK24bj_Qc"
```

```json
{"error":"error message"}
```

#### Endpoint: /v1/drugs

* Path: `/v1/drugs`
* Method: `POST`
* Payload: `{name: string|required, approved: string|required, min_dose: integer|required}, max_dose: integer|required, available_at: string|datetime|required`
* Respuesta: JSON Response.

Descripción:

Registrar nuevo drug

Ejemplo respuesta con estatus 200:

```sh
curl localhost:8080/v1/drugs \ 
-H "Authorization: Bearer <JWT TOKEN>" \
-d '{"name": "cafiaspirina", "approved": true, "min_dose": 1, "max_dose": 4, "available_at": "2024-05-05 13:50:00"}'
```

```json
{"message":"Se ha registrado el nuevo medicamento de manera exitosa"}
```

Ejemplo respuesta con estatus 400:

```json
{"error":"error message"}
```

#### Endpoint: /v1/drugs/{id}

* Path: `/v1/drugs/{id}`
* Path Param:
  * id: integer
* Method: `PUT`
* Payload: `{name: string|required, approved: string|required, min_dose: integer|required}, max_dose: integer|required, available_at: string|datetime|required`
* Respuesta: JSON Response.

Descripción:

Actualizar un registro de drug

Ejemplo:

```sh
curl -X PUT localhost:8080/v1/drugs/1 \ 
-H "Authorization: Bearer <JWT TOKEN>" \
-d '{"name": "Aspirina"}'

curl -X PUT localhost:8080/v1/drugs/2 \
-H "Authorization: Bearer <JWT TOKEN>" \ 
-d '{"available_at": "2024-05-15 12:00:00", "name": "Cafiaspirina"}'
```

Ejemplo respuesta con estatus 200:

```json
{"message":"Se ha actualizado la información del medicamento de manera exitosa"}
```

Ejemplo respuesta con estatus 400:

```json
{"error":"error message"}
```

#### Endpoint: /v1/drugs/{id}

* Path: `/v1/drugs/{id}`
* Path Param:
  * id: integer
* Method: `DELETE`
* Respuesta: JSON Response.

Descripción:

Eliminar un registro de drug

Ejemplo

```sh
curl -X DELETE localhost:8080/v1/drugs/3 \
-H "Authorization: Bearer <JWT TOKEN>"
```

Ejemplo respuesta con estatus 200:

```json
{"message":"Se ha eliminado el medicamento de manera exitosa"}
```

Ejemplo respuesta con estatus 400:

```json
{"error":"error message"}
```

//
### **Vaccinations**
#### Endpoint: /v1/vaccination

* Path: `/v1/vaccination`
* Method: `GET`
* Auth: **JWT Token**
* Respuesta: JSON Response.

Descripción:

Obtiene el listado de drugs

Ejemplo respuesta con estatus 200:

```sh
curl localhost:8080/v1/vaccination \
-H "Authorization: Bearer <JWT TOKEN>"
```

```json
{
"data":[
    {
      "id":2,
      "name":"Jhone Doe",
      "drug":"Cafiaspirina",
      "drug_id":2,
      "dose":1,
      "date":"2024-05-05T13:50:00Z"
    }
  ]
}
```

Ejemplo respuesta con estatus 400:

```json
{"error":"error message"}
```

#### Endpoint: /v1/vaccination

* Path: `/v1/vaccination`
* Method: `POST`
* Payload: `{name: string|required, drug_id: integer|required, dose: integer|required}, applied_at: string|datetime|required`
* Respuesta: JSON Response.

Descripción:

Registrar nuevo vaccination

Ejemplo respuesta con estatus 200:

```sh
curl localhost:8080/v1/vaccination \ 
-H "Authorization: Bearer <JWT TOKEN>" \
-d '{"name": "Jhone Doe", "drug_id": 2, "dose": 1, "applied_at": "2024-05-05 13:50:00"}
```

```json
{"message":"Se ha registrado de manera exitosa"}
```

Ejemplo respuesta con estatus 400:

```json
{"error":"error message"}
```

#### Endpoint: /v1/vaccination/{id}

* Path: `/v1/vaccination/{id}`
* Path Param:
  * id: integer
* Method: `PUT`
* Payload: `{name: string|required, drug_id: integer|required, dose: integer|required}, applied_at: string|datetime|required`
* Respuesta: JSON Response.

Descripción:

Actualizar un registro de vaccination

Ejemplo:

```sh
curl -X PUT localhost:8080/v1/vaccination/1 \ 
-H "Authorization: Bearer <JWT TOKEN>" \
-d '{"name": "Raffaella Carra"}'

curl -X PUT localhost:8080/v1/vaccination/2 \
-H "Authorization: Bearer <JWT TOKEN>" \ 
-d '{"applied_at": "2024-05-15 12:00:00", "name": "Raffaella Carra"}'
```

Ejemplo respuesta con estatus 200:

```json
{"message":"Se ha actualizado la información de manera exitosa"}
```

Ejemplo respuesta con estatus 400:

```json
{"error":"error message"}
```

#### Endpoint: /v1/vaccination/{id}

* Path: `/v1/vaccination/{id}`
* Path Param:
  * id: integer
* Method: `DELETE`
* Respuesta: JSON Response.

Descripción:

Eliminar un registro de vaccination

Ejemplo

```sh
curl -X DELETE localhost:8080/v1/vaccination/3 \
-H "Authorization: Bearer <JWT TOKEN>"
```

Ejemplo respuesta con estatus 200:

```json
{"message":"Se ha eliminado el registro de manera exitosa"}
```

Ejemplo respuesta con estatus 400:

```json
{"error":"error message"}
```
---

Author: Paul Arizpe

