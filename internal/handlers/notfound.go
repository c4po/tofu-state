package handlers

import (
	"html/template"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Tofu State - Page Not Found</title>
		<style>
			body { 
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
				text-align: center;
				padding: 2rem;
				background-color: #f8f9fa;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
			}
			h1 { color: #2c3e50; }
			p { color: #7f8c8d; }
			a { color: #3498db; text-decoration: none; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>404 - Page Not Found</h1>
			<p>The requested URL <code>{{.Path}}</code> was not found on this server.</p>
			<p>Return to <a href="/">Home Page</a></p>
		</div>
	</body>
	</html>
	`

	tpl := template.Must(template.New("404").Parse(html))
	tpl.Execute(w, map[string]interface{}{
		"Path": r.URL.Path,
	})
}
