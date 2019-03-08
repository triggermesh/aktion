## Function to create `TaskRun` objects

This function creates a `TaskRun` object to execute the selected `Task` object  

### Local usage

```
dep ensure -v
go build -o taskrun .
./taskrun
```

### Remote usage

Use `ko`. If you are new to `ko` check out the [demo](https://github.com/sebgoa/kodemo)

Set a registry and make sure you can push to it:

```
export KO_DOCKER_REPO=gcr.io/triggermesh
```

Then `apply` like this:

```
ko apply -f config/
```

To deploy in a different namespace:

```
ko -n nondefault apply -f config/
```
