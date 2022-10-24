{{- define "loadgen.deployment.labels" -}}
helm.sh/chart: {{ include "scaling.chart" . }}
{{ include "loadgen.deployment.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: "todo" # "used to be a template: 'include "scaling.app-version" .'
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{- define "loadgen.deployment.selectorLabels" -}}
app.kubernetes.io/name: loadgen
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: loadgen
app.kubernetes.io/part-of: scaling
{{- end }}


{{- define "loadgen.job.labels" -}}
helm.sh/chart: {{ include "scaling.chart" . }}
{{ include "loadgen.job.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: "todo" # "used to be a template: 'include "scaling.app-version" .'
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{- define "loadgen.job.selectorLabels" -}}
app.kubernetes.io/name: loadgen
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: loadgen
app.kubernetes.io/part-of: scaling
{{- end }}
