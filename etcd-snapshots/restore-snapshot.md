# Steps to restore etcd snapshot

1. Take a etcd backup if dont't find required snapshots [here](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-calibration/-/tree/master/etcd-snapshots)
```shell
kubectl port-forward svc/nexus-etcd 2379:2379 -n <Namespace> 
ETCDCTL_API=3 etcdctl --endpoints http://localhost:2379 snapshot save <snapshot-name>
```

2. Create a Backup PVC and etcd-backup-pod manifest
```shell
echo 'kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: etcd-backup-pvc-2
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
  storageClassName: nfs
---
apiVersion: v1
kind: Pod
metadata:
  name: etcd-backup-pod-2
spec:
  volumes:
    - name: etcd-backup
      persistentVolumeClaim:
        claimName: etcd-backup-pvc-2
  containers:
    - name: inspector
      image: busybox:latest
      command:
        - sleep
        - infinity
      volumeMounts:
        - mountPath: "/backup"
          name: etcd-backup' > etcd-mount.yaml
```

3. Copy the snapshot to etcd-backup-pod
```shell
kubectl cp <snapshot-name> <Namespace>/etcd-backup-pod-2:/backup/<snapshot-name>
```

4. Change backup file permissions
```shell
kubectl exec -it etcd-backup-pod-2 -n <Namespace> -- sh
chown -R 1001 /backup/<snapshot-name>
chmod -R 700 /backup/<snapshot-name>
```

5. Delete the current etcd statefulsets and services
```shell
kubectl delete statefulset nexus-etcd  -n <Namespace>
kubectl delete svc nexus-etcd nexus-etcd-headless -n <Namespace>
kubectl delete pvc -n <Namespace> data-nexus-etcd-0
```

6. Start etcd with snapshot
```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install nexus-etcd bitnami/etcd \
  --set startFromSnapshot.enabled=true \
  --set startFromSnapshot.existingClaim=etcd-backup-pvc-2 \
  --set startFromSnapshot.snapshotFilename=mybackup \
  --set image.debug=true \
  --set containerSecurityContext.runAsUser=1001 \
  --set containerSecurityContext.enabled=true \
  --set containerSecurityContext.allowPrivilegeEscalation=true \
  --set containerSecurityContext.runAsNonRoot=false \
  --set statefulset.replicaCount=1 --namespace <Namespace>
```

7. Update the caBundle value in validatingwebhookconfigurations
```shell
a. kubectl port-forward svc/nexus-api-gw 5000:80 -n <Namespace>
b. kubectl get secrets nexus-validation-tls -n <Namspace> -o yaml | yq '.data.["ca.crt"]'
c. kubectl -s localhost:5000 edit validatingwebhookconfigurations
Update `caBundle` value from above step 7b output
```

8. Update the nginx-ingress service in api-gateway
```shell
kubectl -s localhost:5000 edit svc nginx-ingress
Change this value to --> externalName: nexus-ingress-nginx-controller.<Namespace>.svc
```


