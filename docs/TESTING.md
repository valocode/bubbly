# Testing

This file describes the testing process for bubbly.

## Unit tests

### Golang

Nothing fancy, if you want to rerun all the unit tests just do `go test ./...`

### NodeJS

TODO

## Integration tests

This is where things get a little more interesting.
There are currently two types of integration tests:

1. `incluster` tests are those which run inside a k8s cluster because they have direct dependencies on things like databases or services that are not exposed outside the cluster. Thus, these tests need to be run from within the cluster
2. `extcluster` integration tests run against the externally exposed services, such as the bubbly server API

### Tooling

There are two main technologies to help with running integration tests (besides the actual testing frameworks and general k8s tooling):

1. kind: <https://kind.sigs.k8s.io/>
2. skaffold: <https://skaffold.dev/>

`kind` gives us a local k8s cluster inside a single docker container.
What's great about this is that the workflow we use locally can also be integrated into CI.

`skaffold` provides the connection between our code and a deployment in our k8s cluster (i.e. `kind`).
In a nutshell, it automates the building and tagging of docker images, and then performs a `kubectl apply` of our k8s manifests to deploy those built images.
There are two main commands you need to know about when using `skaffold`:

1. `skaffold dev` to run in development mode (with the ability to reload code changes). This is likely what you want locally for development.
2. `skaffold run` to simply run a deployment and finish. This is what we do in CI, and probably what you want to do before submitting your PR :)

Most of the commands needed to run tests are integrated into the root `Makefile`.

### Workflow

To run this workflow, make sure you have `kind`, `skaffold` and `kubectl` installed.
You might want some more tools for working with k8s but that's optional for you.

#### Development

Whilst you are developing, you probably want to do the following

```bash
# run this once, to boot strap your kind (k8s in docker) environment
make kind-bootstrap

# now you should have a kubernetes cluster created... you can check this with kind
kind get clusters
# output should be "bubbly" which is the name of our cluster

# now start your dev mode, this is the command you run over and over again...
# this will invoke skaffold which will build docker images and apply our k8s
# manifests into our kind cluster
#
# skaffold dev is run with --trigger manual. This means YOU need to press enter
# if you want skaffold to rebuild and redeploy... otherwise it sits in the
# background and goes crazy every time you save :)
#
# NOTE: deploying the nginx ingress will likely look like it's failing on the
# first deploy. Fear not, eventually the k8s resources will align themselves
# and it will work.
make dev

# check that it is running
curl http://localhost/v1/status
```

Now you have bubbly running in your local `kind` cluster :)

```bash
# get all pods running in default namespace
kubectl get pods -n default

# get all pods running
kubectl get pods --all-namespaces
```

Notice that we have pods like postgresql running.

##### Running integration tests locally

What's cool is that you can now run your integration tests locally...
So start another terminal up!

```bash
# set the host and port, and then run our integration tests
BUBBLY_HOST=localhost BUBBLY_PORT=80 go test ./integration -tags=integration
```

Note: this will not run the `incluster` tests. To do that, we need to run the tests from within the cluster (keep reading).

#### Run all integration tests

To run all the integration tests from within the cluster, follow the bootstrapping steps above and then do:

```bash
# this will run "skaffold run" with the integration profile, which will deploy
# bubbly to our k8s cluster and then also deploy a kubernetes Job, that will
# run the integration tests.
#
# it then gets the logs using kubectl so that you can see what happened.
make integration

# if you also want the logs from the bubbly server, do
kubectl logs bubbly-0
```
