Sometimes it's useful to wait until PostgreSQL has finished coming up before
trying to do stuff to it, for example when using Kubernete's `initContainers`
feature. AwaitPostgres is a lightweight docker image that will repeatedly
attempt to connect to a PostgreSQL instance, and exit when it succeeds (or exit
with nonzero status if it can't).

Example usage:
1) start a vanilla postgres instance:
docker run -d --name postgres  --publish 5432:5432 postgres:9.5
2) connect to it!
docker run --env 'POSTGRES_URL=postgresql://postgres:@postgres:5432/postgres?sslmode=disable' --link postgres:postgres bowerswilkins/awaitpostgres

Connection is established via the required environment variable POSTGRES_URL,
which specifies the instance to connect to as well as the user, password, etc.

Other environment variables are
* RETRIES: the number of times to retry connecting to postgres (defaults to 10)
* WAIT_SECS: the number of seconds to wait between connection attempts (defaults to 2)
