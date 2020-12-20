### Bee-crontab
[![GoDoc](https://godoc.org/github.com/sinksmell/bee-crontab?status.svg)](https://godoc.org/github.com/sinksmell/bee-crontab)

> 基于go语言实现的一个分布式的任务调度系统,类似crontab

* 可视化Web后台管理,配置定时任务方便
* 分布式架构,集群化调度,可用性高
* 秒级调度,准确性高
* 追踪任务执行状态,采集任务输出
* 可视化日志查看



#### Bee-crontab 架构图

![](https://i.loli.net/2019/03/18/5c8f6c4881cd1.png)

#### worker节点核心 调度器架构图
![](https://i.loli.net/2019/03/18/5c8f6c61c2145.png)

#### 应用界面

* 主界面
![](https://i.loli.net/2019/03/18/5c8f6cd83b5d9.png)

* 添加任务
![](https://i.loli.net/2019/03/18/5c8f6cdaabc08.png)

* 查看当前worker节点信息
![](https://i.loli.net/2019/03/18/5c8f6cdd15016.png)

* 查看任务执行日志
![](https://i.loli.net/2019/03/18/5c8f6ce01f4ca.png)

#### k8s部署

![](https://i.loli.net/2019/06/06/5cf86f04438e664010.png)

#### 测试集群

* 建立上图对应的k8s资源,检查pod状态

```
master@ubuntu:~/k8s/bee_crontab$ kubectl get po
NAME                              READY   STATUS    RESTARTS   AGE
bc-master-76df4bb9cf-l2wz2        1/1     Running   2          12h
bc-master-76df4bb9cf-vl2h8        1/1     Running   0          18s
bc-worker-58bd56f44c-5gzjl        1/1     Running   2          12h
bc-worker-58bd56f44c-cwrbc        1/1     Running   2          12h
bc-worker-58bd56f44c-d98vq        1/1     Running   2          12h
etcd-operator-84db9bc774-lzvrs    1/1     Running   10         3d
example-etcd-cluster-b64p5wxst8   1/1     Running   0          3m36s
mongo-69d6d44cb4-4lst7            1/1     Running   4          3d18h

```

*  查看Service 的ip

```
master@ubuntu:~/k8s/bee_crontab$ kubectl get svc
NAME                                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
bc-master                             ClusterIP   10.152.183.229   <none>        80/TCP              13h
example-etcd-cluster                  ClusterIP   None             <none>        2379/TCP,2380/TCP   7m2s
example-etcd-cluster-client           ClusterIP   10.152.183.2     <none>        2379/TCP            7m2s
example-etcd-cluster-client-service   NodePort    10.152.183.54    <none>        2379:32379/TCP      2d23h
kubernetes                            ClusterIP   10.152.183.1     <none>        443/TCP             6d23h
mongo                                 ClusterIP   10.152.183.161   <none>        27017/TCP           3d18h
nginx-app                             ClusterIP   10.152.183.159   <none>        80/TCP              6d22h

```

* 根据Service的信息修改master和worker的配置文件中的etcd和mongodb的地址

* 使用浏览器访问master提供的管理界面

#### 测试结果
> 部分图标资源未找到

> * 主界面
![](https://i.loli.net/2019/06/06/5cf879130a4f615338.png)

> * 查看健康节点
![](https://i.loli.net/2019/06/06/5cf879111212a91053.png)

> * 添加任务
![](https://i.loli.net/2019/06/06/5cf879107e68e78035.png)

> * 查看任务执行日志
![](https://i.loli.net/2019/06/06/5cf87913a87e368264.png)
![](https://i.loli.net/2019/06/06/5cf87911bd7a574250.png)
![](https://i.loli.net/2019/06/06/5cf879147480a54313.png)

可以看到集群可以正常工作,任务不会被并发调度。