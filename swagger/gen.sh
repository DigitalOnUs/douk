#!/usr/bin/env bash
# creating clients from spec
java -jar swagger-codegen-cli-2.2.1.jar generate -i swagger.yaml -l go -o ./api -c ./go-config.json
java -jar swagger-codegen-cli-2.2.1.jar generate -i swagger.yaml -l go -o ./api-js
