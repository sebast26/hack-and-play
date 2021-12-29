## Core concepts - exercises

1. Create a new Pod named nginx running the image nginx:1.17.10. Expose the container port 80. The Pod should live in the namespace named ckad.

```bash
kubectl create ns ckad
kubectl run nginx --image=nginx:1.17.10 --restart=Never --port=80 -n ckad
```

**You have to create a namespece prior to creating a pod**

2. Get the details of the Pod including its IP address

```bash
kubectl get pod nginx -o=wide
```

3. Create a temporary Pod that uses the busybox image to execute a wget command inside of the container. The wget command should access the endpoint exposed by the nginx container. You should see the HTML response body rendered in the terminal.

```bash
kubectl run busybox --image=busybox --restart=Never --rm -it -n ckad -- wget -O- 10.1.0.66:80
```

**You can use `-it` and `--rm` to run temporary pod**

**`wget` has an option to output to stdout: `-O-`**

4. Get the logs of the nginx container.

```bash
kubectl logs nginx
```

5. Add the environment variables DB_URL=postgresql://mydb:5432 and DB_USERNAME=admin to the container of the nginx Pod.

**You cannot change the env variables on running pod. You will get following error: Forbidden: pod updates may not change fields other than `spec.containers[*].image`, `spec.initContainers[*].image`, `spec.activeDeadlineSeconds` or `spec.tolerations`**

```bash
kubectl run nginx --image=nginx:1.17.10 --restart=Never --port=80 --env="DB_URL=postgresql://mydb:5432" --env="DB_USERNAME=admin"
```

**Each `env` var needs to be passed separately with `--env` flag**

6. Open a shell for the nginx container and inspect the contents of the current directory ls -l

```bash
kubectl exec -ti nginx -- /bin/sh
```

7. Create a YAML manifest for a Pod named loop that runs the busybox image in a container. The container should run the following command: for i in {1..10}; do echo "Welcome $i times"; done. Create the Pod from the YAML manifest. What’s the status of the Pod?

**For some reason you have to use `for i in 1 2 3 4 5 6 7 8 9` instread of the command given**

```bash
k run loop --image=busybox --restart=Never -o yaml --dry-run=client > loop.yaml -- /bin/sh -c 'for i in 1 2 3 4 5 6 7 8 9; do echo "Welcome $i times";done'
```

8. Edit the Pod named loop. Change the command to run in an endless loop. Each iteration should echo the current date.

```bash
k run loop --image=busybox --restart=Never -o yaml --dry-run=client > loop.yaml -- /bin/sh -c 'while true; do echo `date`; done'
```

9. Inspect the events and the status of the Pod loop.

```bash
kubectl describe pod loop | grep -C 10 Events:
```

10. Delete the namespace ckad and its Pods.

```bash
k delete ns ckad
```