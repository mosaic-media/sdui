#!/usr/bin/env node
/**
 * Fail when a publishable package's version disagrees with the repository's
 * latest git tag.
 *
 * Go modules do not need this: a Go module *is* its tag, and a module can read
 * back whatever version was actually linked (`v1.ModuleVersion`), so the number
 * is structurally true rather than maintained. npm has no equivalent — a
 * `package.json` version is a literal somebody has to remember to change, and
 * a forgotten one publishes under the wrong number or, worse, silently fails to
 * publish at all because the version already exists.
 *
 * This is not hypothetical. Two Go modules in this project drifted from their
 * tags before the constants were removed — one by two releases — and both were
 * wrong in the least visible way available: a stale version looks exactly like
 * a correct one. npm is the same shape of problem with no structural fix, so it
 * gets a check instead.
 *
 * The rule matches what both repositories already do rather than imposing
 * something new: **the latest tag names the version of whatever the repository
 * publishes.** Private packages are exempt — they are not published, so their
 * version is bookkeeping and nothing depends on it being right.
 *
 * Skipped entirely when there are no tags (a repository that has never
 * released) so a fresh clone or a first CI run does not fail on nothing.
 */
import { readFileSync, existsSync, readdirSync } from "node:fs";
import { execSync } from "node:child_process";

/** The newest tag by version order, or null when the repo has never tagged. */
function latestTag() {
  try {
    // Version-sorted rather than chronological: a patch tagged after a minor
    // is still the older release, and date order would call it the newest.
    const out = execSync("git tag --sort=-v:refname", { encoding: "utf8" }).trim();
    return out ? out.split("\n")[0] : null;
  } catch {
    return null;
  }
}

/** Every package.json this repo owns: the root, plus any npm workspaces.
 *
 * Expanded by hand rather than with fs.globSync, which needs Node 22 — CI runs
 * 20, and a version check that only runs on the maintainer's machine is the
 * kind of check that is quietly not running. Workspace patterns in practice are
 * `dir/*`, so that is what this handles, and anything else is reported rather
 * than skipped silently. */
function packageFiles() {
  const files = ["package.json"];
  const root = JSON.parse(readFileSync("package.json", "utf8"));
  for (const pattern of root.workspaces ?? []) {
    if (!pattern.endsWith("/*")) {
      console.warn(`check-versions: unsupported workspace pattern ${pattern} — not checked`);
      continue;
    }
    const base = pattern.slice(0, -2);
    if (!existsSync(base)) continue;
    for (const entry of readdirSync(base, { withFileTypes: true })) {
      if (!entry.isDirectory()) continue;
      const p = `${base}/${entry.name}/package.json`;
      if (existsSync(p)) files.push(p);
    }
  }
  return files;
}

const tag = latestTag();
if (!tag) {
  console.log("check-versions: no tags yet — nothing to check against");
  process.exit(0);
}
const tagVersion = tag.replace(/^v/, "");

const problems = [];
let checked = 0;

for (const file of packageFiles()) {
  const pkg = JSON.parse(readFileSync(file, "utf8"));
  if (pkg.private) continue; // not published; its version means nothing
  checked++;
  if (pkg.version !== tagVersion) {
    problems.push(
      `  ${file}\n    ${pkg.name} is ${pkg.version}, but the latest tag is ${tag}`,
    );
  }
}

if (checked === 0) {
  console.log("check-versions: no publishable packages — nothing to check");
  process.exit(0);
}

if (problems.length > 0) {
  console.error(
    `check-versions: a published package disagrees with the latest tag (${tag}):\n\n` +
      problems.join("\n") +
      `\n\nEither bump the package version to ${tagVersion}, or tag the release\n` +
      `the package version claims. They have to agree: the tag is what a\n` +
      `consumer resolves and the package version is what npm serves, and when\n` +
      `they differ one of them is lying about what the other contains.\n`,
  );
  process.exit(1);
}

console.log(`check-versions: ${checked} published package(s) agree with ${tag}`);
