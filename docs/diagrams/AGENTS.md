# AGENTS.md

Guidance for AI agents working in this diagram sub-project.

## Project overview

This is a **reusable** DAG visualization tool. A generator produces `graph-data.json` — the renderer handles everything else (layout, drawing, interaction) without modification.

## Architecture

```text
Generator (any language)        Renderer (this project)
─────────────────────────       ───────────────────────
Analyze codebase           →    graph-data.json
                                     ↓
                                schema.ts (validate)
                                     ↓
                                render.ts (layout + draw)
                                     ↓
                                index.html + styles.css
```

The only file a generator needs to produce is `graph-data.json`. The renderer reads it at runtime via `fetch()`, validates it against the schema, and renders the diagram.

## File structure

```text
docs/diagrams/
  graph-data.json  → Generator output (the only file that changes per project)
  schema.ts        → TypeScript types + runtime validation (the contract)
  render.ts        → Layout (d3-dag) + SVG rendering (d3) + interactions
  index.html       → Minimal shell (title/legend populated from JSON at runtime)
  styles.css       → Visual styling (dark theme)
  serve.ts         → Bun dev server (bundles TS, serves static files)
  package.json     → Dependencies
  tsconfig.json    → TypeScript config
```

## Commands

```bash
bun install        # install dependencies
bun run dev        # start dev server at http://localhost:3333
bun run build      # bundle for static hosting → dist/render.js
```

## The generator contract (`graph-data.json`)

A generator must produce a JSON file matching this schema:

```json
{
  "title": "Project Name Architecture",
  "subtitle": "Description of what the diagram shows",
  "layers": [
    { "id": "layer-id", "label": "Display Name", "fill": "#1f3a5f", "stroke": "#388bfd" }
  ],
  "nodes": [
    {
      "id": "unique-node-id",
      "parentIds": ["id-of-caller"],
      "label": "Component Name",
      "sublabel": "Brief description",
      "layer": "layer-id"
    }
  ]
}
```

### Rules

* `title` — displayed in the page header
* `subtitle` — displayed below the title
* `layers` — define your own categories (any number, any names)
  * `id` — referenced by nodes
  * `label` — shown in the legend
  * `fill` — dark background color for nodes in this layer
  * `stroke` — bright accent color for borders and indicators
* `nodes` — the graph itself
  * `id` — unique identifier, referenced in other nodes' `parentIds`
  * `parentIds` — IDs of nodes that _call_ this node (edges flow parent → child)
  * Empty `parentIds` marks a root node (entry point)
  * `layer` — must match a layer `id`
* All `parentIds` must reference existing node IDs (validated at runtime)
* No cycles allowed (this is a DAG — d3-dag will throw if cycles exist)

## Critical d3-dag gotchas

### Points are `[x, y]` tuples

```typescript
// WRONG — returns undefined, edges invisible
link.points.map((p) => ({ x: p.x, y: p.y }))

// CORRECT
link.points.map((p: any) => [p[0] as number, p[1] as number])
```

### nodeSize callback receives undefined for dummy nodes

```typescript
.nodeSize((node) => {
  if (node === undefined) return [20, 21];  // routing dummies
  return [NODE_WIDTH + SEP_X, NODE_HEIGHT + SEP_Y];
})
```

### Edges clip under opaque nodes

Offset edge endpoints away from node centers:

```typescript
if (i === 0) return [p[0], p[1] + NODE_HEIGHT / 2 + 2];      // exit bottom
if (i === last) return [p[0], p[1] - NODE_HEIGHT / 2 - 6];   // enter top
```

### graphStratify input

Expects objects with `id: string` and `parentIds: string[]`. Other fields are preserved as `node.data`.

## Bun-specific notes

* Bun binary is at `~/.bun/bin/bun` on this system (not in default PATH)
* The dev server (`serve.ts`) bundles `.ts` files on-the-fly via `Bun.build()`
* `<script src="render.ts">` works because the server intercepts and returns bundled JS
* Top-level `await` works in Bun (used in render.ts for `fetch()`)

## Writing a generator

A generator can be written in any language. It should:

1. Analyze the target codebase (parse AST, read imports, trace call graphs)
2. Identify components and their calling relationships
3. Assign each component to a layer (define layers that make sense for the domain)
4. Output `graph-data.json` conforming to the schema above

Example generator pseudocode:

```python
import json

doc = {
    "title": "My Project Architecture",
    "subtitle": "Module dependencies",
    "layers": [
        {"id": "api", "label": "API Layer", "fill": "#1f3a5f", "stroke": "#388bfd"},
        {"id": "db", "label": "Database", "fill": "#3d1f1f", "stroke": "#f78166"},
    ],
    "nodes": []
}

for module in analyze_project():
    doc["nodes"].append({
        "id": module.name,
        "parentIds": [dep.name for dep in module.callers],
        "label": module.short_name,
        "sublabel": module.description,
        "layer": classify_layer(module),
    })

with open("graph-data.json", "w") as f:
    json.dump(doc, f, indent=2)
```

## Rendering behavior

* Layout: sugiyama (layered, top-to-bottom flow)
* Crossing minimization: two-layer heuristic
* Coordinate assignment: center-aligned
* Zoom/pan: d3.zoom with scroll + drag
* Hover: highlights connected edges, dims unrelated nodes
* Header and legend: populated dynamically from JSON metadata
