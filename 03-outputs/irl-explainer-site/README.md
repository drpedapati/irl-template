# IRL Explainer Site (Draft)

Open `index.html` in a browser.

Handy commands:
- `make serve` (then open `http://localhost:8765`)
- `make smoke` (starts a server briefly and verifies `index.html` loads)
- `make ips` (prints likely LAN IPs)
- `make zt` (shows ZeroTier info, if installed)

Notes:
- The page renders diagrams client-side via the `beautiful-mermaid` CDN, so it needs an internet connection.
- If you want to serve it locally: `python3 -m http.server` from this folder, then open `http://localhost:8000`.

## Access Over ZeroTier (Remote Testing)

1. Start the server: `make serve PORT=8765`
2. Find your ZeroTier interface IP.
   - If you have `zerotier-cli`: `make zt`
   - Otherwise: `ifconfig` and look for a `zt*` interface
3. Open from another device on the same ZeroTier network:
   - `http://<your-zerotier-ip>:8765/`
