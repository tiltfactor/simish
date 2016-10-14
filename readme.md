# Simish

Simish is a sentence similarity and matching microservice service built using golang.
https://golang.org/


# Getting Started

Installing Simish is simple. Assuming you have golang installed simply run.
```bash
go get github.com/tiltfactor/simish
```

To update to new version add the -u flag
```bash
go get -u github.com/tiltfactor/simish
```

Assuming that a MySQL database is being used as storage you can run
```bash
simish init
```
To generate a configuration file. See Databases section for information on extending support for
other databases.

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

Once this file is generated
```bash
simish start
```
will start the simish service. To run it as a background process:
```bash
nohup simish start &
```

# Databases
Simish comes with support for MySQL but can be easily extended to support other
databases and storage methods.

To extend Simish with support for other storage systems simply implement the InputOutputStore
interface (found in domain/InputResponse.go).
```
type InputOutputStore interface {
	SaveInputOutput(InputOutput) error
	Response(string, int64) (InputOutput, float64)
}
```
*If you are extending the features supported by the storage (eg adding upvotes and downvotes storage) please be sure to also extend the InputOutputStore interface. This will make sure that the contract stays the same between storage backends and that the main program does not have to worry about the implementation details of the storage mechanism* 

# Algorithm
To change the soft match algorithm edit the SoftMatch function in the domain/InputResponse file.
Currently the service uses the JaroWinklerDistance to calculate the closeness of two sentences
but this can easily be changed or extended by editing the SoftMatch function in
domain/InputResponse.go


# Getting a match
## Request
```bash
GET http://localhost:8000/api/v1/response?input=Hello&room=1
```

## Response
```json
{
	"input": "Hello",
	"response": "#splat",
	"match": "Hello",
	"room": "1",
	"score":1
}
```
The response consists of the provided input, the response that was found for that input, the existing
input that was found for it, the room number used for searching, and the score (closeness) of the two
inputs.
