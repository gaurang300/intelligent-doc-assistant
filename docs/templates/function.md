# {{ .Name }}

{{ .Description }}

## Usage

```{{ .Language }}
{{ .Example }}
```

## Parameters

{{ range .Parameters }}
- **{{ .Name }}** ({{ .Type }}): {{ .Description }}
{{ end }}

## Returns

{{ .Returns }}

## Source Location

File: `{{ .FilePath }}`
Lines: {{ .StartLine }} - {{ .EndLine }}
