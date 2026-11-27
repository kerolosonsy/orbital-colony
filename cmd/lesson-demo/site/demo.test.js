"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { runRoute, goSource } = require("./demo");

const firstSignal = { Width: 6, Height: 3, Optimal: 3, Start: { X: 0, Y: 0 }, Beacon: { X: 3, Y: 0 } };

test("three East calls reach the beacon using the shortest route", () => {
  assert.deepEqual(runRoute(firstSignal, ["East", "East", "East"]),
    { X: 3, Y: 0, blocked: 0, success: true, optimal: true });
});

test("a run is judged at its final position, even if it visits the target", () => {
  assert.equal(runRoute(firstSignal, ["East", "East", "East", "West"]).success, false);
});

test("north increases Y and moves outside the grid leave the drone in place", () => {
  assert.deepEqual(runRoute(firstSignal, ["West", "North", "East"]),
    { X: 1, Y: 1, blocked: 1, success: false, optimal: false });
});

test("the code preview keeps each direction in program order", () => {
  assert.match(goSource(["East", "North"]), /colony\.Move\(colony\.East\)\n    colony\.Move\(colony\.North\)/);
});
