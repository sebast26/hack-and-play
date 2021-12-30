## Pod Design - exercises

1. Create three Pods that use the image nginx. The names of the Pods should be pod-1, pod-2, and pod-3. Assign the label tier=frontend to pod-1 and the label tier=backend to pod-2 and pod-3. All pods should also assign the label team=artemidis.

```bash
k run pod-1 --restart=Never --image=nginx -l=tier=frontend,team=artemidis
k run pod-2 --restart=Never --image=nginx -l=tier=backend,team=artemidis
k run pod-3 --restart=Never --image=nginx -l=tier=backend,team=artemidis
```

2. Assign the annotation with the key deployer to pod-1 and pod-3. Use your own name as the value.

```bash
k annotate pod pod-1 deployer='SG'
k annotate pod pod-3 deployer='SG'
```

3. From the command line, use label selection to find all Pods with the team artemidis or aircontrol and that are considered a backend service.

```bash
k get pods --show-labels -l='team in (artemidis, aircontrrol), tier=backend'
```

4. Create a new Deployment named server-deployment. The Deployment should control two replicas using the image grand-server:1.4.6.

```bash
k create deploy server-deployment --image=grand-server:1.4.6 --replicas=2
```

5. Inspect the Deployment and find out the root cause for its failure.

```bash
k get pods
```

6. Fix the issue by assigning the image nginx instead. Inspect the rollout history. How many revisions would you expect to see?

```bash
k set image deploy server-deployment grand-server=nginx --record
k rollout history deploy server-deployment
```

7. Create a new CronJob named google-ping. When executed, the Job should run a curl command for google.com. Pick an appropriate image. The excution should occur every two minutes.

```bash
k create cj google-ping --image=busybox --schedule="*/2 * * * *" -- curl google.com
```

8. Tail the logs of the CronJob at runtime. Check the command-line options of the relevant command or consult the Kubernetes documentation.

??

9. Reconfigure the CronJob to retain a history of seven executions.
10. Reconfigure the CronJob to disallow a new execution if the current execution is still running. Consult the Kubernetes documentation for more information.

```bash
k get cj google-ping -o=yaml > cj.yaml
```

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  creationTimestamp: "2021-12-30T21:30:39Z"
  name: google-ping
  namespace: default
  resourceVersion: "918654"
  selfLink: /apis/batch/v1beta1/namespaces/default/cronjobs/google-ping
  uid: 64a4b059-5e36-4ee6-a3ed-c65e5ab60742
spec:
  concurrencyPolicy: Forbid
  failedJobsHistoryLimit: 1
  jobTemplate:
    metadata:
      creationTimestamp: null
      name: google-ping
    spec:
      template:
        metadata:
          creationTimestamp: null
        spec:
          containers:
          - command:
            - curl
            - google.com
            image: busybox
            imagePullPolicy: Always
            name: google-ping
            resources: {}
            terminationMessagePath: /dev/termination-log
            terminationMessagePolicy: File
          dnsPolicy: ClusterFirst
          restartPolicy: OnFailure
          schedulerName: default-scheduler
          securityContext: {}
          terminationGracePeriodSeconds: 30
  schedule: '*/2 * * * *'
  successfulJobsHistoryLimit: 7
  suspend: false
```