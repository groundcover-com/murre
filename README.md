<div align="center">
  <h1>Murre</h1>
  <p align="center">
    <img src="images/logo.png" width="25%" alt="murre" title="murre" />
   </p>
    <h2>On demand Kubernetes metrics at scale</h2>
   <a href="https://groundcover.com/blog/murre"><strong>Read More »</strong></a>
  <p>


  [![slack](https://cdn.bfldr.com/5H442O3W/at/pl546j-7le8zk-e2zxeg/Slack_Mark_Monochrome_White.png?auto=webp&format=png&width=50&height=50)](https://www.groundcover.com/join-slack)
</div>


<p align="center">
<img src="images/demo.gif" width="100%" alt="murre" title="murre" />
</p>

## What is Murre?
Murre is an **on-demand, scalable source of container resource metrics for K8s**.

Murre fetches CPU & memory resource metrics directly from the kubelet on each K8s Node.
Murre also enriches the resources with the relevant K8s requests and limits from each PodSpec.

## Why Murre?
Murre is a simple, stateless and minimalistic approach to K8s resource monitoring that works at any scale.
Murre is free of any third-party dependencies, requiring nothing to be installed on the cluster.

## Installing Murre
```go
go install github.com/groundcover-com/murre@latest
```

## Using Murre
- Detect pods and containers with high CPU or memory utilization
```bash
murre --sortby-cpu-util
```
- Find out how much of CPU and memory does a specific pod consumes
```bash
murre --pod kong-51xst
```
- Focus on the resource consumption metrics in a specific namespace
```bash
murre --namespace production
```

