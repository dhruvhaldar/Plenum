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
      max-width: 1080px;
      margin: 36px auto;
      padding: 0 20px 48px;
      line-height: 1.45;
      background: #fafafa;
      color: #111;
    }
    h1, h2, h3 {
      margin-bottom: 8px;
    }
    code {
      background: #f0f0f0;
      padding: 2px 6px;
      border-radius: 4px;
    }
    .layout {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
      gap: 16px;
    }
    .card {
      border: 1px solid #ddd;
      border-radius: 12px;
      padding: 16px;
      background: white;
    }
    .card p {
      margin-top: 0;
      color: #333;
      font-size: 0.95rem;
    }
    .row {
      display: flex;
      flex-direction: column;
      gap: 6px;
      margin-bottom: 8px;
    }
    label {
      font-size: 0.9rem;
      color: #222;
    }
    input, textarea, select {
      border: 1px solid #bbb;
      border-radius: 8px;
      padding: 8px;
      font: inherit;
    }
    textarea {
      min-height: 72px;
      resize: vertical;
    }
    button {
      margin-top: 8px;
      padding: 10px 14px;
      border-radius: 8px;
      border: 1px solid #222;
      background: #fff;
      cursor: pointer;
    }
    button:hover {
      background: #f7f7f7;
    }
    #output {
      margin-top: 18px;
      background: #111;
      color: #eee;
      padding: 16px;
      border-radius: 10px;
      overflow: auto;
      min-height: 180px;
      white-space: pre-wrap;
    }
  </style>
</head>
<body>
  <h1>Plenum</h1>
  <p>AeroGo preprocessing and telemetry API dashboard.</p>

  <div class="layout">
    <div class="card">
      <h2>Health</h2>
      <p><code>GET /api/health</code></p>
      <button onclick="checkHealth()">Check API health</button>
    </div>

    <div class="card">
      <h2>Arrhenius sweep</h2>
      <p><code>GET /api/chem/arrhenius</code></p>
      <div class="row"><label for="arrA">A</label><input id="arrA" type="number" value="10000000" /></div>
      <div class="row"><label for="arrEa">Ea</label><input id="arrEa" type="number" value="50000" /></div>
      <div class="row"><label for="arrTmin">Tmin</label><input id="arrTmin" type="number" value="300" /></div>
      <div class="row"><label for="arrTmax">Tmax</label><input id="arrTmax" type="number" value="1200" /></div>
      <div class="row"><label for="arrPoints">Points</label><input id="arrPoints" type="number" value="5" min="2" /></div>
      <button onclick="runArrhenius()">Run calculation</button>
    </div>

    <div class="card">
      <h2>OpenFOAM dictionary generation</h2>
      <p><code>POST /api/foam/generate-dict</code></p>
      <div class="row"><label for="foamSolver">Solver</label>
        <select id="foamSolver">
          <option>simpleFoam</option>
          <option>pisoFoam</option>
          <option>rhoCentralFoam</option>
        </select>
      </div>
      <div class="row"><label for="foamDeltaT">deltaT</label><input id="foamDeltaT" type="number" value="0.001" step="0.0001" /></div>
      <div class="row"><label for="foamEndTime">endTime</label><input id="foamEndTime" type="number" value="10" /></div>
      <div class="row"><label for="foamLength">length</label><input id="foamLength" type="number" value="12" /></div>
      <div class="row"><label for="foamWidth">width</label><input id="foamWidth" type="number" value="4" /></div>
      <div class="row"><label for="foamHeight">height</label><input id="foamHeight" type="number" value="3" /></div>
      <div class="row"><label for="foamX">xCells</label><input id="foamX" type="number" value="120" /></div>
      <div class="row"><label for="foamY">yCells</label><input id="foamY" type="number" value="40" /></div>
      <div class="row"><label for="foamZ">zCells</label><input id="foamZ" type="number" value="30" /></div>
      <div class="row"><label for="foamGrading">simpleGrading</label><input id="foamGrading" value="1 1 1" /></div>
      <button onclick="generateFoamDict()">Generate dictionaries</button>
    </div>

    <div class="card">
      <h2>DrivAer bounds</h2>
      <p><code>POST /api/auto/drivaer-bounds</code></p>
      <div class="row"><label for="vehicleLength">vehicleLength</label><input id="vehicleLength" type="number" value="4.6" step="0.1" /></div>
      <div class="row"><label for="blockageRatio">blockageRatio</label><input id="blockageRatio" type="number" value="0.03" step="0.005" /></div>
      <button onclick="runDrivAerBounds()">Calculate bounds</button>
    </div>

    <div class="card">
      <h2>Telemetry update</h2>
      <p><code>POST /api/telemetry/update</code></p>
      <div class="row"><label for="teleJobId">jobId</label><input id="teleJobId" value="job-001" /></div>
      <div class="row"><label for="teleIteration">iteration</label><input id="teleIteration" type="number" value="1" min="0" /></div>
      <div class="row"><label for="teleCourant">courant</label><input id="teleCourant" type="number" value="0.8" step="0.01" /></div>
      <div class="row"><label for="teleResiduals">residuals (JSON map)</label><textarea id="teleResiduals">{"Ux":0.01,"Uy":0.02,"p":0.03}</textarea></div>
      <button onclick="postTelemetry()">Send telemetry</button>
    </div>

    <div class="card">
      <h2>Telemetry query</h2>
      <p><code>GET /api/telemetry/update?jobId=...</code></p>
      <div class="row"><label for="teleGetJobId">jobId</label><input id="teleGetJobId" value="job-001" /></div>
      <button onclick="getTelemetry()">Get telemetry</button>
    </div>
  </div>

  <h2>Output</h2>
  <pre id="output">Waiting...</pre>

  <script>
    function show(data) {
      document.getElementById('output').textContent = typeof data === 'string' ? data : JSON.stringify(data, null, 2);
    }

    async function request(url, options) {
      const res = await fetch(url, options);
      const text = await res.text();
      let body = text;
      try { body = JSON.parse(text); } catch (_) {}
      if (!res.ok) {
        throw new Error('HTTP ' + res.status + ': ' + (typeof body === 'string' ? body : JSON.stringify(body)));
      }
      return body;
    }

    async function checkHealth() {
      try {
        show(await request('/api/health'));
      } catch (err) {
        show(err.message || String(err));
      }
    }

    async function runArrhenius() {
      try {
        const params = new URLSearchParams({
          A: document.getElementById('arrA').value,
          Ea: document.getElementById('arrEa').value,
          Tmin: document.getElementById('arrTmin').value,
          Tmax: document.getElementById('arrTmax').value,
          points: document.getElementById('arrPoints').value,
        });
        show(await request('/api/chem/arrhenius?' + params.toString()));
      } catch (err) {
        show(err.message || String(err));
      }
    }

    async function generateFoamDict() {
      try {
        const payload = {
          solver: document.getElementById('foamSolver').value,
          deltaT: Number(document.getElementById('foamDeltaT').value),
          endTime: Number(document.getElementById('foamEndTime').value),
          length: Number(document.getElementById('foamLength').value),
          width: Number(document.getElementById('foamWidth').value),
          height: Number(document.getElementById('foamHeight').value),
          xCells: Number(document.getElementById('foamX').value),
          yCells: Number(document.getElementById('foamY').value),
          zCells: Number(document.getElementById('foamZ').value),
          simpleGrading: document.getElementById('foamGrading').value,
        };
        show(await request('/api/foam/generate-dict', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(payload),
        }));
      } catch (err) {
        show(err.message || String(err));
      }
    }

    async function runDrivAerBounds() {
      try {
        const payload = {
          vehicleLength: Number(document.getElementById('vehicleLength').value),
          blockageRatio: Number(document.getElementById('blockageRatio').value),
        };
        show(await request('/api/auto/drivaer-bounds', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(payload),
        }));
      } catch (err) {
        show(err.message || String(err));
      }
    }

    async function postTelemetry() {
      try {
        const residuals = JSON.parse(document.getElementById('teleResiduals').value || '{}');
        const payload = {
          jobId: document.getElementById('teleJobId').value,
          iteration: Number(document.getElementById('teleIteration').value),
          courant: Number(document.getElementById('teleCourant').value),
          residuals,
        };
        show(await request('/api/telemetry/update', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(payload),
        }));
      } catch (err) {
        show(err.message || String(err));
      }
    }

    async function getTelemetry() {
      try {
        const jobId = encodeURIComponent(document.getElementById('teleGetJobId').value);
        show(await request('/api/telemetry/update?jobId=' + jobId));
      } catch (err) {
        show(err.message || String(err));
      }
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
