## 使用Kind搭建本地Kubernetes集群

https://kind.sigs.k8s.io/

Kind 是一个使用 Docker 容器作为“节点”来运行本地 Kubernetes 集群的工具。它不需要虚拟机，因此速度极快，资源消耗极低。

* 控制平面 (Master) 运行在一个容器中。
* 工作节点 (Worker) 也运行在容器中。
* 整个集群都封装在您本机的 Docker 引擎中。

### 步骤一：安装 Kind

```bash
 go install sigs.kubernetes.io/kind@v0.30.0
```

### 步骤二：创建集群

- 创建单节点集群

```bash
 kind create cluster
```

- 创建多节点集群

```bash
kind create cluster --name kind --config kind-multi-node.yaml
```

- 删除集群

```bash
 kind delete cluster --name kind
```

- 查看集群列表

```bash
 kind get clusters
```

### 步骤三：验证集群

```bash
 kubectl cluster-info --context kind-kind
```
