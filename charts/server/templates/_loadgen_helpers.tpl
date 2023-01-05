{{- define "loadgen.deployment.labels" -}}
helm.sh/chart: {{ include "scaling.chart" . }}
{{ include "loadgen.deployment.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: "todo" # "used to be a template: 'include "loadgen.app-version" .'
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{- define "loadgen.deployment.selectorLabels" -}}
app.kubernetes.io/name: loadgen
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: loadgen-deployment
app.kubernetes.io/part-of: loadgen
{{- end }}


{{- define "loadgen.job.labels" -}}
helm.sh/chart: {{ include "scaling.chart" . }}
{{ include "loadgen.job.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: "todo" # "used to be a template: 'include "loadgen.app-version" .'
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{- define "loadgen.job.selectorLabels" -}}
app.kubernetes.io/name: loadgen
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: loadgen-job
app.kubernetes.io/part-of: loadgen
{{- end }}
