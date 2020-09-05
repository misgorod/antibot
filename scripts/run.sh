#!/bin/bash

for i in $(seq 1 $1); do
    curl -H 'X-Forwarded-For: 123.255.1.2' http://localhost:8080/api
done