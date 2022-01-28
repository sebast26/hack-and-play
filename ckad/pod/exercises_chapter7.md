# Services and Networking - exercises

1. Create a new Pod named frontend that uses the image nginx. Assign the labels tier=frontend and app=nginx. Expose the container port 80.

```bash
k run frontend --restart=Never --image=nginx -l=tier=frontend,app=nginx --port 80
```

2. Create a new Pod named backend that uses the image nginx. Assign the labels tier=backend and app=nginx. Expose the container port 80.

```bash
k run backend --restart=Never --image=nginx -l=tier=backend,app=nginx --port 80
```

3. Create a new Service named nginx-service of type ClusterIP. Assign the port 9000 and the target port 80. The label selector should use the criteria tier=backend and deployment=app.

```bash
k create service clusterip nginx-service --tcp=9000:80 --dry-run=client -o=yaml > svc.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: nginx-service
  name: nginx-service
spec:
  ports:
  - port: 9000
    protocol: TCP
    targetPort: 80
  selector:
    tier: backend
    deployment: app
  type: ClusterIP
status:
  loadBalancer: {}
```

4. Try to access the set of Pods through the Service from within the cluster. Which Pods does the Service select?

```bash
k run busybox -it --rm --image=busybox --restart=Never -- wget -O- 10.101.207.209:9000
```

Can't connect. Connection refused.

5. Fix the Service assignment to properly select the backend Pod and assign the correct target port.

```yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: nginx-service
  name: nginx-service
spec:
  ports:
  - port: 9000
    protocol: TCP
    targetPort: 80
  selector:
    tier: backend
    app: nginx
  type: ClusterIP
status:
  loadBalancer: {}
```

```bash
k run busybox -it --rm --image=busybox --restart=Never -- wget -O- 10.108.166.243:9000
```

Now it can be accessed.

6. Expose the Service to be accessible from outside of the cluster. Make a call to the Service.

```yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: nginx-service
  name: nginx-service
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
    nodePort: 30200
  selector:
    tier: backend
    app: nginx
  type: NodePort
status:
  loadBalancer: {}
```

7. Assume an application stack that defines three different layers: a frontend, a backend, and a database. Each of the layers runs in a Pod. You can find the definition in the YAML file app-stack.yaml

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: frontend
  namespace: app-stack
  labels:
    app: todo
    tier: frontend
spec:
  containers:
  - name: frontend
    image: nginx

---

kind: Pod
apiVersion: v1
metadata:
  name: backend
  namespace: app-stack
  labels:
    app: todo
    tier: backend
spec:
  containers:
  - name: backend
    image: nginx

---

kind: Pod
apiVersion: v1
metadata:
  name: database
  namespace: app-stack
  labels:
    app: todo
    tier: database
spec:
  containers:
  - name: database
    image: mysql
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: example
```

```bash
k create ns app-stack
```

8. Create a network policy in the file app-stack-network-policy.yaml. The network policy should allow incoming traffic from the backend to the database but disallow incoming traffic from the frontend.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: database-policy
spec:
  policyTypes:
  - Ingress
  - Egress
  podSelector:
    matchLabels:
      app: todo
      tier: database
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: todo
          tier: backend
```

9. Reconfigure the network policy to only allow incoming traffic to the database on TCP port 3306 and no other port.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: database-policy
spec:
  policyTypes:
  - Ingress
  - Egress
  podSelector:
    matchLabels:
      app: todo
      tier: database
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: todo
          tier: backend
    ports:
    - port: 3306
```