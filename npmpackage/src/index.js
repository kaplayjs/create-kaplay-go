#!/usr/bin/env node
import path from "path";
import fs from "fs";
import { spawnSync } from "child_process";
import downloadEXE from "./download.js";
import { fileURLToPath } from "url";
import crypto from "crypto";

const __dirname = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  ".."
); // stupid mjs limits
const platform = `${process.platform}_${process.arch}`;
const input = process.argv.slice(2);

const file = fs.readFileSync(path.join(__dirname, "app.json"), "utf8");
const app = JSON.parse(file);

if (!app[platform]) {
  console.error("Platform not supported");
  process.exit(1);
}

const appInfo = app[platform];
const filePath = path.join(__dirname, appInfo.executable);

if (fs.existsSync(filePath)) {
  const buffer = fs.readFileSync(filePath);
  const hash = crypto.createHash("sha256");
  hash.update(buffer);
  const digest = hash.digest("hex");
  if (digest !== appInfo.hash) {
    await downloadEXE(appInfo.url, filePath, appInfo.hash, appInfo.name);
  }
} else {
  await downloadEXE(appInfo.url, filePath, appInfo.hash, appInfo.name);
}

const packagePath = filePath;

spawnSync(packagePath, input, {
  stdio: "inherit",
  shell: true,
});
