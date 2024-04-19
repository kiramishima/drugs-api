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
    env_file:
      - ./.env
    networks:
      - credits
    healthcheck:
      test: [ "CMD-SHELL", "pg_isready -U ${POSTGRES_USER}" ]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    image: api_drugs:v1
    container_name: api
    ports:
      - 3000:3000
    links:
      - database
    env_file:
      - ./.env
    networks:
      - credits

  migrate:
    image: migrate/migrate:latest
    profiles: [ "tools" ]
    volumes:
      - ../migrations:/migrations
    entrypoint:
            [
              "migrate",
              "-path",
              "/migrations",
              "-database",
              "postgres://root:root@database:5432/credits?sslmode=disable",
            ]
    command: [ "up" ]
    depends_on:
      database:
        condition: service_healthy
    restart: on-failure

networks:
  credits:
    driver: bridge
```

Para levantar los servicios ejecute en una terminal `docker compose up -d`

PAra ejecutar la migracion:

```shell
docker compose -f docker-compose.yml --profile tools run --rm migrate up
```
# Deploy en local

- Instalar [golang](https://golang.org/dl)
- Instalar [PostgreSQL](https://www.postgresql.org/)
  - Ejecutar la migración con [migrate](https://github.com/golang-migrate/migrate)
    - Crear la base de datos `credits`
    - Si tiene ya instalado task, ejecutar el comando `task db_up`.
- Instalar [Task CLI](https://taskfile.dev/) para ejecutar las tareas del taskfile.
    - Ejecuta el comando `task run` para levantar el servicio. Default port is 8080

---
## Endpoints

### Endpoint: Credit Assignment

* Path: `/v1/credits/credit-assignment`
* Method: `POST`
* Payload: {investment: integer}
* Respuesta: JSON Response.

Descripción:

Toma una inversión multiplo de 100 y retorna las asignaciones de los prestamos.

Ejemplo respuesta con estatus 200:

```sh
 curl -i -d '{"investment": 6700}' localhost:8080/v1/credits/credit-assignment
```

```json
{"credit_type_300":20,"credit_type_500":0,"credit_type_700":1}
```

Ejemplo respuesta con estatus 400:

```sh
 curl -i -d '{"investment": 6700}' localhost:8080/v1/credits/credit-assignment
```
```json
{"error":"investment needs be multiply of 100"}
```

### Endpoint: Statistics

* Path: `/v1/credits/statistics`
* Method: `GET`
* Response: JSON Response.

Description:

Retorna información general sobre el total de asignaciones realizadas, total de asignaciones exitosas, 
total de asignaciones no exitosas, promedio de inversión exitosa, promedio de inversión no exitosa.

Ejemplo de Respuesta:
```sh
curl -i localhost:8080/v1/credits/statistics
```
```json
{"total_assigns":27,"total_success_assigns":16,"total_fail_assigns":11,"avg_success_assigns":70.08,"avg_fail_assigns":29.92}
```


---

Author: Paul Arizpe

