# dummy size option

{{- define "etcd-dev" }}
              cpu: 500m
{{- end }}

{{- define "etcd-stage" }}
              cpu: 480m
{{- end }}

{{- define "etcd-prod" }}
              cpu: 480m
{{- end }}

{{- define "k8s-api-server-dev" }}
              cpu: 250m
{{- end }}

{{- define "k8s-api-server-stage" }}
              cpu: 480m
{{- end }}

{{- define "k8s-api-server-prod" }}
              cpu: 480m
{{- end }}

{{- define "k8s-api-server-dev-request" }}
              cpu: 250m
              memory: 500Mi
{{- end }}

{{- define "k8s-ctrl-mgr-dev" }}
              cpu: 100m
              memory: 500Mi
{{- end }}

{{- define "k8s-ctrl-mgr-stage" }}
              cpu: 490m
              memory: 512Mi
{{- end }}

{{- define "k8s-ctrl-mgr-prod" }}
              cpu: 490m
              memory: 512Mi
{{- end }}

{{- define "etcd_resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to default
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.etcd }}
              cpu: {{ .Values.global.resources.etcd.cpu }}
              {{- else if eq .Values.global.resources.clustertype "dev" }}
                {{- template "etcd-dev" . }}
              {{- else if eq .Values.global.resources.clustertype "stage" }}
                {{- template "etcd-stage" . }}
              {{- else }}
                {{- template "etcd-prod" . }}
              {{- end }}
            {{- else }}
              {{- template "etcd-prod" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.etcd }}
              cpu: {{  .Values.global.resources.etcd.cpu }}
              {{- else if eq .Values.global.resources.clustertype "dev" }}
                {{- template "etcd-dev" . }}
              {{- else if eq .Values.global.resources.clustertype "stage" }}
                {{- template "etcd-stage" . }}
              {{- else }}
                {{- template "etcd-prod" . }}
              {{- end }}
            {{- else }}
              {{- template "etcd-prod" . }}
            {{- end }}
{{- end }}

{{- define "kube_controllermanager_resources" }}
        resources:
          limits:
          # this is to check if the override value is present if not we will set it to default
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubecontrollermanager }}
            cpu: {{ .Values.global.resources.kubecontrollermanager.cpu }}
            memory: {{ .Values.global.resources.kubecontrollermanager.memory }}
            {{- else if eq .Values.global.resources.clustertype "dev" }}
              {{- template "k8s-ctrl-mgr-dev" . }}
            {{- else if eq .Values.global.resources.clustertype "stage" }}
              {{- template "k8s-ctrl-mgr-stage" . }}
            {{- else }}
              {{- template "k8s-ctrl-mgr-prod" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-ctrl-mgr-prod" . }}
          {{- end }}
          requests:
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubecontrollermanager }}
            cpu: {{ .Values.global.resources.kubecontrollermanager.cpu }}
            memory: {{  .Values.global.resources.kubecontrollermanager.memory }}
            {{- else if eq .Values.global.resources.clustertype "dev" }}
              {{- template "k8s-ctrl-mgr-dev" . }}
            {{- else if eq .Values.global.resources.clustertype "stage" }}
              {{- template "k8s-ctrl-mgr-stage" . }}
            {{- else }}
              {{- template "k8s-ctrl-mgr-prod" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-ctrl-mgr-prod" . }}
          {{- end }}
{{- end }}

{{- define "kube_apiserver_resources" }}
        resources:
          limits:
          # this is to check if the override value is present if not we will set it to default
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubeapiserver }}
            cpu: {{ .Values.global.resources.kubeapiserver.cpu }}
            {{- else if eq .Values.global.resources.clustertype "dev" }}
              {{- template "k8s-api-server-dev" . }}
            {{- else if eq .Values.global.resources.clustertype "stage" }}
              {{- template "k8s-api-server-stage" . }}
            {{- else }}
              {{- template "k8s-api-server-prod" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-api-server-prod" . }}
          {{- end }}
          requests:
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubeapiserver }}
            cpu: {{ .Values.global.resources.kubeapiserver.cpu }}
            {{- else if eq .Values.global.resources.clustertype "dev" }}
              {{- template "k8s-api-server-dev-request" . }}
            {{- else if eq .Values.global.resources.clustertype "stage" }}
              {{- template "k8s-api-server-stage" . }}
            {{- else }}
              {{- template "k8s-api-server-prod" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-api-server-prod" . }}
          {{- end }}
{{- end }}


# dummy size option
{{- define "graphql-dev" }}
                    cpu: 500m
                    memory: 128Mi
{{- end }}

{{- define "graphql-stage" }}
                    cpu: 490m
                    memory: 2Gi
{{- end }}

{{- define "graphql-prod" }}
                    cpu: 490m
                    memory: 2Gi
{{- end }}

{{- define "graphql-dev-request" }}
                    cpu: 10m
                    memory: 64Mi
{{- end }}

{{- define "graphql_resources" }}
                resources:
                  limits:
                  # this is to check if the override value is present if not we will set it to default
                  {{- if .Values.global.resources }}
                    {{- if .Values.global.resources.graphql }}
                        cpu: {{ .Values.global.resources.graphql.cpu }}
                        memory: {{ .Values.global.resources.graphql.memory }}
                    {{- else if eq .Values.global.resources.clustertype "dev" }}
                      {{- template "graphql-dev" . }}
                    {{- else if eq .Values.global.resources.clustertype "stage" }}
                      {{- template "graphql-stage" . }}
                    {{- else }}
                      {{- template "graphql-prod" . }}
                    {{- end }}
                  {{- else }}
                    {{- template "graphql-prod" . }}
                  {{- end }}
                  requests:
                  {{- if .Values.global.resources }}
                    {{- if .Values.global.resources.graphql }}
                        cpu: {{ .Values.global.resources.graphql.cpu }}
                        memory: {{   .Values.global.resources.graphql.memory }}
                    {{- else if eq .Values.global.resources.clustertype "dev" }}
                      {{- template "graphql-dev-request" . }}
                    {{- else if eq .Values.global.resources.clustertype "stage" }}
                      {{- template "graphql-stage" . }}
                    {{- else }}
                      {{- template "graphql-prod" . }}
                    {{- end }}
                  {{- else }}
                    {{- template "graphql-prod" . }}
                  {{- end }}
{{- end }}


{{- define "argoJobHook" }}
    argocd.argoproj.io/hook: Sync
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
{{- end }}

{{- define "tolerations" }}
 {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
{{- include "common.affinity-and-toleration" . | nindent 6 }}
{{- end }}