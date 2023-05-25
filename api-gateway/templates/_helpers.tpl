# dummy size option
{{- define "api-gw-dev" }}
              cpu: 250m
              memory: 128Mi
{{- end }}

{{- define "api-gw-dev-request" }}
              cpu: 100m
              memory: 128Mi
{{- end }}

{{- define "api-gw-stage" }}
              cpu: 490m
              memory: 512Mi
{{- end }}

{{- define "api-gw-prod" }}
              cpu: 490m
              memory: 512Mi
{{- end }}

{{- define "api_gw_resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to default
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.api_gateway }}
              cpu: {{ .Values.global.resources.api_gateway.cpu }}
              memory: {{ .Values.global.resources.api_gateway.memory }}
              {{- else if eq .Values.global.resources.clustertype "dev" }}
                {{- template "api-gw-dev" . }}
              {{- else if eq .Values.global.resources.clustertype "stage" }}
                {{- template "api-gw-stage" . }}
              {{- else }}
                {{- template "api-gw-prod" . }}
              {{- end }}
            {{- else }}
              {{- template "api-gw-prod" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.api_gateway }}
              cpu: {{ .Values.global.resources.api_gateway.cpu }}
              memory: {{  .Values.global.resources.api_gateway.memory }}
              {{- else if eq .Values.global.resources.clustertype "dev" }}
                {{- template "api-gw-dev-request" . }}
              {{- else if eq .Values.global.resources.clustertype "stage" }}
                {{- template "api-gw-stage" . }}
              {{- else }}
                {{- template "api-gw-prod" . }}
              {{- end }}
            {{- else }}
              {{- template "api-gw-prod" . }}
            {{- end }}
{{- end }}



{{- define "tolerations" }}
 {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
{{- include "common.affinity-and-toleration" . | nindent 6 }}
{{- end }}