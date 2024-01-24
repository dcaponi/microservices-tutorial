# Kubernetes Local Dev 101
## Config
`kubectl config` is where things like `context` and `namespace` live.

`namespace` is a logical separation of kubernetes concepts. It's basically a virtual cluster in the actual cluster the company pays for. For a solo developer this likely has little use.

`context` pretty much the cluster, but more specifically, a collection of cluster, user, and default namespace. This would likely be local-context and prod-context for a solo dev.

## Objects

### Deployments
- specify how many pods, their resource allocations, ports, and most importantly, images
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: <name>
spec:
  replicas: <how many pods you want>
  selector:
    matchLabels:
      app: <label (I use the name)> # Finds pods by labels to which this deployment applies
  template:
    metadata:
      labels:
        app: <metadata label name> # labels to be applied to pods (for grouping in services)
    spec:
      containers:
      - name: <container name> # identifies containers
        image: <docker-repo/image-name:ver or latest>
        resources:
          limits: # maximum resource the container can use
            memory: "128Mi"
            cpu: "500m"
          requests: # minimum resource the container needs to start and run
              memory: "128Mi"
              cpu: "500m"
        ports:
        - containerPort: <whatever port your app server is mapped to>
```

#### Requests v Limits
##### Requests
- minimum
- may go over if node has free resources to spare
- k8s uses this to figure which node to put the pod on
##### Limits
- cap
- k8s will throttle the pod if it goes over this amount or oomkill it
- prevents a greedy pod from hogging all the node resources
  
### Services
- defines a logical set of pods via `spec:selector:app` and the method how to define them
```yaml
apiVersion: v1
kind: Service
metadata:
  name: <name of service>
spec:
  selector:
    app: <pod label to which service applies> # should match deployment:spec:template:metadata:labels:app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8082
  type: NodePort
```
#### Types
##### ClusterIP
- default
- service (group of pods) only reachable from within the cluster
##### NodePort
- expose the service on each node's ip.
- access from outside via <nodeIP>:<nodePort>
##### LoadBalancer
- expose service using a cloud provider's load balancer
##### ExternalName
- maps the service to the contents of the externalName field by returning a CNAME record with the value

### Ingress
- Defines a high level way to access apps. Instead of an ip:port combo you can use names like my-app.local
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: platform-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: my-app.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: <name of service> # matches service:metadata:name
            port:
              number: 80
```
#### Minikube + Apple Silicon specific ingress notes
- Get an [ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)
- modify `/etc/hosts` to include
```sh
127.0.0.1 <service-name>.local
127.0.0.1 <other-service-name>.local
```
This is so we can `curl http://<service-name>.local` and that means we can reference it like that in other code running locally.
- Enable the minikube ingress addon
  - `minikube addons enable ingress`
  - if `minikube addons enable ingress` hangs
    - `minikube ssh docker pull k8s.gcr.io/ingress-nginx/controller:v1.2.1` should pull the correct nginx [ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)
- `minikube tunnel` [see m1 specific comment](https://github.com/kubernetes/minikube/issues/13510#issuecomment-1193160620)

## Local Startup
### kubectl raw doggin' it 🌭
```sh
minikube start
minikube addons enable ingress
minikube ssh docker pull k8s.gcr.io/ingress-nginx/controller:v1.2.1
```
- For each service, make sure you build an image and push it up to docker hub
`docker build . -t repo_name/image_name:version && docker push repo_name/image_name:version`
- Create a `kubernetes` folder and add `ingress.yaml` based on the [Ingress](#ingress) section
- create a `service-name` folder for each microservice and in each add the `deployment.yaml` and `service.yaml` based on [deployment](#deployments) and [service](#services) sections
```sh
kubectl apply -f ./kubernetes/<service_name> # applies deployments and services
kubectl apply -f ./kubernetes/ingress.yaml # applies yaml
```
- Update `/etc/hosts`
```sh
# /etc/hosts
# add this line for each service
127.0.0.1 <service_name>.local 
```
run `minikube tunnel` in another terminal window
`curl http://<service_name>.local` and it should happen 🎉

### Tilt - automagic 🪄🎩
- get tilt. it makes it so you dont have to build and deploy images manually or run `kubectl apply -f <abc>` to deploy apps
```sh
curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
```
- in the `kubernetes folder` run `tilt up` it will make a Tiltfile
- Modify the `Tiltfile` so it has the following
```python
k8s_yaml(
    [
        "<app_name>/deployment.yaml",
        "<app_name>/service.yaml",
        "<app_name2>/deployment.yaml",
        "<app_name2>/service.yaml",
    ]
)
docker_build("repo/image_name", "../<app_name>/") # or directory where this image's dockerfile lives
docker_build("repo/image_name", "../<app_name>/")
```

*Startup*
```sh
# window 1
minikube start && minikube tunnel # you'll be asked for a password
# window 2
cd <your_frontend_app> && npm run dev
# window 3
cd <kubernetes repo w/ tiltfile> && tilt up
```
