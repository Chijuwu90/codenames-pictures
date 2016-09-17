package codenames

import (
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
)

const tpl = `
<!DOCTYPE html>
<html>
    <head>
        <title>Codenames</title>
        <script src="/js/lib/browser.min.js"></script>
        <script src="/js/lib/react.min.js"></script>
        <script src="/js/lib/react-dom.min.js"></script>
        <script src="/js/lib/jquery-3.0.0.min.js"></script>

        <script type="text/babel">
             {{if .SelectedGameID}}
             window.selectedGameID = "{{.SelectedGameID}}";
             {{end}}
             window.autogeneratedGameID = "{{.AutogeneratedGameID}}";
        </script>

        {{range .JSScripts}}
            <script type="text/babel" src="/js/{{ . }}"></script>
        {{end}}
        {{range .Stylesheets}}
            <link rel="stylesheet" type="text/css" href="/css/{{ . }}" />
        {{end}}
        <style type="text/css">
			html, body {
			  margin: 0;
			  font-family: 'Roboto', verdana, sans-serif;
			}

			#application { margin: 1em; }

			#topbar {
			  padding: 1em;
			  margin-bottom: 1em;
			}

			h1 {
			  text-transform: uppercase;
			  letter-spacing: 0.2em;
			  text-align: center;
			  font-family: "Courier New", monospace;
			}

			h1, h2, h3 {
				font-family: 'Courier New', monospace;
			}

			#game-view, .loading {
				width: 700px;
				margin: 0 auto;
			}
        </style>

        <link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
    </head>
    <body>
		<div id="app">
			<div id="application">
				<div id="topbar">
					<h1>Codenames</h1>
				</div>
			</div>
		</div>
        <script type="text/babel">
            ReactDOM.render(<window.App />, document.getElementById('app'));
        </script>
    </body>
</html>
`

type templateParameters struct {
	SelectedGameID      string
	AutogeneratedGameID string
	JSLibs              []string
	JSScripts           []string
	Stylesheets         []string
}

func (s *Server) handleIndex(rw http.ResponseWriter, req *http.Request) {
	dir, id := filepath.Split(req.URL.Path)
	if dir != "" && dir != "/" {
		http.NotFound(rw, req)
		return
	}

	first := s.words[rand.Intn(len(s.words))]
	second := s.words[rand.Intn(len(s.words))]
	autogeneratedID := strings.Replace(first+"-"+second, " ", "", -1)
	autogeneratedID = strings.ToLower(autogeneratedID)

	err := s.tpl.Execute(rw, templateParameters{
		SelectedGameID:      id,
		AutogeneratedGameID: autogeneratedID,
		JSLibs:              s.jslib.RelativePaths(),
		JSScripts:           s.js.RelativePaths(),
		Stylesheets:         s.css.RelativePaths(),
	})
	if err != nil {
		http.Error(rw, "error rendering", http.StatusInternalServerError)
	}
}
