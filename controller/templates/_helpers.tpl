# dummy size option
{{- define "controller-default" }}
              cpu: 100m
              memory: 128Mi
{{- end }}

{{- define "controller-default-request" }}
              cpu: 50m
              memory: 128Mi
{{- end }}

{{- define "small" }}
              cpu: 500m
              memory: 128Mi
{{- end }}

{{- define "controller_resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to default
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.nexus_controller }}
              cpu: {{ .Values.global.resources.nexus_controller.cpu }}
              memory: {{ .Values.global.resources.nexus_controller.memory }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "controller-default" . }}
            {{- end }}
            requests:
              {{- if .Values.global.resources }}
                {{- if .Values.global.resources.nexus_controller }}
              cpu: {{ .Values.global.resources.nexus_controller.cpu }}
              memory: {{  .Values.global.resources.nexus_controller.memory }}
                {{- else }}
                  {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                  {{- end }}
                {{- end }}
              {{- else }}
              {{- template "controller-default-request" . }}
              {{- end }}
{{- end }}
