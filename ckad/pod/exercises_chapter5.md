## Observability - exercises

1. Define a new Pod named web-server with the image nginx in a YAML manifest. Expose the container port 80. Do not create the Pod yet.

```bash
k run web-server --image=nginx --restart=Never --dry-run=client -o=yaml --port=80 > pod.yaml
```

2. For the container, declare a startup probe of type httpGet. Verify that the root context endpoint can be called. Use the default configuration for the probe.

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: web-server
  name: web-server
spec:
  containers:
  - image: nginx
    name: web-server
    ports:
    - containerPort: 80
    resources: {}
    startupProbe:
      httpGet:
        path: /
        port: 80
  dnsPolicy: ClusterFirst
  restartPolicy: Never
status: {}
```

3. For the container, declare a readiness probe of type httpGet. Verify that the root context endpoint can be called. Wait five seconds before checking for the first time.

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: web-server
  name: web-server
spec:
  containers:
  - image: nginx
    name: web-server
    ports:
    - containerPort: 80
    resources: {}
    startupProbe:
      httpGet:
        path: /
        port: 80
    readinessProbe:
      httpGet:
        path: /
        port: 80
      initialDelaySeconds: 5
  dnsPolicy: ClusterFirst
  restartPolicy: Never
status: {}
```

4. For the container, declare a liveness probe of type httpGet. Verify that the root context endpoint can be called. Wait 10 seconds before checking for the first time. The probe should run the check every 30 seconds.

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: web-server
  name: web-server
spec:
  containers:
  - image: nginx
    name: web-server
    ports:
    - containerPort: 80
    resources: {}
    startupProbe:
      httpGet:
        path: /
        port: 80
    readinessProbe:
      httpGet:
        path: /
        port: 80
      initialDelaySeconds: 5
    livenessProbe:
      httpGet:
        path: /
        port: 80
      initialDelaySeconds: 10
      periodSeconds: 30
  dnsPolicy: ClusterFirst
  restartPolicy: Never
status: {}
```

5. Create the Pod and follow the lifecycle phases of the Pod during the process.

6. Inspect the runtime details of the probes of the Pod.

```bash
kubectl describe pod web-server
```

7. Retrieve the metrics of the Pod (e.g., CPU and memory) from the metrics server.

```bash
kubectl top pod web-server
```

8. Create a Pod named custom-cmd with the image busybox. The container should run the command top-analyzer with the command-line flag --all.

```bash
kubectl run custom-cmd --image=busybox --restart=Never -- /bin/sh -c "top-analyzer --all"
kubectl get pod custom-cmd
```

9. Inspect the status. How would you further troubleshoot the Pod to identify the root cause of the failure? 

```bash
kubectl logs custom-cmd
```