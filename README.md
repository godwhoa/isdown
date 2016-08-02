isdown
======

A service that let's you test if a site is down.


## Architecture

the architecture is all distributed. in-theory multiple people could host the "minion" nodes.

### Master server:
+ Maintains a list of minion nodes.
+ Checks minion's health every 1sec.
+ `/register` for the minions to register. (POST)
+ `/list` gives list of minion in json for client to use. (GET)

### Minion server:
+ Registers with the master server.
+ `/health` for master srv. to do health checks. (GET)
+ `/isdown` for sending in sites to test. (POST)

## TODO
+ ~~Make no assumption about which port minions run on~~
+ Client