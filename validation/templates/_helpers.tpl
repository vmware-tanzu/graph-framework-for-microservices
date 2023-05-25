# dummy size option
{{- define "validation-dev" }}
              cpu: 250m
              memory: 128Mi
{{- end }}

{{- define "validation-dev-request" }}
              cpu: 100m
              memory: 128Mi
{{- end }}

{{- define "validation-stage" }}
              cpu: 490m
              memory: 480Mi
{{- end }}

{{- define "validation-prod" }}
              cpu: 490m
              memory: 480Mi
{{- end }}

{{- define "argoJobHook" }}
    argocd.argoproj.io/hook: Sync
    argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
{{- end }}

{{- define "validation_resources" }}
          resources:
            limits:
            # this is to check if the override value is present if not we will set it to prod
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.validation }}
              cpu: {{ .Values.global.resources.validation.cpu }}
              memory: {{ .Values.global.resources.validation.memory }}
              {{- else if eq .Values.global.resources.clustertype "dev" }}
                {{- template "validation-dev" . }}
              {{- else if eq .Values.global.resources.clustertype "stage" }}
                {{- template "validation-stage" . }}
              {{- else }}
                {{- template "validation-prod" . }}
              {{- end }}
            {{- else }}
              {{- template "validation-prod" . }}
            {{- end }}
            requests:
            {{- if .Values.global.resources }}
              {{- if .Values.global.resources.validation }}
              cpu: {{ .Values.global.resources.validation.cpu }}
              memory: {{  .Values.global.resources.validation.memory }}
              {{- else if eq .Values.global.resources.clustertype "dev" }}
                {{- template "validation-dev-request" . }}
              {{- else if eq .Values.global.resources.clustertype "stage" }}
                {{- template "validation-stage" . }}
              {{- else }}
                {{- template "validation-prod" . }}
              {{- end }}
            {{- else }}
              {{- template "validation-prod" . }}
            {{- end }}
{{- end }}

{{- define "tolerations" }}
 {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
{{- include "common.affinity-and-toleration" . | nindent 6 }}
{{- end }}