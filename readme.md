# Simish

Simish is a sentence matching microservice built for the crowded-dungeon game.

To install Simish just run.
```bash
go get github.com/tiltfactor/simish
```

To update to new version add the -u flag
```bash
go get -u github.com/tiltfactor/simish
```

To run the program make sure you have a db_cfg.json file in the directory where you're running
Simish.

```json
{
  "username": "db_username",
  "password": "db_password",
  "ip_addr": "db_ipAddress eg 127.0.0.1",
  "database": "db_name",
  "db_port": "9000",
  "server_port": "9000"
}
```

and then run
```
simish
```

To change the soft match algorithm edit the SoftMatch function in the domain/InputResponse file.


# Getting a match
```bash
GET http://localhost:8000/api/v1/response?input=Hello&room=ExampleRoom
```
