# prometheus-collector

A service running on k8s cluster and does tasks based on the following requirements:

```
A service written in python or golang that queries 2 urls (https://httpstat.us/503 & https://httpstat.us/200)
The service will check the external urls (https://httpstat.us/503 & https://httpstat.us/200 ) are up (based on http status code 200) and response time in milliseconds
The service will run a simple http service that produces  metrics (on /metrics) and output a prometheus format when hitting the service /metrics url
Expected response format:
sample_external_url_up{url="https://httpstat.us/503 "}  = 0
sample_external_url_response_ms{url="https://httpstat.us/503 "}  = [value]
sample_external_url_up{url="https://httpstat.us/200 "}  = 1
sample_external_url_response_ms{url="https://httpstat.us/200 "}  = [value]
Looking for:
Code in python or golang (with tests)
Dockerfile to build image
Kubernetes Deployment Specification to deploy Image to Kubernetes Cluster
Artifacts should be uploaded to github.com
```

=============================================================
# Deployemnt Instruction
Use command to build docker image inside of app folder:
```
docker build -t prometheus_collector:yanyma . -f Dockerfile
```

If your build the docker image outside your k8s cluster, use the following command to load your image on your k8s:
```
docker save prometheus_collector:yanyma > prometheus_collector.tar
scp prometheus_collector.tar {account}@{your k8s cluster}:{path}
```

on you k8s cluster, load the docker image:
```
docker load < {path}/prometheus_collector.tar
```

Use command to bring up k8s service and k8s deployment:
```
kubectl apply -f prometheusCollector_service.yaml 
kubectl apply -f prometheusCollector_deployment.yaml
```
===============================================================
# Sample Output

```
YANYMA-M-T2S2:prometheus yanyma$ curl http://10.193.219.62:32112/metrics
# HELP sample_external_url_response_ms_200 Response time of https://httpstat.us/200 in milliseconds.
# TYPE sample_external_url_response_ms_200 gauge
sample_external_url_response_ms_200{url="https://httpstat.us/200"} 7
# HELP sample_external_url_response_ms_503 Response time of https://httpstat.us/503 in milliseconds.
# TYPE sample_external_url_response_ms_503 gauge
sample_external_url_response_ms_503{url="https://httpstat.us/503"} 5
# HELP sample_external_url_up_200 Up/Down status of https://httpstat.us/200.
# TYPE sample_external_url_up_200 gauge
sample_external_url_up_200{url="https://httpstat.us/200"} 1
# HELP sample_external_url_up_503 Up/Down status of https://httpstat.us/503.
# TYPE sample_external_url_up_503 gauge
sample_external_url_up_503{url="https://httpstat.us/503"} 0
```

===============================================================
# Future Enhancement

1. The port number of http service to produce metrics currently is hardcoded as 2112.
This port number could be configured as environment variable in k8s deployment file.

2. Again, the port number 2112 in k8s deployment/service yaml file could be set in config file, like Ansible Playbook config file

3. Unit Test now just cover the CollectFunc where we have the main logic. It could be extended to have more code coverage.

4. Dock image is built from golang:1.13 and its size is quite big. It could be built from golang:alpine, but git installation is required since it is not included in golang:alpine

5. Use glide instead if manually pull depent package

6. Scalability consideration.
