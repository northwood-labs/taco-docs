/**
 * DAG Diagram Renderer (generic, data-driven)
 *
 * Reads graph-data.json, validates it, computes a sugiyama layout via d3-dag,
 * and renders an interactive SVG with d3. The data file is the only thing
 * that changes between projects — this renderer is reusable as-is.
 */

import * as d3 from "d3";
import { coordCenter, decrossTwoLayer, graphStratify, sugiyama } from "d3-dag";
import { validate, type GraphDocument, type LayerDefinition, type NodeData } from "./schema";

// ─── Load Data ─────────────────────────────────────────────────────────────────

const response = await fetch("graph-data.json");
const doc: GraphDocument = await response.json();

const errors = validate(doc);
if (errors.length > 0) {
  document.body.innerHTML = `<pre style="color:#f85149;padding:2em;">Validation errors:\n${errors.join("\n")}</pre>`;
  throw new Error(`Graph data invalid: ${errors.join("; ")}`);
}

// Build lookup maps
const layerMap = new Map<string, LayerDefinition>(doc.layers.map((l) => [l.id, l]));

// ─── Build Dynamic Header ──────────────────────────────────────────────────────

const titleEl = document.querySelector("header h1");
if (titleEl) titleEl.textContent = doc.title;

const subtitleEl = document.querySelector("header .subtitle");
if (subtitleEl) subtitleEl.textContent = doc.subtitle;

const legendEl = document.getElementById("legend");
if (legendEl) {
  legendEl.innerHTML = "";
  for (const layer of doc.layers) {
    const item = document.createElement("div");
    item.className = "legend-item";
    item.innerHTML = `<div class="legend-swatch" style="background:${layer.stroke}"></div><span>${layer.label}</span>`;
    legendEl.appendChild(item);
  }
}

// ─── Configuration ─────────────────────────────────────────────────────────────

const NODE_WIDTH = 180;
const NODE_HEIGHT = 52;
const NODE_SEP_X = 40;
const NODE_SEP_Y = 70;

// ─── Build DAG ─────────────────────────────────────────────────────────────────

const stratify = graphStratify();
const dag = stratify(doc.nodes);

// ─── Layout ────────────────────────────────────────────────────────────────────

const layout = sugiyama()
  .decross(decrossTwoLayer())
  .coord(coordCenter())
  .nodeSize((node) => {
    if (node === undefined) {
      return [NODE_SEP_X * 0.5, NODE_SEP_Y * 0.3];
    }
    return [NODE_WIDTH + NODE_SEP_X, NODE_HEIGHT + NODE_SEP_Y];
  });

const { width, height } = layout(dag);

// ─── SVG Setup ─────────────────────────────────────────────────────────────────

const container = d3.select("#graph");
const svg = container.append("svg").attr("viewBox", `0 0 ${width + 100} ${height + 100}`);

const g = svg.append("g").attr("transform", "translate(50, 50)");

const zoom = d3
  .zoom<SVGSVGElement, unknown>()
  .scaleExtent([0.3, 3])
  .on("zoom", (event) => {
    g.attr("transform", event.transform);
  });

svg.call(zoom);

// Initial transform: fit and center
const containerEl = document.getElementById("graph")!;
const initialScale = Math.min(containerEl.clientWidth / (width + 100), containerEl.clientHeight / (height + 100), 1.0);
const initialX = (containerEl.clientWidth - (width + 100) * initialScale) / 2;
const initialY = 20;

svg.call(zoom.transform, d3.zoomIdentity.translate(initialX, initialY).scale(initialScale));

// ─── Arrow Marker ──────────────────────────────────────────────────────────────

svg.append("defs").append("marker").attr("id", "arrowhead").attr("viewBox", "0 0 10 10").attr("refX", 10).attr("refY", 5).attr("markerWidth", 8).attr("markerHeight", 8).attr("orient", "auto-start-reverse").append("path").attr("d", "M 0 0 L 10 5 L 0 10 Z").attr("fill", "#8b949e");

// ─── Render Edges ──────────────────────────────────────────────────────────────

const line = d3
  .line<[number, number]>()
  .x((d) => d[0])
  .y((d) => d[1])
  .curve(d3.curveBasis);

const edgeGroup = g.append("g").attr("class", "edges");

for (const link of dag.links()) {
  // d3-dag points are [x, y] tuples
  const rawPoints: [number, number][] = link.points.map((p: any) => [p[0] as number, p[1] as number]);

  // Offset endpoints so edges exit/enter node boundaries, not centers
  const points: [number, number][] = rawPoints.map((p, i) => {
    if (i === 0) return [p[0], p[1] + NODE_HEIGHT / 2 + 2];
    if (i === rawPoints.length - 1) return [p[0], p[1] - NODE_HEIGHT / 2 - 6];
    return p;
  });

  edgeGroup.append("path").attr("class", "edge").attr("d", line(points)!).attr("fill", "none").attr("stroke", "#8b949e").attr("stroke-width", 1.5).attr("marker-end", "url(#arrowhead)").attr("data-source", link.source.data.id).attr("data-target", link.target.data.id);
}

// ─── Render Nodes ──────────────────────────────────────────────────────────────

const nodeGroup = g.append("g").attr("class", "nodes");

for (const node of dag.nodes()) {
  const d = node.data as NodeData;
  const layer = layerMap.get(d.layer)!;

  const nodeG = nodeGroup
    .append("g")
    .attr("class", `node node-${d.layer}`)
    .attr("transform", `translate(${node.x - NODE_WIDTH / 2}, ${node.y - NODE_HEIGHT / 2})`)
    .attr("data-id", d.id);

  // Background rect
  nodeG.append("rect").attr("width", NODE_WIDTH).attr("height", NODE_HEIGHT).attr("rx", 6).attr("ry", 6).attr("fill", layer.fill).attr("stroke", layer.stroke).attr("stroke-width", 1.5);

  // Layer indicator bar
  nodeG
    .append("rect")
    .attr("x", 0)
    .attr("y", 6)
    .attr("width", 3)
    .attr("height", NODE_HEIGHT - 12)
    .attr("rx", 1.5)
    .attr("fill", layer.stroke)
    .attr("opacity", 0.8);

  // Primary label
  nodeG.append("text").attr("x", 14).attr("y", 22).attr("class", "label").attr("font-size", "11px").attr("font-weight", "600").attr("fill", "#f0f6fc").attr("font-family", "SF Mono, Menlo, Monaco, monospace").text(d.label);

  // Sublabel
  nodeG.append("text").attr("x", 14).attr("y", 38).attr("class", "sublabel").attr("font-size", "9px").attr("fill", "#8b949e").attr("font-family", "-apple-system, BlinkMacSystemFont, sans-serif").text(d.sublabel);

  // Hover interactions
  nodeG.on("mouseenter", () => highlightConnections(d.id, true));
  nodeG.on("mouseleave", () => highlightConnections(d.id, false));
}

// ─── Interaction Helpers ───────────────────────────────────────────────────────

function highlightConnections(nodeId: string, active: boolean): void {
  const highlightColor = "#58a6ff";
  const defaultColor = "#8b949e";

  edgeGroup.selectAll("path").each(function () {
    const el = d3.select(this);
    const source = el.attr("data-source");
    const target = el.attr("data-target");

    if (source === nodeId || target === nodeId) {
      el.attr("stroke", active ? highlightColor : defaultColor).attr("stroke-width", active ? 2.5 : 1.5);
    } else {
      el.attr("opacity", active ? 0.2 : 1);
    }
  });

  const connectedIds = new Set<string>([nodeId]);
  for (const link of dag.links()) {
    if (link.source.data.id === nodeId) connectedIds.add(link.target.data.id);
    if (link.target.data.id === nodeId) connectedIds.add(link.source.data.id);
  }

  nodeGroup.selectAll(".node").each(function () {
    const el = d3.select(this);
    const id = el.attr("data-id");
    el.attr("opacity", active ? (connectedIds.has(id) ? 1 : 0.3) : 1);
  });
}
