module github.com/ancalabrese/gotth/example

go 1.24.1

require github.com/ancalabrese/gotth v0.0.0-20250525102643-f20b3cf8622c

require (
	github.com/a-h/templ v0.3.865
)

replace (
	github.com/ancalabrese/gotth => ../.
)
