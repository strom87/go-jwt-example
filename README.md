# Go jwt example

A simple example on how to use jwt (Json Web Tokens) with go 

### Tutorial
Drag down the project
```sh
$ go get github.com/strom87/go-jwt-example
```
Run the main.go file to start the server.
```sh
$ go run main.go
```
The site listens to port 1337.  
There is two routes.
* [localhost:1337/api](http://localhost:1337/api)
* [localhost:1337/token](http://localhost:1337/token)

The /api route is protected and needs a valid token to be accessed.  
To test if it works or not, run the following curl commands in your terminal window:
```sh
$ curl localhost:1337/api
```
Now you see the message "Not authenticated, route protected", we need to get a token before we can access the route.  
Run the following command:
```sh
$ curl localhost:1337/token
```
A token is returned, copy this token and replace the INSERT_TOKEN part and run the following command:
```sh
$ curl -H "Authorization: Bearer INSERT_TOKEN" localhost:1337/api
```
Now you should have access to the api.

## Libraries used
* [https://github.com/strom87/middle] (https://github.com/strom87/middle)
* [https://github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
