# Kubernetes 资源对象速查表

在 Kubernetes（k8s）里，所有东西几乎都是通过 **资源对象（Resource Object）** 来描述和管理的。资源对象就是对集群中某些状态的声明，k8s
的控制器会不断地去“调谐（reconcile）”实际状态与期望状态。

下面我给你按 **常见分类** 来梳理几种主要的资源对象及其作用：

---

## 1. **工作负载类（Workloads）**

用于定义和运行容器化应用。

* **Pod**

    * 最小的调度单元，包含一个或多个紧密耦合的容器（共享网络、存储）。
    * 一般不直接创建 Pod，而是通过上层控制器来管理。

* **ReplicaSet**

    * 确保指定数量的 Pod 副本始终运行。
    * 很少单独使用，通常由 Deployment 来管理。

* **Deployment**

    * 用于无状态应用，支持 **滚动更新、回滚**。
    * 最常用的工作负载对象。

* **StatefulSet**

    * 用于有状态应用（如数据库、消息队列）。
    * Pod 有固定的网络 ID 和存储卷（如 `pod-0`, `pod-1`）。

* **DaemonSet**

    * 保证每个节点都运行一个副本（例如日志收集、监控 agent）。

* **Job**

    * 一次性任务，执行完成后结束。

* **CronJob**

    * 定时任务，基于 Cron 表达式调度 Job。

---

## 2. **服务发现与负载均衡类（Service & Networking）**

解决 Pod 的网络访问和服务发现问题。

* **Service**

    * 稳定的访问入口，屏蔽 Pod 的 IP 变化。
    * 类型：

        * `ClusterIP`（默认，只能集群内访问）
        * `NodePort`（暴露到节点端口）
        * `LoadBalancer`（结合云厂商的 LB 使用）

* **Ingress**

    * 提供 HTTP/HTTPS 层的七层路由和反向代理。
    * 需要搭配 Ingress Controller 使用（如 Nginx Ingress）。

* **Endpoint / EndpointSlice**

    * 记录 Service 对应的 Pod IP 列表。

* **NetworkPolicy**

    * 定义 Pod 之间或 Pod 与外部的网络访问规则（相当于防火墙策略）。

---

## 3. **存储类（Storage）**

持久化数据存储。

* **Volume**

    * Pod 内的存储挂载。生命周期和 Pod 绑定。

* **PersistentVolume（PV）**

    * 集群级别的持久化存储资源，由管理员预先配置（NFS、云盘等）。

* **PersistentVolumeClaim（PVC）**

    * 用户对存储的申请声明，类似“存储租赁”。

* **StorageClass**

    * 定义不同类型的存储（SSD、HDD、分布式存储），支持动态创建 PV。

---

## 4. **配置与密钥类（Config & Secret）**

管理应用的配置与敏感信息。

* **ConfigMap**

    * 存储非敏感配置信息（配置文件、环境变量）。

* **Secret**

    * 存储敏感数据（密码、证书、API Key）。
    * 数据会做 base64 编码存储。

---

## 5. **集群管理类（Cluster Resources）**

* **Namespace**

    * 提供资源隔离，常用于区分不同环境/团队。

* **Node**

    * 集群中的工作节点资源对象。

* **ResourceQuota**

    * 限制 Namespace 下 CPU、内存、存储等资源用量。

* **LimitRange**

    * 为 Pod/容器设置资源请求和限制的默认值。

* **HorizontalPodAutoscaler (HPA)**

    * 根据指标（如 CPU/内存使用率）自动伸缩 Pod 副本数。

* **VerticalPodAutoscaler (VPA)**

    * 自动调整 Pod 的资源请求值。

* **CustomResourceDefinition (CRD)**

    * 用户自定义资源类型，可以扩展 Kubernetes API。

---

✅ **总结**：

* **工作负载**：运行应用（Deployment、StatefulSet、DaemonSet、Job、CronJob）。
* **网络与服务**：服务发现与流量管理（Service、Ingress、NetworkPolicy）。
* **存储**：持久化和动态卷管理（PV、PVC、StorageClass）。
* **配置管理**：应用配置和敏感信息（ConfigMap、Secret）。
* **集群资源**：命名空间、节点、资源配额、自动伸缩。

---

# kubectl 命令完全指南

## 📁 一、命名空间 (Namespace) 操作

```bash
# 查看所有命名空间
kubectl get namespaces
kubectl get ns

# 创建命名空间
kubectl create namespace <namespace-name>

# 删除命名空间
kubectl delete namespace <namespace-name>

# 设置默认命名空间
kubectl config set-context --current --namespace=<namespace-name>
```

## 🚀 二、Pod 操作

```bash
# 查看Pod
kubectl get pods [-n <namespace>]
kubectl get pods -o wide          # 显示更多信息
kubectl get pods --watch         # 实时监控
kubectl get pods -A              # 所有命名空间

# 查看Pod详情
kubectl describe pod <pod-name> [-n <namespace>]

# 查看Pod日志
kubectl logs <pod-name> [-n <namespace>]
kubectl logs -f <pod-name>       # 实时日志
kubectl logs --tail=100 <pod-name> # 最后100行
kubectl logs -p <pod-name>       # 之前崩溃的容器日志

# 进入Pod执行命令
kubectl exec -it <pod-name> -- /bin/bash
kubectl exec -it <pod-name> -- /bin/sh
kubectl exec <pod-name> -- <command>  # 执行单条命令

# 删除Pod
kubectl delete pod <pod-name> [-n <namespace>]

# 强制删除Pod（卡住时使用）
kubectl delete pod <pod-name> --force --grace-period=0
```

## 📦 三、部署 (Deployment) 操作

```bash
# 查看部署
kubectl get deployments
kubectl get deploy

# 查看部署详情
kubectl describe deployment <deployment-name>

# 创建/更新部署
kubectl apply -f deployment.yaml

# 删除部署
kubectl delete deployment <deployment-name>

# 扩缩容
kubectl scale deployment <deployment-name> --replicas=3

# 滚动重启
kubectl rollout restart deployment <deployment-name>

# 查看部署状态
kubectl rollout status deployment <deployment-name>

# 回滚部署
kubectl rollout undo deployment <deployment-name>
kubectl rollout undo deployment <deployment-name> --to-revision=2

# 查看部署历史
kubectl rollout history deployment <deployment-name>
```

## 🌐 四、服务 (Service) 操作

```bash
# 查看服务
kubectl get services
kubectl get svc

# 查看服务详情
kubectl describe service <service-name>

# 创建服务
kubectl expose deployment <deployment-name> --port=80 --target-port=8080 --type=NodePort

# 端口转发（本地访问服务）
kubectl port-forward service/<service-name> 8080:80
kubectl port-forward pod/<pod-name> 8080:80

# 删除服务
kubectl delete service <service-name>
```

## 🚪 五、Ingress 操作

```bash
# 查看Ingress
kubectl get ingress
kubectl get ing

# 查看Ingress详情
kubectl describe ingress <ingress-name>

# 查看Ingress控制器
kubectl get pods -n ingress-nginx
```

## 💾 六、配置管理 (ConfigMap & Secret)

```bash
# ConfigMap 操作
kubectl get configmaps
kubectl get cm
kubectl describe configmap <configmap-name>
kubectl create configmap <name> --from-file=<file-path>
kubectl create configmap <name> --from-literal=key=value

# Secret 操作
kubectl get secrets
kubectl describe secret <secret-name>
kubectl create secret generic <name> --from-literal=password=secret
```

## 📊 七、状态检查与监控

```bash
# 查看节点状态
kubectl get nodes
kubectl describe node <node-name>

# 查看资源使用情况
kubectl top nodes          # 节点资源使用
kubectl top pods           # Pod资源使用
kubectl top pods --containers # 容器资源使用

# 查看事件
kubectl get events
kubectl get events --sort-by=.metadata.creationTimestamp
kubectl get events -n <namespace>

# 查看资源配额
kubectl get resourcequotas
kubectl get limitranges
```

## 🔧 八、调试与故障排查

```bash
# 查看API资源
kubectl api-resources

# 查看API版本
kubectl api-versions

# 集群信息
kubectl cluster-info
kubectl cluster-info dump

# 查看当前配置
kubectl config view
kubectl config current-context

# 切换上下文
kubectl config use-context <context-name>

# 调试命令
kubectl run debug-pod --image=busybox --rm -it --restart=Never -- /bin/sh
```

## 📋 九、文件操作

```bash
# 从文件创建资源
kubectl apply -f <file.yaml>
kubectl apply -f <directory/>          # 目录下所有文件
kubectl apply -f <file1.yaml> -f <file2.yaml>

# 删除文件定义的资源
kubectl delete -f <file.yaml>

# 验证YAML文件
kubectl apply -f <file.yaml> --dry-run=client
kubectl apply -f <file.yaml> --dry-run=server

# 导出资源配置
kubectl get deployment <name> -o yaml > deployment.yaml
kubectl get pod <name> -o yaml > pod.yaml

# 编辑资源
kubectl edit deployment <deployment-name>
kubectl edit pod <pod-name>
```

## 🔄 十、批量操作

```bash
# 删除所有Pod
kubectl delete pods --all

# 删除命名空间下所有资源
kubectl delete all --all -n <namespace>

# 批量获取资源
kubectl get pods,services,deployments
kubectl get all -n <namespace>

# 标签操作
kubectl get pods --show-labels
kubectl label pods <pod-name> env=prod
kubectl get pods -l app=nginx
```

## 🎯 十一、常用组合命令

```bash
# 一键查看所有
kubectl get pods,services,deployments -n <namespace>

# 持续监控Pod状态
watch kubectl get pods -n <namespace>

# 查看Pod并排序
kubectl get pods --sort-by=.metadata.creationTimestamp

# 查看最近创建的Pod
kubectl get pods --sort-by=.metadata.creationTimestamp | tail -5

# 快速进入第一个Pod
kubectl exec -it $(kubectl get pods -o name | head -1) -- /bin/sh
```

## 📝 十二、实用技巧

```bash
# 使用别名（添加到 ~/.bashrc 或 ~/.zshrc）
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgs='kubectl get services'
alias kgd='kubectl get deployments'
alias kaf='kubectl apply -f'
alias kdf='kubectl delete -f'

# JSON路径查询
kubectl get pods -o jsonpath='{.items[*].metadata.name}'
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.podIP}{"\n"}{end}'

# 自定义输出列
kubectl get pods -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,NODE:.spec.nodeName
```