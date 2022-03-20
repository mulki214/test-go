# READ ME

## Requirement
- go 1.18
- echo 4.7.2

## Installation
- go mod tidy
- go run server.go

## API URL

Url: localhost:8081/login \
Method: POST \
Form-Data \
username: admin \
password: admin

Url: localhost:8081/api/competition \
Method: GET \
Header \
Autorization: Bearer {TokenJWT} \
queryString \
competition: UEFA Champions League \
page: 1

Url: localhost:8081/api/schedule \
Method: GET \
Header \
Autorization: Bearer {TokenJWT} \
queryString \
team: Chelsea \
year: 2011 \
page: 1


