# Kubernetes 日志收集系统

本项目使用 ECK (Elastic Cloud on Kubernetes) 和 Fluent Bit 搭建一个完整的 Kubernetes 日志收集系统。

## 系统架构

```
[K8s Pod 日志] 
      |
      v
[Fluent Bit DaemonSet] (收集、解析、过滤日志)
      |
      | (通过 HTTP/TLS 输出)
      v
[Elasticsearch Cluster (ECK)] (存储和索引日志)
      |
      v
[Kibana (ECK)] (可视化和查询日志)
```

## 目录结构

```
/elastic-system/ 
├── eck/                      # ECK 相关配置 
│   ├── kibana.yaml           # Kibana 配置 
│   ├── operator.yaml         # ECK 操作符配置 
│   └── elasticsearch.yaml    # Elasticsearch 资源定义 
├── fluent-bit/               # Fluent Bit 配置 
│   ├── configmap.yaml        # Fluent Bit 配置映射 
│   ├── daemonset.yaml        # Fluent Bit DaemonSet 
│   └── rbac.yaml             # Fluent Bit RBAC 权限 
├── sample-app/               # 示例应用 
│   └── counter.yaml          # 示例日志生成应用 
└── namespace.yaml            # 命名空间定义 
```

## 安装步骤

### 1. 创建命名空间

```bash
kubectl apply -f namespace.yaml
```

### 2. 安装 ECK 操作符

首先安装 ECK Operator，它将管理 Elasticsearch 和 Kibana。

```bash
# 创建独立的命名空间
kubectl create namespace elastic-system

# 安装 ECK CRDs
kubectl apply -f https://download.elastic.co/downloads/eck/2.10.0/crds.yaml

# 安装 ECK Operator
kubectl apply -f https://download.elastic.co/downloads/eck/2.10.0/operator.yaml

# 确认 Operator 运行正常
kubectl get pods -n elastic-system
```

### 3. 部署 Elasticsearch

```bash
kubectl apply -f eck/elasticsearch.yaml
```

等待 Elasticsearch 部署完成：

```bash
kubectl get elasticsearch -n elastic-system
kubectl get pods -n elastic-system -l elasticsearch.k8s.elastic.co/cluster-name=elasticsearch
```

启动失败可能是内存限制

```bash
kubectl get elasticsearch quickstart -o yaml | grep -A 5 -B 5 resources
```

### 4. 部署 Kibana

```bash
kubectl apply -f eck/kibana.yaml
```

等待 Kibana 部署完成：

```bash
kubectl get kibana -n elastic-system
kubectl get pods -n elastic-system -l kibana.k8s.elastic.co/name=kibana
```

### 5. 部署 Fluent Bit

```bash
# 创建 RBAC 权限
kubectl apply -f fluent-bit/fluent-bit-rbac.yaml

# 创建配置映射
kubectl apply -f fluent-bit/fluent-bit-config.yaml

# 部署 DaemonSet
kubectl apply -f fluent-bit/fluent-bit-daemonset.yaml
```

### 6. 部署示例应用

```bash
kubectl apply -f sample-app/counter.yaml
```

## 访问 Kibana

获取 Elasticsearch 的默认用户密码：

```bash
kubectl get secret elasticsearch-es-elastic-user -n elastic-system -o=jsonpath='{.data.elastic}' | base64 --decode; echo

```

使用端口转发访问 Kibana：

```bash
kubectl port-forward service/kibana-kb-http -n elastic-system 5601:5601
```

在浏览器中访问 https://localhost:5601，使用用户名 `elastic` 和上面获取的密码登录。

## 查看日志

在 Kibana 中：

1. 进入 **Stack Management** > **Index Patterns**
2. 创建索引模式 `fluent-bit-*`
3. 选择 `@timestamp` 作为时间字段
4. 进入 **Discover** 页面查看日志

## 故障排除

### 检查 Fluent Bit 状态

```bash
kubectl get pods -n elastic-system -l k8s-app=fluent-bit
kubectl logs -f -n elastic-system -l k8s-app=fluent-bit
```

### 检查 Elasticsearch 状态

```bash
kubectl get elasticsearch -n elastic-system
kubectl describe elasticsearch -n elastic-system
```

### 检查 Kibana 状态

```bash
kubectl get kibana -n elastic-system
kubectl describe kibana -n elastic-system
```

## 注意事项

- 本配置针对 kind 本地集群进行了资源优化，减少了存储和内存需求
- 生产环境中应根据实际需求调整资源配置
- Elasticsearch 和 Kibana 的版本应保持一致