{{- define "postgres.environment" }}
- name: PGHOST
  value: {{ .Values.postgres.host | quote }}
- name: PGPORT
  value: "5432"
- name: PGUSER
  value: {{ .Values.postgres.user | quote }}
- name: PGPASSWORD
  value: {{ .Values.postgres.password | quote }}
- name: PGDATABASE
  value: {{ .Values.postgres.dbname | quote }}
{{- end }}
