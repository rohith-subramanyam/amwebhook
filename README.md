# amwebhook
Sample webhook for Prometheus Alertmanager

## What is it?
This demonstrates the workflow to receive alerts from Prometheus
Alertmanager to a custom generic receiver. This is achieved using Prometheus
Alertmanager's [webhook_config](https://prometheus.io/docs/alerting/configuration/#webhook_config).
This webhook receives HTTP POST requests from Prometheus Alertmanager in the
following JSON format to its webhook endpoint:
```json
{
  "version": "4",
  "groupKey": "<string>",    // key identifying the group of alerts (e.g. to deduplicate)
  "status": "<resolved|firing>",
  "receiver": "<string>",
  "groupLabels": "<object>",
  "commonLabels": "<object>",
  "commonAnnotations": "<object>",
  "externalURL": "<string>",  // backlink to the Alertmanager.
  "alerts": [
    {
      "status": "<resolved|firing>",
      "labels": "<object>",
      "annotations": "<object>",
      "startsAt": "<rfc3339>",
      "endsAt": "<rfc3339>",
      "generatorURL": "<string>" // identifies the entity that caused the alert
    },
    ...
  ]
}
```
It simply logs the alerts.

## Build
### Prerequisites
Since this uses a [Dockerfile](Dockerfile) with [multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/), it needs Docker 17.05 or higher.

### Steps
```shell
$ git clone https://github.com/rohith-subramanyam/amwebhook.git
$ cd amwebhook
$ docker build -t amwebhook .  # Build the docker image. The Dockerfile is a multi-stage build.
$ docker tag amwebhook rohithvsm/amwebhook:1.0.0  # Tag the image.
$ docker login  # Login to docker with your docker ID.
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: rohithvsm
Password:
Login Succeeded
$ docker push rohithvsm/amwebhook:1.0.0  # Push the image to docker hub.
```

## Deploy
Create a deployment in MSP/K8s cluster using [this deployment yaml](k8s/deploy.yaml) which uses the image we just pushed above.
```shell
$ kubectl apply -f deploy.yaml
```
The deployment points to the image we just pushed:
```shell
$ git grep "image:" deploy.yaml
          image: rohithvsm/rohithvsm:alertmgrwebhook
```
Create a service to expose the app so that alertmanager can send the alerts to its webhook endpoint using [this service yaml](service.yaml):
```shell
$ kubectl apply -f service.yaml
```

## Configure Alertmanager
Configure the alertmanager to send the alerts to our app using the [webhook_configs](https://prometheus.io/docs/alerting/configuration/#webhook_config). You can check an example in [alertmanager.yaml](alertmanager/alertmanager.yaml).
```shell
$ git grep -B1 -A1 webhook_configs alertmanager.yaml
- name: ntnx-system-webhook
  webhook_configs:
  - url: "http://am-webhook.ntnx-system/webhook"
```
For a MSP cluster, base64 encode the alertmanager.yaml and create a secret and apply that. Check [alertmanager-secret.yaml](alertmanager/alertmanager-secret.yaml)
```shell
$ kubectl apply -f alertmanager-secret.yaml
```

## Verify
Verify that our app is receiving the alerts from Prometheus Alertmanager
```shell
$ kubectl logs am-webhook-7f4b7ccf6c-t47nm -n ntnx-system
2019/11/07 17:26:58 listening on: :80
2019/11/12 23:22:23 Alerts: GroupLabels=map[], CommonLabels=map[prometheus:ntnx-system/k8s]
2019/11/12 23:22:23 Alert: status=firing,Labels=map[alertname:KubeCPUOvercommit prometheus:ntnx-system/k8s severity:warning],Annotations=map[message:Cluster has overcommitted CPU resource requests for Pods and cannot tolerate node failure. runbook_url:https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubecpuovercommit]
2019/11/12 23:22:23 Alert: status=firing,Labels=map[alertname:KubeControllerManagerDown prometheus:ntnx-system/k8s severity:critical],Annotations=map[message:KubeControllerManager has disappeared from Prometheus target discovery. runbook_url:https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubecontrollermanagerdown]
2019/11/12 23:22:23 Alert: status=firing,Labels=map[alertname:KubeMemOvercommit prometheus:ntnx-system/k8s severity:warning],Annotations=map[message:Cluster has overcommitted memory resource requests for Pods and cannot tolerate node failure. runbook_url:https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubememovercommit]
2019/11/12 23:22:23 Alert: status=firing,Labels=map[alertname:KubeSchedulerDown prometheus:ntnx-system/k8s severity:critical],Annotations=map[message:KubeScheduler has disappeared from Prometheus target discovery. runbook_url:https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubeschedulerdown]
2019/11/12 23:22:23 Alert: status=firing,Labels=map[alertname:CPUThrottlingHigh namespace:ntnx-system pod_name:fluent-bit-rgj4q prometheus:ntnx-system/k8s severity:warning],Annotations=map[runbook_url:https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-cputhrottlinghigh message:Throttling 93% of CPU in namespace ntnx-system for .]
2019/11/12 23:22:23 Alert: status=firing,Labels=map[alertname:CPUThrottlingHigh container_name:fluent-bit namespace:ntnx-system pod_name:fluent-bit-rgj4q prometheus:ntnx-system/k8s severity:warning],Annotations=map[message:Throttling 98% of CPU in namespace ntnx-system for fluent-bit. runbook_url:https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-cputhrottlinghigh]
2019/11/12 23:33:29 Alerts: GroupLabels=map[job:kubelet], CommonLabels=map[alertname:KubeClientErrors instance:10.45.22.66:10250 job:kubelet prometheus:ntnx-system/k8s severity:warning]
2019/11/12 23:33:29 Alert: status=firing,Labels=map[severity:warning alertname:KubeClientErrors instance:10.45.22.66:10250 job:kubelet prometheus:ntnx-system/k8s],Annotations=map[message:Kubernetes API server client 'kubelet/10.45.22.66:10250' is experiencing 1% errors.' runbook_url:https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubeclienterrors]
```
