# Function to create `TaskRun` objects and POST it on the kubernetes API server

The function creates a `TaskRun` object to execute the selected `Task` object  

## Function Deploy

Deploy the function with buildtemplate and env variables ```tm deploy service taskrun -f . --build-template https://raw.githubusercontent.com/triggermesh/knative-lambda-runtime/master/go-1.x/buildtemplate.yaml --env TASK_NAME=yourtaskname --env NAMESPACE=yourNamespace --wait ```
