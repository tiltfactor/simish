# Simish

Simish is a sentence matching microservice built for the crowded-dungeon game.


```bash
go get github.com/tiltfactor/simish

# will run on port 8765
simish
```

# Getting a match
```bash
http://localhost:8000/api/v1/response?input=Hello&room=ExampleRoom
```

## WIP - Not currently working
# Upvoting a match
```bash
http://localhost:8000/api/v1/upvote?input=Hello&match=Hey&room=ExampleRoom
```

# Downvoting a match
```bash
http://localhost:8000/api/v1/downvote?input=Hello&match=Hey&room=ExampleRoom
```
