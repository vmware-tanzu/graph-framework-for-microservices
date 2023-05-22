# dummy size option

{{- define "etcd-default" }}
              cpu: 500m
{{- end }}

{{- define "etcd-prod" }}
              cpu: 480m
{{- end }}

{{- define "k8s-api-server-default" }}
              cpu: 250m
{{- end }}

{{- define "k8s-api-server-prod" }}
              cpu: 480m
{{- end }}

{{- define "k8s-api-server-default-request" }}
              cpu: 250m
              memory: 500Mi
{{- end }}

{{- define "k8s-ctrl-mgr-default" }}
              cpu: 100m
              memory: 500Mi
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
              {{- else if eq .Values.global.resources.nexussizing "prod" }}
                {{- template "etcd-prod" . }}
              {{- else }}
                {{- template "etcd-default" . }}
              {{- end }}
            {{- else }}
              {{- template "etcd-default" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.etcd }}
              cpu: {{  .Values.global.resources.etcd.cpu }}
              {{- else if eq .Values.global.resources.nexussizing "prod" }}
                {{- template "etcd-prod" . }}
              {{- else }}
                {{- template "etcd-default" . }}
              {{- end }}
            {{- else }}
              {{- template "etcd-default" . }}
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
            {{- else if eq .Values.global.resources.nexussizing "prod" }}
              {{- template "k8s-ctrl-mgr-prod" . }}
            {{- else }}
              {{- template "k8s-ctrl-mgr-default" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-ctrl-mgr-default" . }}
          {{- end }}
          requests:
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubecontrollermanager }}
            cpu: {{ .Values.global.resources.kubecontrollermanager.cpu }}
            memory: {{  .Values.global.resources.kubecontrollermanager.memory }}
            {{- else if eq .Values.global.resources.nexussizing "prod" }}
              {{- template "k8s-ctrl-mgr-prod" . }}
            {{- else }}
              {{- template "k8s-ctrl-mgr-default" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-ctrl-mgr-default" . }}
          {{- end }}
{{- end }}

{{- define "kube_apiserver_resources" }}
        resources:
          limits:
          # this is to check if the override value is present if not we will set it to default
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubeapiserver }}
            cpu: {{ .Values.global.resources.kubeapiserver.cpu }}
            {{- else if eq .Values.global.resources.nexussizing "prod" }}
              {{- template "k8s-api-server-prod" . }}
            {{- else }}
              {{- template "k8s-api-server-default" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-api-server-default" . }}
          {{- end }}
          requests:
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubeapiserver }}
            cpu: {{ .Values.global.resources.kubeapiserver.cpu }}
            {{- else if eq .Values.global.resources.nexussizing "prod" }}
              {{- template "k8s-api-server-prod" . }}
            {{- else }}
              {{- template "k8s-api-server-default-request" . }}
            {{- end }}
          {{- else }}
            {{- template "k8s-api-server-default-request" . }}
          {{- end }}
{{- end }}


# dummy size option
{{- define "graphql-default" }}
                    cpu: 500m
                    memory: 128Mi
{{- end }}

{{- define "graphql-prod" }}
                    cpu: 490m
                    memory: 2Gi
{{- end }}

{{- define "graphql-default-request" }}
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
                    {{- else if eq .Values.global.resources.nexussizing "prod" }}
                      {{- template "graphql-prod" . }}
                    {{- else }}
                      {{- template "graphql-default" . }}
                    {{- end }}
                  {{- else }}
                    {{- template "graphql-default" . }}
                  {{- end }}
                  requests:
                  {{- if .Values.global.resources }}
                    {{- if .Values.global.resources.graphql }}
                        cpu: {{ .Values.global.resources.graphql.cpu }}
                        memory: {{   .Values.global.resources.graphql.memory }}
                    {{- else if eq .Values.global.resources.nexussizing "prod" }}
                      {{- template "graphql-prod" . }}
                    {{- else }}
                      {{- template "graphql-default-request" . }}
                    {{- end }}
                  {{- else }}
                    {{- template "graphql-default-request" . }}
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