# dummy size option
{{- define "validation-default" }}
              cpu: 250m
              memory: 128Mi
{{- end }}

{{- define "validation-default-request" }}
              cpu: 100m
              memory: 128Mi
{{- end }}

{{- define "small" }}
              cpu: 500m
              memory: 128Mi
{{- end }}

{{- define "argoJobHook" }}
    argocd.argoproj.io/hook: Sync
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
{{- end }}

{{- define "validation_resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to default
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.validation }}
              cpu: {{ .Values.global.resources.validation.cpu }}
              memory: {{ .Values.global.resources.validation.memory }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "validation-default" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.validation }}
              cpu: {{ .Values.global.resources.validation.cpu }}
              memory: {{  .Values.global.resources.validation.memory }}
              {{- else }}
                {{- if eq .Values.global.size "small" }}
              {{- template "small" . }}
                {{- end }}
              {{- end }}
            {{- else }}
              {{- template "validation-default-request" . }}
            {{- end }}
{{- end }}
