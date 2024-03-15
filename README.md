# GoRateLimiterFC
This work is about a implementation of a Rate Limiter that should limit a number of requests according to the restrictions of number the requeser per second and the limit of time of each token or ip.

## Instructions to project

To run this project you need docker installed in you envroment.

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