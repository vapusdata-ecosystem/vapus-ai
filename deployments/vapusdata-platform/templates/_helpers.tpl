{{/*
Expand the name of the chart.
*/}}
{{- define "vapusdata-platform.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 70 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "vapusdata-platform.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 70 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "vapusdata.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "vapusdata.platform.labels" -}}
app.kubernetes.io/part-of: VapusData
app.kubernetes.io/component: platform
{{- end }}

{{- define "vapusdata.platform.selectorLabels" -}}
app.kubernetes.io/component: platform
{{- end }}

{{- define "vapusdata.nabhikserver.labels" -}}
app.kubernetes.io/part-of: VapusData
app.kubernetes.io/component: nabhikserver
{{- end }}

{{- define "vapusdata.nabhikserver.selectorLabels" -}}
app.kubernetes.io/component: nabhikserver
{{- end }}

{{- define "vapusdata.webapp.labels" -}}
app.kubernetes.io/part-of: VapusData
app.kubernetes.io/component: webapp
{{- end }}

{{- define "vapusdata.webapp.selectorLabels" -}}
app.kubernetes.io/component: webapps
{{- end }}


{{- define "vapusdata.common.labels" -}}
app.kubernetes.io/part-of: VapusData
app.kubernetes.io/component: common
{{- end }}

{{- define "vapusdata.common.selectorLabels" -}}
app.kubernetes.io/component: common
{{- end }}

{{- define "vapusdata.aistudio.labels" -}}
app.kubernetes.io/part-of: VapusData
app.kubernetes.io/component: aistudio
{{- end }}

{{/*
Selector labels
*/}}
{{- define "vapusdata.aistudio.selectorLabels" -}}
app.kubernetes.io/component: aistudio
{{- end }}

{{- define "vapusdata.nabrunners.labels" -}}
app.kubernetes.io/part-of: VapusData
app.kubernetes.io/component: nabrunners
{{- end }}

{{- define "vapusdata.nabrunners.selectorLabels" -}}
app.kubernetes.io/component: nabrunners
{{- end }}

{{- define "vapusdata.vapusaigateway.labels" -}}
app.kubernetes.io/part-of: VapusData
app.kubernetes.io/component: vapusaigateway
{{- end }}

{{- define "vapusdata.vapusaigateway.selectorLabels" -}}
app.kubernetes.io/component: vapusaigateway
{{- end }}