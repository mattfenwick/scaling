{{- define "loadgen.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}


{{- define "loadgen.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}


{{- define "loadgen.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}


{{- define "loadgen.deployment.labels" -}}
helm.sh/chart: {{ include "loadgen.chart" . }}
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
helm.sh/chart: {{ include "loadgen.chart" . }}
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
