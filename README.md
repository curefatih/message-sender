<h1 align="center">📬 message-sender</h1>

Basic Message sender is an app that tries to 
send message for every N minute per unprocessed X message.

## Run

### Create `webhook.site` that returns a message

* Go to [https://webhook.site/](https://webhook.site/)

* Create url with given configuration by clicking the top left `Edit` button.

![Creating webhook site url](/assets/creating-webhooksite-response.png)

|  Column | Value | 
|---|---|
| Status Code  | 202 | 
| Content type | `application/json` | 

to create a random id with message field - use this payload

```json
{
  "message": "Accepted",
  "messageId": "$request.uuid$"
}
```

### Setting environment variables

Copy `config.example.yml` as `config.yml` and set the variables inside.

By default postgresql and redis credentials are in

* `docker/local.postgresql.env` and, 
* `docker/local.redis.env`

You can also set other environment variables by checking the 

* https://hub.docker.com/r/bitnami/redis for the redis image
* https://hub.docker.com/_/postgres for the postgresql

but beware that some configs also need code changes.

### Compose up

To run the entire application in detach mode

```sh
docker compose -f docker/docker-compose.local.yml up -d
```

if you just want to stop

```sh
docker compose -f docker/docker-compose.local.yml stop
```

and if you want to cleanup 

```sh
docker compose -f docker/docker-compose.local.yml down
```

## Usage

There is a swagger that can be accessible with link http://localhost:8080/swagger/index.html#/ for the local setup.

Also, for the database, there is an adminer container that can be 
accesible with link http://localhost:8082/

## Limitations

For every iteration, we need to make sure that the processing rate is more or at least the same as the creation rate. 

Initially, it is designed for 2 messages per 2 minutes. Higher limits might block the entire application.
