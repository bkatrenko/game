## Motivation
To prepare an application for deployment, we have to make a local setup where we could play around, experiment or test new changes if there are any.
Currently, k8s folder contains only "local" folder, which could be used for debugging purposes and to play around with the application.
Configuration well-tested with MacOS docker-desktop only local cluster.

## How to
Everything is there, to deploy an application locally, you just should build a docker-containers for an application (currently it is expected container name "hockey-server" and a "0.0.5" tag). 
After that, it is possible to just deploy the application, everything including the namespace defined in ./local/*.yaml files.
I would suggest deploying the ingress firstly (it is a slightly changed copy of the default setup from  https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.2.0/deploy/static/provider/cloud/deploy.yaml). We have to change it while we have to support UDP as well as HTTP.

To deploy the NGINX ingress controller, use 
```
kubectl apply -f ingress-nginx.yaml
```
The deployment will take some time, so in this case, we could execute 
```
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```
command to wait until ingress pods will be ready. We could also manually check the state components in nginx-ingress namespace with the command: 
```
kubectl get all -n ingress-nginx
```
The expected output will look like this:
```
NAME                                            READY   STATUS      RESTARTS   AGE
// nginx-admission is a the service for the validating webhook that ingress-nginx includes.
// After it is done we'll have an admission service that will be responsible for this part
pod/ingress-nginx-admission-create-vn2nt        0/1     Completed   0          4h7m
pod/ingress-nginx-admission-patch-cbltn         0/1     Completed   0          4h7m
pod/ingress-nginx-controller-6dbbc469fb-mkzqt   1/1     Running     0          4h7m

NAME                                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
// Here is TCP & UDP services that are responsible for routing data inside the cluster
// Currently to keep things simple we use only one server (monolith) to serve HTTP & UDP endpoints 
service/ingress-nginx-controller             NodePort    10.111.191.15    <none>        80:30080/TCP,443:31156/TCP   4h7m
service/ingress-nginx-controller-admission   ClusterIP   10.96.97.186     <none>        443/TCP                      4h7m
service/ingress-nginx-controller-udp         NodePort    10.107.167.236   <none>        8087:30087/UDP               4h7m

NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ingress-nginx-controller   1/1     1            1           4h8m

NAME                                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/ingress-nginx-controller-6dbbc469fb   1         1         1       4h8m

NAME                                       COMPLETIONS   DURATION   AGE
job.batch/ingress-nginx-admission-create   1/1           15s        4h8m
job.batch/ingress-nginx-admission-patch    1/1           16s        4h8m
```
After we ensure that everything works well, we could run the application. We could check logs with 
```
kubectl logs hockey-server-85c77755f7-xz9xw -n hockey
```
Command. There should be an indication that server started successfully:
```
{
  "level": "info",
  "http_address": ":8080",
  "udp_address": ":8081",
  "time": 1656342038,
  "message": "start hockey server"
}
```
After we ensure that server works as expected, it's time to start the client with the following command:

```
UDP_SERVER_HOST_PORT=hockey.dev:30087 HTTP_SERVER_HOST_PORT=http://hockey.dev:30080 PLAYER_ID=bkatrenko GAME_ID=1 PLAYER_NUMBER=0 ./game
```
Where we could set any GAME_ID and PLAYER_ID. The first player that will join will create a new game.
