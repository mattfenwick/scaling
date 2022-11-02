{{- define "webserver.labels" -}}
helm.sh/chart: {{ include "scaling.chart" . }}
{{ include "webserver.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: "todo" # "used to be a template: 'include "scaling.app-version" .'
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{- define "webserver.selectorLabels" -}}
app.kubernetes.io/name: webserver
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: webserver
app.kubernetes.io/part-of: scaling
{{- end }}


{{- define "webserver.serviceAccountName" -}}
{{- if .Values.webserver.serviceAccount.create }}
{{- default (include "scaling.fullname" .) .Values.webserver.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.webserver.serviceAccount.name }}
{{- end }}
{{- end }}
