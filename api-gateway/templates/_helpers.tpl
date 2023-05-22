# dummy size option
{{- define "api-gw-default" }}
              cpu: 250m
              memory: 128Mi
{{- end }}

{{- define "api-gw-default-request" }}
              cpu: 100m
              memory: 128Mi
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
              {{- else if eq .Values.global.resources.nexussizing "prod" }}
                {{- template "api-gw-prod" . }}
              {{- else }}
                {{- template "api-gw-default" . }}
              {{- end }}
            {{- else }}
              {{- template "api-gw-default" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.api_gateway }}
              cpu: {{ .Values.global.resources.api_gateway.cpu }}
              memory: {{  .Values.global.resources.api_gateway.memory }}
              {{- else if eq .Values.global.resources.nexussizing "prod" }}
                {{- template "api-gw-prod" . }}
              {{- else }}
                {{- template "api-gw-default-request" . }}
              {{- end }}
            {{- else }}
              {{- template "api-gw-default-request" . }}
            {{- end }}
{{- end }}
