# dummy size option
{{- define "graphql-default" }}
              cpu: 500m
              memory: 128Mi
{{- end }}

{{- define "graphql-default-request" }}
              cpu: 10m
              memory: 64Mi
{{- end }}

{{- define "small" }}
              cpu: 500m
              memory: 128Mi
{{- end }}

{{- define "graphql_resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to default
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.graphql }}
              cpu: {{ .Values.global.resources.graphql.cpu }}
              memory: {{ .Values.global.resources.graphql.memory }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "graphql-default" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.graphql }}
              cpu: {{ .Values.global.resources.graphql.cpu }}
              memory: {{   .Values.global.resources.graphql.memory }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "graphql-default-request" . }}
            {{- end }}
{{- end }}
