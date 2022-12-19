# dummy size option

{{- define "etcd-default" }}
              cpu: 500m
{{- end }}

{{- define "k8s-api-server-default" }}
              cpu: 500m
{{- end }}

{{- define "k8s-api-server-default-request" }}
              cpu: 250m
              memory: 500Mi
{{- end }}

{{- define "k8s-ctrl-mgr-default" }}
              cpu: 100m
              memory: 500Mi
{{- end }}

{{- define "small" }}
              cpu: 500m
              memory: 128Mi
{{- end }}

{{- define "etcd_resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to default
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.etcd }}
              cpu: {{ .Values.global.resources.etcd.cpu }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "etcd-default" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.etcd }}
              cpu: {{  .Values.global.resources.etcd.cpu }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
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
            {{- else }}
              {{- if eq .Values.global.size "small" }}
          {{- template "small" . }}
              {{- end }}
            {{- end }}
          {{- else }}
          {{- template "k8s-ctrl-mgr-default" . }}
          {{- end }}
          requests:
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubecontrollermanager }}
            cpu: {{ .Values.global.resources.kubecontrollermanager.cpu }}
            memory: {{  .Values.global.resources.kubecontrollermanager.memory }}
            {{- else }}
              {{- if eq .Values.global.size "small" }}
            {{- template "small" . }}
              {{- end }}
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
            {{- else }}
              {{- if eq .Values.global.size "small" }}
          {{- template "small" . }}
             {{- end }}
            {{- end }}
          {{- else }}
          {{- template "k8s-api-server-default" . }}
          {{- end }}
          requests:
          {{- if .Values.global.resources }}
            {{- if .Values.global.resources.kubeapiserver }}
            cpu: {{ .Values.global.resources.kubeapiserver.cpu }}
            {{- else }}
              {{- if eq .Values.global.size "small" }}
            {{- template "small" . }}
              {{- end }}
            {{- end }}
          {{- else }}
            {{- template "k8s-api-server-default-request" . }}
          {{- end }}
{{- end }}
