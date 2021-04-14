#!/bin/bash

# A collection of simple curl requests that can be used to manually test endpoints before and while writing automated tests

curl localhost:9090/signin -X POST -d '{"username":"sickboy", "password":"ubi123"}'
curl localhost:9090/signup -X POST -d '{"username":"gamer", "firstName":"Ubi", "lastName":"vius", "email":"ubivius@udes.com", "password":"test123"}'
