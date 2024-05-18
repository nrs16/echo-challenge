# Coding Challenge task

## Disclaimer
- This is the first time I use echo framework, so I played around a bit and added some middleware functions for my own curiosity, just to see what functionalitites it has and what I can do with it beyond the scope of the task
- The additional middleware functions do not affect the functionality of the code and the task. You should be able to test the task normally and get the expected output

## How to get and run the code
### Through github

- clone the project using ```git clone https://github.com/nrs16/echo-challenge.git```
- run the command ```go mod download```
- run the command ```go run main.go```
- use this curl to test the code: 

```
curl --location 'http://localhost:8080/routes' \
--header 'x-correlation-id: jweygfjkegdf' \
--header 'Content-Type: application/json' \
--data '[["LAX","DXB"],["JFK","LAX"], ["SFO","SJC"], ["DXB","SFO"]]'

```
You can remove x-correlation-id and Content-Type headers but you must send the body


### Through docker

- get the image using ```git pull nrs16/echoserver```
this might take a while to download because it has golang:1.21 image as base

- run the image inside a container using ```docker run -p 8000:8080 --name routes nrs16/echoserver```
    - Note that you can change port 8000 to whichever port you want on your machine.
    - To run the container in the background add -d to the commad so ```docker run -d -p 8000:8080 --name routes nrs16/echoserver```
- use the below curl to test the code:
```
curl --location 'http://localhost:8000/routes' \
--header 'x-correlation-id: jweygfjkegdf' \
--header 'Content-Type: application/json' \
--data '[["LAX","DXB"],["JFK","LAX"], ["SFO","SJC"], ["DXB","SFO"]]'

```
You can remove x-correlation-id and Content-Type headers but you must send the body
