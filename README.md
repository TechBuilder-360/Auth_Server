# Auth_Server
This application is responsible for authenticating, and validating user token.


### How to set up project
1. Clone project: `git clone https://github.com/TechBuilder-360/Auth_Server.git`
1. Run `go get .`
1. Run `go build . && go run main.go`


### How does authentication works
1. User clicks an auth button, then they get redirected to an auth page that has a form for a username, password and social sign in.
2. Depending on the response:
    1. Success - return token to user, closes browser tab.
    2. Failure - display error message



### Server flow
1. Register user
    1. email and password (encrypted password)
    2. Social authentication (pass social data to portal server)
1. Authenticate and generate token
    1. Get email and password from user
    1. Social authentication 
1. Generate token from refresh token
2. Validate token
1. Verify client ID and request hash in the middleware.