# Function to create `TaskRun` objects and POST it on the kubernetes API server

The function creates a `TaskRun` object to execute the selected `Task` object  


### Local usage

```
dep ensure -v
go build -o taskrun .
./taskrun
```

### Local Docker Usage

```
docker build -t taskrun:latest . 
docker run -ti -e TASK_NAME="your_task_name" \
               -e NAMESPACE="your_namespace" \
               taskrun:latest
```