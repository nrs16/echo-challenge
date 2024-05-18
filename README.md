# Coding Challenge task
## How to get and run the code
### Through github
- clone the project using ```git clone https://github.com/nrs16/echo-challenge.git```
- go run main.go
- use this curl to test the code: 
```
curl --location 'http://localhost:8080/routes' \
--header 'x-correlation-id: jweygfjkegdf' \
--header 'Content-Type: application/json' \
--data '[["LAX","DXB"],["JFK","LAX"], ["SFO","SJC"], ["DXB","SFO"]]'

```
You can remove x-correlation-id and Content-Type headers but you mist send the body


### Through docker