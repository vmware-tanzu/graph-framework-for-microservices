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
              cpu: {{ default "500m" .Values.global.resources.nexus-controller.cpu }}
              memory: {{ default "128Mi" .Values.global.resources.nexus-controller.memory }}
            {{- else }}
            {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
            {{- else }}
              {{- template "default" . }}
            {{- end }}
            {{- end }}
            requests:
              {{- if .Values.global.resources }}
              cpu: {{ default "500m" .Values.global.resources.nexus-controller.cpu "500m" }}
              memory: {{  default "128Mi"  .Values.global.resources.nexus-controller.memory }}
              {{- else }}
              {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
              {{- else }}
              {{- template "default" . }}
            {{- end }}
              {{- end }}
{{- end }}