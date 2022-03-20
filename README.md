# READ ME

## Requirement
- go 1.18
- echo 4.7.2

## Installation
- go mod tidy
- go run server.go

## API URL

url: localhost:8081/login \
method: POST \
form-data \
username: admin \
password: admin

url: localhost:8081/api/competion \
method: GET \
Header \
Autorization: Bearer {TokenJWT} \
queryString \
competition: UEFA Champions League \
page: 1

url: localhost:8081/api/schedule \
method: GET \
Header \
Autorization: Bearer {TokenJWT} \
queryString \
team: Chelsea \
year: 2011 \
page: 1


