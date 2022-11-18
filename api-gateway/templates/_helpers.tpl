# dummy size option
{{- define "default" }}
              cpu: 500m
              memory: 128Mi
{{- end }}

{{- define "small" }}
              cpu: 500m
              memory: 128Mi
{{- end }}

{{- define "resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to default
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.api_gateway }}
              cpu: {{ .Values.global.resources.api_gateway.cpu }}
              memory: {{ .Values.global.resources.api_gateway.memory }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "default" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.api_gateway }}
              cpu: {{ .Values.global.resources.api_gateway.cpu }}
              memory: {{  .Values.global.resources.api_gateway.memory }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "default" . }}
            {{- end }}
{{- end }}
