"use strict";

const TREE_CARD_WIDTH = 132, TREE_CARD_HEIGHT = 44;

function layoutBounds(positions) {
  const points = Object.values(positions);
  if (!points.length) return { minX: 0, minY: 0, maxX: 0, maxY: 0 };
  return {
    minX: Math.min(...points.map((p) => p.x)) - TREE_CARD_WIDTH / 2,
    minY: Math.min(...points.map((p) => p.y)) - TREE_CARD_HEIGHT / 2,
    maxX: Math.max(...points.map((p) => p.x)) + TREE_CARD_WIDTH / 2,
    maxY: Math.max(...points.map((p) => p.y)) + TREE_CARD_HEIGHT / 2,
  };
}

// Intersecting the two axis intervals gives the radii where cards would overlap.
function overlapRadii(direction, point) {
  let start = -Infinity, end = Infinity;
  for (const [axis, size] of [["x", TREE_CARD_WIDTH], ["y", TREE_CARD_HEIGHT]]) {
    if (Math.abs(direction[axis]) < 1e-12) {
      if (Math.abs(point[axis]) >= size) return null;
    } else {
      const a = (point[axis] - size) / direction[axis], b = (point[axis] + size) / direction[axis];
      start = Math.max(start, Math.min(a, b)); end = Math.min(end, Math.max(a, b));
    }
  }
  return start < end ? { start, end } : null;
}

function ringRadius(directions, positions, minimum) {
  let radius = minimum;
  const blocked = [];
  for (let i = 0; i < directions.length; i++) {
    const a = directions[i];
    for (const b of directions.slice(0, i)) {
      radius = Math.max(radius, Math.min(TREE_CARD_WIDTH / Math.abs(a.x - b.x), TREE_CARD_HEIGHT / Math.abs(a.y - b.y)));
    }
    for (const point of Object.values(positions)) {
      const interval = overlapRadii(a, point);
      if (interval) blocked.push(interval);
    }
  }
  for (const interval of blocked.sort((a, b) => a.start - b.start)) {
    if (radius > interval.start - 1e-8 && radius < interval.end) radius = interval.end + 1e-6;
  }
  return radius + 1e-6;
}

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
  const weights = branches.map((branch) => Math.max(1, ...[...rings.values()].map((ring) => (ring.get(branch) || []).length)));
  const totalWeight = weights.reduce((sum, weight) => sum + weight, 0);
  const sectors = new Map();
  let start = 0;
  branches.forEach((branch, i) => {
    const span = 2 * Math.PI * weights[i] / totalWeight;
    sectors.set(branch, { start, span }); start += span;
  });
  const anglesByID = new Map();
  for (const depth of [...rings.keys()].sort((a, b) => a - b)) {
    const directions = [];
    let minimum = 0;
    branches.forEach((branch) => {
      const peers = rings.get(depth).get(branch) || [];
      const sector = sectors.get(branch), middle = sector.start + sector.span / 2;
      const parentAngle = (node) => {
        const parents = node.after.map((p) => anglesByID.get(p.id)).filter((angle) => angle !== undefined);
        if (!parents.length) return middle;
        const mean = Math.atan2(parents.reduce((sum, angle) => sum + Math.sin(angle), 0), parents.reduce((sum, angle) => sum + Math.cos(angle), 0));
        return middle + Math.atan2(Math.sin(mean - middle), Math.cos(mean - middle));
      };
      const sorted = peers.map((node) => ({ node, angle: parentAngle(node) })).sort((a, b) => a.angle - b.angle || (a.node.id < b.node.id ? -1 : a.node.id > b.node.id ? 1 : 0));
      sorted.forEach(({ node }, i) => {
        const angle = sector.start + sector.span * (i + 1) / (peers.length + 1);
        anglesByID.set(node.id, angle);
        directions.push({ id: node.id, x: Math.cos(angle), y: Math.sin(angle) });
        for (const p of node.after) minimum = Math.max(minimum, Math.hypot(positions[p.id].x, positions[p.id].y) + 1e-6);
      });
    });
    const radius = ringRadius(directions, positions, minimum);
    for (const n of directions) {
      positions[n.id] = { x: radius * n.x, y: radius * n.y, depth };
    }
  }
  return positions;
}
