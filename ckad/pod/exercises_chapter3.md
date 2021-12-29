## Configuration - exercises

1. Create a directory with the name config. Within the directory, create two files. The first file should be named db.txt and contain the key-value pair password=mypwd. The second file is named ext-service.txt and should define the key-value pair api_key=LmLHbYhsgWZwNifiqaRorH8T.

```bash
mkdir config
echo -n "password=mypwd" > config/db.txt
echo -n "api_key=LmLHbYhsgWZwNifiqaRorH8T" > config/ext-service.txt
```

2. Create a Secret named ext-service-secret that uses the directory as data source and inspect the YAML representation of the object.

```bash
kubectl create secret generic ext-service-secret --from-file=config
kubectl get secret ext-service-secret -o=yaml
```

3. Create a Pod named consumer with the image nginx and mount the Secret as a Volume with the mount path /var/app. Open an interactive shell and inspect the values of the Secret.

```bash
kubectl run consumer --restart=Never --image=nginx -o=yaml --dry-run=client > pod.yaml
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: consumer
  name: consumer
spec:
  containers:
  - image: nginx
    name: consumer
    volumeMounts:
    - name: ext-service-secrets-vol
      mountPath: /var/app
  volumes:
  - name: ext-service-secrets-vol
    secret:
      secretName: ext-service-secret
```

```bash
kubectl exec -ti consumer -- /bin/sh
```

4. Use the declarative approach to create a ConfigMap named ext-service-configmap. Feed in the key-value pairs api_endpoint=https://myapp.com/api and username=bot as literals.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ext-service-configmap
data:
  api_endpoint: https://myapp.com/api
  username: bot
```

5. Inject the ConfigMap values into the existing Pod as environment variables. Ensure that the keys conform to typical naming conventions of environment variables.

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: consumer
  name: consumer
spec:
  containers:
  - image: nginx
    name: consumer
    volumeMounts:
    - name: ext-service-secrets-vol
      mountPath: /var/app
    env:
    - name: API_ENDPOINT
      valueFrom:
        configMapKeyRef:
          key: api_endpoint
          name: ext-service-configmap
    - name: USERNAME
      valueFrom:
        configMapKeyRef:
          key: username
          name: ext-service-configmap
  volumes:
  - name: ext-service-secrets-vol
    secret:
      secretName: ext-service-secret
```

6. Open an interactive shell and inspect the values of the ConfigMap.

```bash
kubectl exec -ti consumer -- env
```

7. Define a security context on the container level of a new Pod named security-context-demo that uses the image alpine. The security context adds the Linux capability CAP_SYS_TIME to the container. Explain if the value of this security context can be redefined in a Pod-level security context.

```bash
k run security-context-demo --restart=Never --image=alpine --dry-run=client -o=yaml > pod.yaml
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: security-context-demo
  name: security-context-demo
spec:
  containers:
  - image: alpine
    name: security-context-demo
    securityContext:
      capabilities:
        add: ["CAP_SYS_TIME"]
    command: ["sleep", "3600"]
```

8. Define a ResourceQuota for the namespace project-firebird. The rules should constrain the count of Secret objects within the namespace to 1.

```bash
k create ns project-firebird
k create quota project-firebird-quota --hard=secrets=1 -n project-firebird
```

```bash
k describe quota -n project-firebird
```

10. Create a new Service Account named monitoring and assign it to a new Pod with an image of your choosing. Open an interactive shell and locate the authentication token of the assigned Service Account.

```bash
k create sa monitoring
k run sa-demo --restart=Never --image=nginx --dry-run=client -o=yaml > pod.yaml
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: sa-demo
  name: sa-demo
spec:
  serviceAccountName: monitoring
  containers:
  - image: nginx
    name: sa-demo
```

```bash
k describe pod sa-demo | grep -A 1 Mounts:
k exec -ti sa-demo -- /bin/bash
cat /var/run/secrets/kubernetes.io/serviceaccount/token
```





































