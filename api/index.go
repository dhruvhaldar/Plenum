package handler

import (
	"net/http"
)

const page = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Plenum</title>
  <style>
    body {
      font-family: system-ui, sans-serif;
      max-width: 900px;
      margin: 48px auto;
      padding: 0 20px;
      line-height: 1.5;
    }
    code {
      background: #f2f2f2;
      padding: 2px 6px;
      border-radius: 4px;
    }
    .card {
      border: 1px solid #ddd;
      border-radius: 12px;
      padding: 20px;
      margin: 16px 0;
    }
    button {
      padding: 10px 14px;
      border-radius: 8px;
      border: 1px solid #222;
      cursor: pointer;
    }
    pre {
      background: #111;
      color: #eee;
      padding: 16px;
      border-radius: 10px;
      overflow: auto;
    }
  </style>
</head>
<body>
  <h1>Plenum</h1>
  <p>AeroGo preprocessing and telemetry API dashboard.</p>

  <div class="card">
    <h2>Health</h2>
    <button onclick="checkHealth()">Check API health</button>
  </div>

  <div class="card">
    <h2>Arrhenius sweep</h2>
    <button onclick="runArrhenius()">Run sample calculation</button>
  </div>

  <h2>Output</h2>
  <pre id="output">Waiting...</pre>

  <script>
    async function checkHealth() {
      const output = document.getElementById('output');
      try {
        const res = await fetch('/api/health');
        if (!res.ok) {
          throw new Error('Health check failed with status ' + res.status);
        }

        const data = await res.json();
        output.textContent = JSON.stringify(data, null, 2);
      } catch (err) {
        output.textContent = err.message || String(err);
      }
    }

    async function runArrhenius() {
      const url = '/api/chem/arrhenius?A=10000000&Ea=50000&Tmin=300&Tmax=1200&points=5';
      const res = await fetch(url);
      const data = await res.json();
      document.getElementById('output').textContent = JSON.stringify(data, null, 2);
    }
  </script>
</body>
</html>`

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}
