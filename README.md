Sometimes it's useful to wait until PostgreSQL has finished coming up before
trying to do stuff to it, for example when using Kubernete's `initContainers`
feature. AwaitPostgres is a lightweight docker image that will repeatedly
attempt to connect to a PostgreSQL instance, and exit politely when it succeeds (or exit
with nonzero status if it can't).

### Basic Example Usage:
1) start a vanilla postgres instance:
```
docker run -d --name postgres  --publish 5432:5432 postgres:9.5
```
2) connect to it!
```
docker run --env 'POSTGRES_URL=postgresql://postgres:@postgres:5432/postgres?sslmode=disable' --link postgres:postgres bowerswilkins/awaitpostgres
```

Connection is established via the required environment variable POSTGRES_URL,
which specifies the instance to connect to as well as the user, password, etc.

Other environment variables are
* `RETRIES`: the number of times to retry connecting to postgres (defaults to 10)
* `WAIT_SECS`: the number of seconds to wait between connection attempts (defaults to 2)


### Example Usage With Kubernetes' initContainers
This repo includes some example code that demonstrates using AwaitPostgres in conjunction with Kubernetes. Kubernetes has an
[initContainers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) feature that allows you to execute some
code in a separate Docker image before running the deployment or job you're actually interested in. Let's say we have two pods
we want Kubernetes to run- one that runs a Postgres container, and a separate job that adds some data to that Postgres
instance. We'll use AwaitPostgres in the `initContainers` block of the data-adding job to make sure that Postgres has fully initialized before we try to add data to it.

First, please make sure you have [Docker](https://www.docker.com/), [Minikube](https://github.com/kubernetes/minikube/releases), [kubectl](https://kubernetes.io/docs/tasks/kubectl/install/), and [watch](http://brewformulas.org/Watch) installed.

Open a terminal and run `watch kubectl get pods`. This way you can monitor pods' status as they go up and down.
Next, open a separate terminal, navigate to your AwaitPostgres directory, and run
```
kubectl create -f kubernetes/toy_job_with_await.yml
```
In your `watch` terminal, you should see the pod come up and wait with status `Init:0/1`. It's using AwaitPostgres to
wait until Postgres comes up before it runs a job to inject some toy data. Now run
```
kubectl create -f kubernetes/postgres/
```
This will start a Postgres deployment. In the `watch` terminal you can see that once Postgres finishes coming up the job will
finally execute and the pod will disappear. If you want, you can attach to the Postgres container and see that the data
was inserted (how to do this is left as an exercise for the reader).
If you want to see what the behavior looks like without AwaitPostgres, do the above steps but use `toy_job_without_await.yml`
instead. The job will eventually run, but it will take a lot longer and probably enter a few crash loop backoffs.

#### Spin up with tmux

Prerequisites (macOS):

```bash
brew install watch tmux
```

Watch spin up:

```bash
./tmux.sh
kubectl create -f kubernetes/toy_job_with_await.yml
kubectl create -f kubernetes/postgres/
```

Tear down tmux session:

```bash
tmux kill-session -t awaitpostres
```
