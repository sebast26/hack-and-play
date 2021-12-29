## Multi-Container Pods

1. Create a YAML manifest for a Pod named complex-pod. The main application container named app should use the image nginx and expose the container port 80. Modify the YAML manifest so that the Pod defines an init container named setup that uses the image busybox. The init container runs the command wget -O- google.com.

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: complex-pod
  name: complex-pod
spec:
  initContainers:
  - name: setup
    image: busybox
    command: ["wget", "-O-", "google.com"]
  containers:
  - image: nginx
    name: app
    ports:
    - containerPort: 80
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Never
status: {}
```

2. Create the Pod from the YAML manifest.

```bash
kubectl create -f pod.yaml
```

3. Download the logs of the init container. You should see the output of the wget command.

```bash
k logs complex-pod --container setup
```

4. Open an interactive shell to the main application container and run the ls command. Exit out of the container.

```bash
k exec -ti complex-pod --container app -- /bin/sh
```

5. Force-delete the Pod.

```bash
kubectl delete pod complex-pod --grace-period=0 --force
```

6. Create a YAML manifest for a Pod named data-exchange. The main application container named main-app should use the image busybox. The container runs a command that writes a new file every 30 seconds in an infinite loop in the directory /var/app/data. The filename follows the pattern {counter++}-data.txt. The variable counter is incremented every interval and starts with the value 1.

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: data-exchange
  name: data-exchange
spec:
  containers:
  - image: busybox
    name: main-app
    args:
    - /bin/sh
    - -c
    - mkdir -p /var/app/data; i=1; while true; do touch /var/app/data/$i-data.txt; i=$((i+1)); sleep 30; done
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Never
```

7. Modify the YAML manifest by adding a sidecar container named sidecar. The sidecar container uses the image busybox and runs a command that counts the number of files produced by the main-app container every 60 seconds in an infinite loop. The command writes the number of files to standard output.

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: data-exchange
  name: data-exchange
spec:
  containers:
  - image: busybox
    name: main-app
    args:
    - /bin/sh
    - -c
    - mkdir -p /var/app/data; i=1; while true; do touch /var/app/data/$i-data.txt; i=$((i+1)); sleep 30; done
    resources: {}
  - image: busybox
    name: sidecar
    args:
    - /bin/sh
    - -c
    - while true; do i=`ls /var/app/data | grep data.txt | wc -l`; echo $i; sleep 60; done
  dnsPolicy: ClusterFirst
  restartPolicy: Never
```

8. Define a Volume of type emptyDir. Mount the path /var/app/data for both containers.

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: data-exchange
  name: data-exchange
spec:
  containers:
  - image: busybox
    name: main-app
    args:
    - /bin/sh
    - -c
    - mkdir -p /var/app/data; i=1; while true; do touch /var/app/data/$i-data.txt; i=$((i+1)); sleep 30; done
    resources: {}
    volumeMounts:
    - name: app-data
      mountPath: /var/app/data
  - image: busybox
    name: sidecar
    args:
    - /bin/sh
    - -c
    - while true; do i=`ls /var/app/data | grep data.txt | wc -l`; echo $i; sleep 60; done
    volumeMounts:
    - name: app-data
      mountPath: /var/app/data
  dnsPolicy: ClusterFirst
  restartPolicy: Never
  volumes:
  - name: app-data
    emptyDir: {}
```

9. Create the Pod. Tail the logs of the sidecar container.

```bash
k logs -f data-exchange --container=sideca
```

10. Delete the Pod.

```bash
k delete pod data-exchange --grace-period=0 --force
```