# GoRateLimiterFC
This work is about a implementation of a Rate Limiter that should limit a number of requests according to the restrictions of number the requester per second and the limit of time of each token or ip.

## Instructions to project

1. To run this project you need docker installed in you environment.
2. If necessary to configure the environment variables of project, you can access in path **cmd/server/.env**
   
   2.1. Inside this env file has the token and ip configuration, if you want to change the execution parameters.

### Instructions to execute project

To run this project you need have docker/docker-compose installed, after this just execute the commands bellow:

1. **docker-compose build** - Build the images of redis and app.
2. **docker-compose up** - Run the application on localhost:8080

## Instructions to middleware

The implementation of the rete limiter is done through middleware, then is easily to attach to your system. Bellow are instructions for
using the rate limiter.

1. Import the depency of middleware in your project.
  `"github.com/GoExpertCurso/GoRateLimiterFC/internal/infra/web/middleware" `
2. You must have your own router or a third-party router in your code.
3. Install the middleware in you code like the example below:
   ` wrappedMux := mid.RateLimitMiddleware(mux, dbClient, conf) `
4. You have to pass three parameters to middleware, router, dbClient and configurations.
5. When you up you server adding the middleaware:
   ` http.ListenAndServe(":8080", wrappedMux) `