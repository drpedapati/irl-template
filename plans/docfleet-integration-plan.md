# Docfleet Integration Plan

## Goal

Move the IRL explainer site (`draft.html`) from this repo into the **docfleet** deployment system so it can be served as a production static site with HTTPS and a real hostname.

---

## Current State

### IRL Template Repo (`irl-template`)
- `03-outputs/irl-explainer-site/draft.html` — the combined page (2,972 lines, 134 KB), self-contained HTML with inline CSS and JS, 13 sections, 8 interactive toys
- `03-outputs/irl-explainer-site/toys.html` — standalone toy playground (now redundant since toys are integrated into draft.html)
- `03-outputs/irl-explainer-site/index.html` — original explainer page (superseded by draft.html)
- `03-outputs/irl-explainer-site/draft-outline-backup.html` — outline backup (development artifact)

### Docfleet Repo (`docfleet`)
- Docker container running nginx:alpine
- Kamal 2.9 deployment to a VPS with auto-SSL via Let's Encrypt
- Hostname-based routing: each site gets a `server {}` block in `nginx.conf`
- Sites live in `sites/<site-name>/` as static files
- Currently has two example sites: `example-memory` (howmemoryworks.com) and `example-placeholder`

---

## Integration Steps

### Step 1: Choose a hostname

Pick a domain/subdomain for the IRL explainer. Options:
- A subdomain like `irl.yourdomain.com`
- A new domain
- Reuse an existing domain

**Decision needed from author.**

### Step 2: Prepare the site files

The IRL explainer is a single self-contained HTML file. For docfleet:

```
docfleet/sites/irl-explainer/
└── index.html          ← renamed from draft.html
```

No build step needed — it's pure static HTML with inline CSS/JS. Just copy and rename.

Optional: also include `toys.html` as a standalone page if you want to keep it accessible:
```
docfleet/sites/irl-explainer/
├── index.html          ← draft.html (main page)
└── toys.html           ← standalone toy playground
```

### Step 3: Add nginx server block

Add a new `server {}` block to `docfleet/nginx.conf`:

```nginx
# IRL Explainer Site
server {
    listen 80;
    server_name HOSTNAME;    # ← replace with chosen hostname
    root /sites/irl-explainer;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### Step 4: Add hostname to Kamal config

In `docfleet/config/deploy.yml`, add the hostname to the proxy hosts list:

```yaml
proxy:
  hosts:
    - howmemoryworks.com
    - docs.example.com
    - HOSTNAME              # ← add the new hostname
  ssl: true
```

### Step 5: Point DNS

Create an A record pointing the chosen hostname to the docfleet server IP.

### Step 6: Deploy

```bash
cd docfleet
kamal deploy
```

Kamal builds the Docker image (which copies `sites/irl-explainer/` into the container), pushes it to the registry, and deploys with zero downtime. Let's Encrypt provisions the SSL certificate automatically.

### Step 7: Update links

- Update the "View Source" link in the topbar if needed
- Update any internal references (the brand link `href="#"` could point to the live URL)

---

## Ongoing Workflow

After initial deployment, the update cycle is:

1. Edit `draft.html` in this IRL template repo (where the plan-driven workflow lives)
2. When ready to publish, copy the updated file to docfleet:
   ```bash
   cp 03-outputs/irl-explainer-site/draft.html ~/Documents/GitHub/docfleet/sites/irl-explainer/index.html
   ```
3. Commit and deploy:
   ```bash
   cd ~/Documents/GitHub/docfleet
   git add sites/irl-explainer/index.html
   git commit -m "Update IRL explainer site"
   kamal deploy
   ```

Alternatively, a script or Makefile target could automate steps 2-3.

---

## Files to Copy

| Source (irl-template) | Destination (docfleet) | Notes |
|---|---|---|
| `03-outputs/irl-explainer-site/draft.html` | `sites/irl-explainer/index.html` | Main page, rename to index.html |
| `03-outputs/irl-explainer-site/toys.html` | `sites/irl-explainer/toys.html` | Optional, standalone playground |

Files NOT to copy (development artifacts):
- `index.html` — superseded by draft.html
- `draft-outline-backup.html` — outline backup
- `Makefile` — local dev only
- `README.md` — local dev only

---

## Decisions Needed

1. **Hostname** — what domain/subdomain should serve the IRL explainer?
2. **Include toys.html?** — keep the standalone toy playground as a second page, or just publish the combined page?
3. **Server IP and registry** — the deploy.yml needs real values (currently placeholder `YOUR_SERVER_IP`)
4. **Automation** — want a Makefile target or script to automate the copy-and-deploy cycle?
