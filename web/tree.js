"use strict";

function layoutTree(nodes) {
  const positions = {};
  if (nodes.length === 0) return positions;
  const ordered = [...nodes].sort((a, b) => a.id < b.id ? -1 : a.id > b.id ? 1 : 0);
  const byID = new Map(ordered.map((n) => [n.id, n]));
  const roots = ordered.filter((n) => n.after.length === 0);
  if (byID.size !== nodes.length || roots.length !== 1) throw new Error("Tree needs unique IDs and one root");
  const root = roots[0];
  const depths = new Map([[root.id, 0]]);
  const visiting = new Set();
  function depthOf(id) {
    if (depths.has(id)) return depths.get(id);
    if (!byID.has(id)) throw new Error("Unknown prerequisite: " + id);
    if (visiting.has(id)) throw new Error("Cycle at node: " + id);
    visiting.add(id);
    const depth = 1 + Math.max(...byID.get(id).after.map((p) => depthOf(p.id)));
    visiting.delete(id);
    depths.set(id, depth);
    return depth;
  }
  const branches = [...new Set(ordered.filter((n) => n !== root).map((n) => n.branch))].sort();
  const rings = new Map();
  for (const n of ordered) {
    const depth = depthOf(n.id);
    if (depth === 0) continue;
    if (!rings.has(depth)) rings.set(depth, new Map());
    const ring = rings.get(depth);
    if (!ring.has(n.branch)) ring.set(n.branch, []);
    ring.get(n.branch).push(n);
  }
  positions[root.id] = { x: 0, y: 0, depth: 0 };
  const R = 180;
  const minimumDistance = 96;
  const sector = 2 * Math.PI / branches.length;
  let radius = 0;
  for (const depth of [...rings.keys()].sort((a, b) => a - b)) {
    const angles = [];
    branches.forEach((branch, b) => {
      const peers = rings.get(depth).get(branch) || [];
      peers.forEach((n, i) => angles.push({ id: n.id, angle: sector * (b + (i + 1) / (peers.length + 1)) }));
    });
    radius += R;
    // Expanding the whole ring preserves parent ordering across branches and later depths.
    if (angles.length > 1) {
      for (let i = 0; i < angles.length; i++) {
        const next = angles[(i + 1) % angles.length].angle;
        const gap = (next - angles[i].angle + 2 * Math.PI) % (2 * Math.PI);
        radius = Math.max(radius, (minimumDistance + 1) / (2 * Math.sin(gap / 2)));
      }
    }
    for (const n of angles) {
      positions[n.id] = { x: radius * Math.cos(n.angle), y: radius * Math.sin(n.angle), depth };
    }
  }
  return positions;
}
