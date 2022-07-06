# Connector

1. Supports on-demand replication of state between two endpoints.
2. Watches for ReplicationConfig object created in the Nexus DM and learns the source node and remote endpoint information.
3. Replicates the source object and its immediate children and NOT the entire sub-graph. 
4. Replicates only the spec and ignores the relationships. 

In the below-example, if user has enabled replication for the object "Cluster Bar", then connector syncs ClusterBar and its immediate children GNSBinding Bar and Policy A. 

![Connector](.content/images/connector.png)

## Connect Datamodel

![ConnectDatamodel](.content/images/connect-dm.png)

## Sample NexusEndpoint Object

```
apiVersion: connect.nexus.org/v1
kind: NexusEndpoint
metadata:
  name: default
  labels:
     nexus/is_name_hashed: "false"
     connects.connect.nexus.org: default
spec:
    host: https://8DBD4DFA95D3F56BD5A31D1090E6D244.gr7.us-east-2.eks.amazonaws.com
    port: "443"
    cert: xyz
```

where **host** indicates destination cluster URL. 

## Sample ReplicationConfig Object

```
apiVersion: connect.nexus.org/v1
kind: ReplicationConfig
metadata:
  name: defaultnew
  labels:
     nexus/is_name_hashed: "false"
     connects.connect.nexus.org: default
spec:
  accessToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IlZaMVRDcHc0bEVMbjJ5eHNZbDRfV2o1VGxhNXdLQllNRnJCN09qQlB3T2sifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tcmh4ZGQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjIwZGFhZDRhLTJjMjItNDgwNy1hNDkzLWQzMWIyZjQ5ZGZkZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.GQzwvDJpIF7oxWL_5gkvJFK9n-WPfClZjjmA7KqrZBf7JEqRJI4AxjnGKsxtOQyngvAPj98S6GyYwk7ZVVWfthml5o6JtPwabsSSTFStL8xljoLbq_JbeBNIAxcw4XT-5wL7w2J-gi3-XfTkY_FZbUjIFu4ckeSQV-l_evpigW67wE2sUGUz0l1b48FzmxkfEiQx6Fki4SPHPo0qUJ3ttjNwqLh6AAQaeGonxxkO76utvIsvYoJsOCQNVhsI1pBfcEOuL-PfohLvzOCaGpWEMM-RqHe4j8AUCKslb8N0fiUMQVWhyrQCw5F5hnAGQyJWrCnOpS4g16ZbBfRUYl0i0g"
  remoteEndpointGvk:
    group: connect.nexus.org
    kind: NexusEndpoint
    name: default
  source:
    group: config.mazinger.com
    kind: ApiDevSpace
    name: c6d587ef6e09d12ef72b2caacf12885f97a28b6c
```

where,
- **accessToken:** base64 decoded value of the service account access token present in the K8s secret on the destination cluster.
- **source:** indicates the replication object's information that the user is interested in.
- **remoteEndpointGvk:** indicates link to the desired destination object information.

## How to fetch the access token

**Step 1:** In the destination cluster, fetch the secret info from the service account.

```
% kubectl get sa default -oyaml
apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: "2022-07-04T10:03:37Z"
  name: default
  namespace: default
  resourceVersion: "416"
  uid: 20daad4a-2c22-4807-a493-d31b2f49dfdd
secrets:
- name: default-token-rhxdd
```

**Step 2:** Fetch the token from the secret.

```
% kubectl get secret default-token-rhxdd -oyaml
apiVersion: v1
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeU1EY3dOREV3TURNd01Gb1hEVE15TURjd01URXdNRE13TUZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTXVyCmRlWUlqdlUvVmhUVnBlYVZsK1pya09pM3dPcHZCdDNBSUhPTVJPN2U5WHVaTmpNaU85T0d2T0ZQYUw4QjFWNmIKb1ZRMTEwYTlZWEdVYWNKbEFIZnVzenN0NU1DWXJkTER3NUdiNXhsK1FVazBpdS94cXk2SEVSMW93ZDVlN0xWUQplUGFTbGZVdURjeUNDeWZXbm9pMU9Ca3dkb3RUSlFUU1o5Si9zQ05sWmMwc2ZFblFzdmovM0dHaWhaeWdtQzBwCjFmSzliOVhQa3MvM29DSklaU2dkdkdwUytjTC9kcm1mYkxXakVwRFZJYitDY0JIRkg4dzZyTEtvWWxKalNYQkwKK0d1RXNlWTlaMGVOTWhFUXVFYWluOW1NMWxnNldFQU1ZRjhJNzFFUWV0NFFwWlJxYkJjMUkyK0NUYnMvVnFzMQpxWTZSRWo2WlNkM252YjdIVHJVQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZQMTlzSFpNa0dWMUY1S3c5OUlRcDVxa0h2R0NNQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFERDBjbWJmZzczdFVlSlVoenNNeXF0dlQ0aStya0hSbTE0SHovVlA2NlNlR2Z4MytSdwpGemVIcmI1U2FONHVnT0daM3JNcGVCSTJPRnVZOGo1ZTM0RkZmRWJXQmVxVURGSjNNbGlTRWwyQ3NqODJVb20xCjB6ZG1aZlgxN1c2ckN1ZDc1akFrb1NYVk92eUdHSjZwK3JEVDBWT3ZEVXBaOE1nNm5TeDZTcDZkbG96NjFWNGQKcGl3RUI3dG9WZUg1QkFCVTBtV2M0ejU1aE1XeGhoWE9iZFJRWUJGbEFQUGcwS2V4OFZUNjdOZE4ya00yYlY5RgpXV0RhdlhFanhBK1cyUXdEWFlmWFdCbGJXbXp0VC92SU9qN3hjNFlyZlVUNmw0YVZMd2U1MWdZeEdLNVZCcVZpClFIY3Y4dVRkeEVUODRQWUx4UmVRbk5zTmhFWUdDdHVwVmRhRQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  namespace: ZGVmYXVsdA==
  token: ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklsWmFNVlJEY0hjMGJFVk1iako1ZUhOWmJEUmZWMm8xVkd4aE5YZExRbGxOUm5KQ04wOXFRbEIzVDJzaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUprWldaaGRXeDBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpXTnlaWFF1Ym1GdFpTSTZJbVJsWm1GMWJIUXRkRzlyWlc0dGNtaDRaR1FpTENKcmRXSmxjbTVsZEdWekxtbHZMM05sY25acFkyVmhZMk52ZFc1MEwzTmxjblpwWTJVdFlXTmpiM1Z1ZEM1dVlXMWxJam9pWkdWbVlYVnNkQ0lzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJakl3WkdGaFpEUmhMVEpqTWpJdE5EZ3dOeTFoTkRrekxXUXpNV0l5WmpRNVpHWmtaQ0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwa1pXWmhkV3gwT21SbFptRjFiSFFpZlEuR1F6d3ZESnBJRjdveFdMXzVna3ZKRks5bi1XUGZDbFpqam1BN0txclpCZjdKRXFSSkk0QXhqbkdLc3h0T1F5bmd2QVBqOThTNkd5WXdrN1pWVldmdGhtbDVvNkp0UHdhYnNTU1RGU3RMOHhsam9MYnFfSmJlQk5JQXhjdzRYVC01d0w3dzJKLWdpMy1YZlRrWV9GWmJVaklGdTRja2VTUVYtbF9ldnBpZ1c2N3dFMnNVR1V6MGwxYjQ4RnpteGtmRWlReDZGa2k0U1BIUG8wcVVKM3R0ak53cUxoNkFBUWFlR29ueHhrTzc2dXR2SXN2WW9Kc09DUU5WaHNJMXBCZmNFT3VMLVBmb2hMdnpPQ2FHcFdFTU0tUnFIZTRqOEFVQ0tzbGI4TjBmaVVNUVZXaHlyUUN3NUY1aG5BR1F5SldyQ25PcFM0ZzE2WmJCZlJVWWwwaTBn
kind: Secret
metadata:
  annotations:
    kubernetes.io/service-account.name: default
    kubernetes.io/service-account.uid: 20daad4a-2c22-4807-a493-d31b2f49dfdd
  creationTimestamp: "2022-07-04T10:03:37Z"
  name: default-token-rhxdd
  namespace: default
  resourceVersion: "413"
  uid: 03ab6040-2bcb-4293-9231-f8dc5c83c734
type: kubernetes.io/service-account-token
```

## Perform the following steps to deploy nexus-connector

**Step 1:** Apply nexus-connector deployment and configuration yaml files.

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: connector
  namespace: mazinger
  labels:
    control-plane: connector
spec:
  selector:
    matchLabels:
      control-plane: connector
  replicas: 1
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/default-container: manager
      labels:
        control-plane: connector
    spec:
      initContainers:
        - name: check-nexus-proxy-container
          image: "gcr.io/mesh7-public-images/tools:latest"
          command:
            - /bin/bash
            - -c
            - |
              #!/bin/bash
              set -x
              URL="http://nexus-proxy-container/api/v1/namespaces"
              max_retries=20
              counter=0
              while [[ $counter -lt $max_retries ]]; do
                    status=$(curl -s -o /dev/null -I -w "%{http_code}" -XGET $URL)
                    if [ $status == "200" ]; then
                        echo "$URL is reachable"
                        exit 0
                    else
                        counter=$((counter +1))
                        sleep 5
                    fi
              done
      containers:
        -
          image: 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus/connector:1cae9251381d4eb299daa05899a134fae6ffaa79
          name: connector
          env:
            - name: KUBECONFIG
              value: /config/kubeconfig
          imagePullPolicy: IfNotPresent
          securityContext:
            allowPrivilegeEscalation: false
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8081
            initialDelaySeconds: 15
            periodSeconds: 20
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8081
            initialDelaySeconds: 5
            periodSeconds: 10
          # TODO(user): Configure the resources accordingly based on the project requirements.
          # More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
          resources:
            limits:
              cpu: 500m
              memory: 128Mi
            requests:
              cpu: 10m
              memory: 64Mi
          volumeMounts:
            - mountPath: /config
              name: config
      volumes:
        - name: config
          configMap:
            name: connector-kubeconfig-local
      terminationGracePeriodSeconds: 10
      securityContext:
        runAsUser: 0
        runAsGroup: 0
---
apiVersion: v1
kind: Service
metadata:
  name: nexus-connector
  namespace: mazinger
spec:
  ports:
  - protocol: TCP
    name: http
    port: 80
    targetPort: 80
  - protocol: TCP
    name: https
    port: 443
    targetPort: 443
  selector:
    control-plane: connector
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: connector-kubeconfig-local
  namespace: mazinger
data:
  kubeconfig: |
    current-context: localapiserver
    apiVersion: v1
    kind: Config
    clusters:
    - cluster:
        api-version: v1
        server: http://nexus-proxy-container:80
        insecure-skip-tls-verify: true
      name: localapiserver
    contexts:
    - context:
        cluster: localapiserver
      name: localapiserver
  connector-config: |
    dispatcher:
        workerTTL: "15s"
        maxWorkerCount: 100
        closeRequestsQueueSize: 15
        eventProcessedQueueSize: 100
    ignoredNamespaces:
        matchNames:
            - "kube-public"
            - "kube-system"
            - "kube-node-lease"
            - "istio-system"
            - "ibm-system"
            - "ibm-operators"
            - "ibm-cert-store"
```

**Note:** For kind cluster, use harbor image:
```
harbor-repo.vmware.com/nexus/connector: 1cae9251381d4eb299daa05899a134fae6ffaa79
```

**Step 2:** Apply NexusEndpoint and ReplicationConfig manually.

**Step 3:** Start creating the objects and verify the connector logs if the replication was successful. 

```
% kubectl logs -n mazinger      connector-569998cf98-9n8zf -f
Default container name "manager" not found in a pod
1.6570129302043152e+09	INFO	controller-runtime.metrics	Metrics server is starting to listen	{"addr": ":8080"}
1.657012930204715e+09	INFO	setup	starting manager
1.6570129302054465e+09	INFO	Starting server	{"path": "/metrics", "kind": "metrics", "addr": "[::]:8080"}
1.6570129302055562e+09	INFO	Starting server	{"kind": "health probe", "addr": "[::]:8081"}
1.6570129302059379e+09	INFO	controller.customresourcedefinition	Starting EventSource	{"reconciler group": "apiextensions.k8s.io", "reconciler kind": "CustomResourceDefinition", "source": "kind source: *v1.CustomResourceDefinition"}
1.6570129302059662e+09	INFO	controller.customresourcedefinition	Starting Controller	{"reconciler group": "apiextensions.k8s.io", "reconciler kind": "CustomResourceDefinition"}
time="2022-07-05T09:22:10Z" level=info msg="Starting Controllers..."
time="2022-07-05T09:22:10Z" level=info msg="Starting ReplicationConfig Controller..."
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-replicationconfigs.connect.nexus.org
1.6570129303069077e+09	INFO	controller.customresourcedefinition	Starting workers	{"reconciler group": "apiextensions.k8s.io", "reconciler kind": "CustomResourceDefinition", "worker count": 1}
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-oidcs.authentication.nexus.org
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-workloadruntimes.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apilos.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-projects.apix.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Received Create event for Replication Config defaultnew"
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apiruntimes.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apis.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-workloadbindings.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-workloads.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-testdatas.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-workloadruntimes.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-workloadnamespaces.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-nexuses.api.nexus.org
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-workloadnamespaces.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apispecversions.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apidevspaces.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apiversions.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apinamespaces.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apiversions.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-roots.apix.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-contracttests.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-clusters.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-gatewayruntimes.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-runtimes.apix.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apitests.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apis.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apigateways.apigateway.nexus.org
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-clusters.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apiversions.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-namespaces.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-workloads.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-inventories.apix.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-gatewayruntimes.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apis.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-configs.apix.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-mocks.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-namespaces.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apilos.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-contracts.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apidevspaces.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apispecifications.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-routes.route.nexus.org
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apicollaborationspaces.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apibindings.config.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apinamespaces.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apicollaborationspaces.inventory.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-apitestresults.runtime.mazinger.com
time="2022-07-05T09:22:10Z" level=info msg="Starting controller" controller=controller-integrations.config.mazinger.com
time="2022-07-05T09:51:52Z" level=info msg="Replication is enabled for the resource: roots.apix.mazinger.com/root1/projects.apix.mazinger.com/p1/configs.apix.mazinger.com/c1/apicollaborationspaces.config.mazinger.com/ac1/apidevspaces.config.mazinger.com/ad1, creating..."


time="2022-07-05T09:52:24Z" level=info msg="Replication is enabled for the resource: roots.apix.mazinger.com/root1/projects.apix.mazinger.com/p1/configs.apix.mazinger.com/c1/apicollaborationspaces.config.mazinger.com/ac1/apidevspaces.config.mazinger.com/ad1/integrations.config.mazinger.com/in1, creating..."
time="2022-07-05T09:52:24Z" level=info msg="Replication is enabled for the resource: roots.apix.mazinger.com/root1/projects.apix.mazinger.com/p1/configs.apix.mazinger.com/c1/apicollaborationspaces.config.mazinger.com/ac1/apidevspaces.config.mazinger.com/ad1, updating..."
time="2022-07-05T09:52:29Z" level=info msg="Replication is enabled for the resource: roots.apix.mazinger.com/root1/projects.apix.mazinger.com/p1/configs.apix.mazinger.com/c1/apicollaborationspaces.config.mazinger.com/ac1/apidevspaces.config.mazinger.com/ad1/integrations.config.mazinger.com/in2, creating..."
time="2022-07-05T09:52:29Z" level=info msg="Replication is enabled for the resource: roots.apix.mazinger.com/root1/projects.apix.mazinger.com/p1/configs.apix.mazinger.com/c1/apicollaborationspaces.config.mazinger.com/ac1/apidevspaces.config.mazinger.com/ad1, updating..."
time="2022-07-05T09:53:46Z" level=info msg="Replication is enabled for the resource: roots.apix.mazinger.com/root1/projects.apix.mazinger.com/p1/configs.apix.mazinger.com/c1/apicollaborationspaces.config.mazinger.com/ac1/apidevspaces.config.mazinger.com/ad1/integrations.config.mazinger.com/in2, updating..."
```
