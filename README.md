<div align="center">
  <h1>Murre</h1>
    <h2>K8s on demand metrics, free of third-party dependencies</h2>
   <a href="link_to_blog"><strong>Read more »</strong></a>
  <p>

  [![slack](https://cdn.bfldr.com/5H442O3W/at/pl546j-7le8zk-e2zxeg/Slack_Mark_Monochrome_White.png?auto=webp&format=png&width=50&height=50)](https://www.groundcover.com/join-slack)
</div>

<p align="center">
<img src="images/demo.gif" width="100%" alt="murre" title="murre" />
</p>

## What is murre?
murre is a on-demand, scaleable source of container resource metrics for K8s.

murre fetchs resource metrics directly from the kubelet on each K8s Node. murre also enrich the resources with the relevant limits from each PodSpec.

## What is murre useful for?
- Detect pods and containers with high cpu/memory utilization
- Find out how much of cpu and memory does specific pod consumes

## Why murre?
somtimes, you just want to go simple. stateless. minimil as possible like murre.
