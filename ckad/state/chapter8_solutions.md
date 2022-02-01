## State persistence

1. Create a Pod YAML file with two containers that use the image alpine:3.12.0. Provide a command for both containers that keep them running forever.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: alpines
spec:
  containers:
  - image: alpine:3.12.0
    name: first
    command: ["/bin/sh"]
    args:
    - -c
    - "while true; do sleep 60; done;"
  - image: alpine:3.12.0
    name: second
    command: ["/bin/sh"]
    args:
    - -c
    - "while true; do sleep 60; done;"
```

2. Define a Volume of type emptyDir for the Pod. Container 1 should mount the Volume to path /etc/a, and container 2 should mount the Volume to path /etc/b.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: alpines
spec:
  containers:
  - image: alpine:3.12.0
    name: first
    command: ["/bin/sh"]
    args:
    - -c
    - "while true; do sleep 60; done;"
    volumeMounts:
    - mountPath: /etc/a
      name: vol
  - image: alpine:3.12.0
    name: second
    command: ["/bin/sh"]
    args:
    - -c
    - "while true; do sleep 60; done;"
    volumeMounts:
    - mountPath: /etc/b
      name: vol
  volumes:
  - name: vol
    emptyDir: {}
```

3. Open an interactive shell for container 1 and create the directory data in the mount path. Navigate to the directory and create the file hello.txt with the contents “Hello World.” Exit out of the container.

```bash
k exec alpines -c first -it -- /bin/sh
mkdir /etc/a/data
echo "Hello World" > /etc/a/data/hello.txt
```

4. Open an interactive shell for container 2 and navigate to the directory /etc/b/data. Inspect the contents of file hello.txt. Exit out of the container.

```bash
k exec alpines -c second -it -- /bin/sh
cat /etc/b/data/hello.txt
```

5. Create a PersistentVolume named logs-pv that maps to the hostPath /var/logs. The access mode should be ReadWriteOnce and ReadOnlyMany. Provision a storage capacity of 5Gi. Ensure that the status of the PersistentVolume shows Available.

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: logs-pv
spec:
  capacity:
    storage: 5Gi
  accessModes:
  - ReadWriteOnce
  - ReadWriteMany
  hostPath:
    path: /var/logs
```

6. Create a PersistentVolumeClaim named logs-pvc. The access it uses is ReadWriteOnce. Request a capacity of 2Gi. Ensure that the status of the PersistentVolume shows Bound.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: logs-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests: 
      storage: 2Gi
```

7. Mount the PersistentVolumeClaim in a Pod running the image nginx at the mount path /var/log/nginx.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - image: nginx
    name: nginx
    volumeMounts:
    - mountPath: /var/log/nginx
      name: vol
  volumes:
  - name: vol
    persistentVolumeClaim:
      claimName: logs-pvc
```

8. Open an interactive shell to the container and create a new file named my-nginx.log in /var/log/nginx. Exit out of the Pod.

```bash
k exec nginx -it -- /bin/sh
echo "logs..." > /var/log/nginx/my-nginx.log
```

9. Delete the Pod and re-create it with the same YAML manifest. Open an interactive shell to the Pod, navigate to the directory /var/log/nginx, and find the file you created before.

```bash
k delete pod nginx --grace-period=0 --force
```