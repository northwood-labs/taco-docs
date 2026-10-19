/**
 * DAG Diagram Schema
 *
 * This is the contract between generators and the renderer. A generator
 * produces a JSON file conforming to GraphDocument. The renderer consumes
 * it without modification.
 */

/**
 * A single layer (category) that nodes can belong to.
 * Layers control visual grouping and color-coding.
 */
export interface LayerDefinition {
  /** Unique identifier for this layer. */
  id: string;
  /** Human-readable display name (shown in legend). */
  label: string;
  /** Dark fill color for node backgrounds. */
  fill: string;
  /** Brighter accent color for node borders and indicators. */
  stroke: string;
}

/**
 * A node in the directed acyclic graph.
 */
export interface NodeData {
  /** Unique identifier. Used in parentIds references. */
  id: string;
  /** IDs of nodes that call/invoke this node (edges flow parent → child). */
  parentIds: string[];
  /** Primary display text (component name). */
  label: string;
  /** Secondary description text. */
  sublabel: string;
  /** Layer ID this node belongs to. Must match a LayerDefinition.id. */
  layer: string;
}

/**
 * Top-level document structure produced by generators.
 * This is the shape of graph-data.json.
 */
export interface GraphDocument {
  /** Diagram title displayed in the header. */
  title: string;
  /** Subtitle / description shown below the title. */
  subtitle: string;
  /** Layer definitions — order determines legend display order. */
  layers: LayerDefinition[];
  /** All nodes in the graph. */
  nodes: NodeData[];
}

/**
 * Validate a GraphDocument at runtime. Returns an array of error messages
 * (empty if valid).
 */
export function validate(doc: GraphDocument): string[] {
  const errors: string[] = [];

  if (!doc.title) errors.push("Missing 'title' field");
  if (!doc.layers || doc.layers.length === 0) errors.push("No layers defined");
  if (!doc.nodes || doc.nodes.length === 0) errors.push("No nodes defined");

  const layerIds = new Set(doc.layers?.map((l) => l.id) ?? []);
  const nodeIds = new Set(doc.nodes?.map((n) => n.id) ?? []);

  // Check for duplicate node IDs
  const seen = new Set<string>();
  for (const node of doc.nodes ?? []) {
    if (seen.has(node.id)) {
      errors.push(`Duplicate node ID: "${node.id}"`);
    }
    seen.add(node.id);
  }

  for (const node of doc.nodes ?? []) {
    // Validate layer reference
    if (!layerIds.has(node.layer)) {
      errors.push(`Node "${node.id}" references unknown layer "${node.layer}"`);
    }

    // Validate parent references
    for (const parentId of node.parentIds) {
      if (!nodeIds.has(parentId)) {
        errors.push(`Node "${node.id}" references non-existent parent "${parentId}"`);
      }
    }
  }

  return errors;
}
