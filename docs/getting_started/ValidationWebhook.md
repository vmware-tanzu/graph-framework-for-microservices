## Creating a Validation webhook in the NexusAPIServer

This Workflow will provide steps to add a validation webhook with custom microservice on a nexus runtime

### Prerequisites

1. Please install nexus runtime in a namespace using nexus CLI
   
   #### **[Nexus Runtime installation](RuntimeWorkflow.md)**


2. Create a microservice with validation webhook guidelines provided by K8s 

   #### **[Writing a backend service](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#write-an-admission-webhook-server)**
   #### **[Writing a backend service using golang SDK](https://developers.redhat.com/articles/2021/09/17/kubernetes-admission-control-validating-webhooks#bootstrap_with_the_operator_sdk)**



### Steps 

1. Create self signed certificates to use for TLS/SSL termination for communication from APIServer to Microservice

```
#!/usr/bin/env bash
set -ex
namespace=$POD_NAMESPACE
usage() {
    cat <<EOF
Generate certificate suitable for use with an webhook service.

This script uses k8s' CertificateSigningRequest API to a generate a
certificate signed by k8s CA suitable for use with sidecar-injector webhook
services. This requires permissions to create and approve CSR. See
https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster for
detailed explantion and additional instructions.

The server key/cert k8s CA cert are stored in a k8s secret.

usage: ${0} [OPTIONS]

The following flags are required.

        --service          Service name of webhook.
EOF
    exit 1
}

while [[ $# -gt 0 ]]; do
    case ${1} in
        --service)
            service="$2"
            shift
            ;;
        *)
            usage
            ;;
    esac
    shift
done

if [ ! -x "$(command -v openssl)" ]; then
    echo "openssl not found"
    exit 1
fi

csrName=${service}
tmpdir=$(mktemp -d)
echo "creating certs in tmpdir ${tmpdir} "

cat <<EOF >> ${tmpdir}/csr.conf
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${service}
DNS.2 = ${service}.${namespace}
DNS.3 = ${service}.${namespace}.svc
DNS.4 = localhost
EOF

openssl genrsa -out ${tmpdir}/server-key.pem 2048
openssl req -new -x509 -nodes -days 365000 -key  ${tmpdir}/server-key.pem -out ${tmpdir}/ca-cert.pem -subj "/C=US/ST=Denial/L=PaloAlto/O=Dis/CN=nexus"
openssl req -new -key ${tmpdir}/server-key.pem -config ${tmpdir}/csr.conf -subj "/CN=${service}" -out ${tmpdir}/server.csr
openssl x509 -req -days 365 -in ${tmpdir}/server.csr -CAcreateserial -sha256 -out ${tmpdir}/server.crt -CA ${tmpdir}/ca-cert.pem -CAkey ${tmpdir}/server-key.pem
openssl x509 -in ${tmpdir}/server.crt -out ${tmpdir}/server-cert.pem -outform PEM

```

save the above script as create_self_signed_cert.sh and run script with name of backend service

```
bash create_self_signed_cert.sh --service <backend-service-name>
```

Note: Please create self signed certificate if your service is exposed at HTTP, For HTTPS service , please use the Certificate already generated.

3. Create validation webhook configuration on nexus APIServer

```
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
name: "nexus-validation.webhook.svc"
webhooks:
- name: "nexus-validation-crd.webhook.svc"
    failurePolicy: Fail
    rules:
    - apiGroups: [ "api.nexus.vmware.com" ] ---> Add required apiGroups(CRD group/base groups)
        apiVersions: [ "v1" ]
        operations: [ "CREATE" ]  ---> Add the events you want to intercept
        resources: [ "*" ]
        scope: "*"
    clientConfig:
    url: https://nexus-validation/validate
    caBundle: __CA_BUNDLE__   ---> Please replace the content with ca-cert.pem contents
    admissionReviewVersions: [ "v1", "v1beta1" ]
    sideEffects: None
```

Edit the above yaml appropriate to your Validation webhook config and create on nexus Runtime

```
kubectl port-forward -n $NAMESPACE -lcontrol-plane=api-gw 5000:80 &
kubectl -s localhost:5000 apply -f <webhook>.yaml
```

