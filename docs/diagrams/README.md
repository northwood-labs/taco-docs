# Architecture Diagram

Reusable DAG (directed acyclic graph) visualization tool. A generator produces `graph-data.json` — the renderer handles layout, drawing, and interaction without modification.

## How it works

```text
Generator (any language)     →   graph-data.json   →   Renderer (this project)
```

The renderer:

1. Fetches `graph-data.json` at runtime
2. Validates the data against the schema
3. Computes a sugiyama layout (layered, top-to-bottom)
4. Renders an interactive SVG with zoom, pan, and hover highlighting

## What the diagram shows

* **Nodes** represent components (packages, types, functions)
* **Edges** show directed calling relationships (arrow = caller → callee)
* **Colors** indicate layers defined in the data file

## Running locally

```bash
bun install
bun run dev
# Open http://localhost:3333
```

## Building for static hosting

```bash
bun run build
# Output: dist/render.js — serve alongside index.html, styles.css, and graph-data.json
```

## Producing `graph-data.json`

A generator can be written in any language. It outputs a JSON file with this structure:

```json
{
  "title": "Project Architecture",
  "subtitle": "Component calling relationships",
  "layers": [
    { "id": "layer-id", "label": "Display Name", "fill": "#1f3a5f", "stroke": "#388bfd" }
  ],
  "nodes": [
    {
      "id": "unique-id",
      "parentIds": ["caller-id"],
      "label": "Component Name",
      "sublabel": "Brief description",
      "layer": "layer-id"
    }
  ]
}
```

Rules:

* `parentIds` lists nodes that _call into_ this node (edges flow parent → child)
* Empty `parentIds` marks a root/entry point
* All `parentIds` must reference existing node IDs
* No cycles (this is a DAG)
* Layer IDs in nodes must match a defined layer

See `AGENTS.md` for the full specification and implementation details.

## Interaction

* **Pan** — click and drag
* **Zoom** — scroll wheel or pinch
* **Hover** — highlights connected edges, dims unrelated nodes
