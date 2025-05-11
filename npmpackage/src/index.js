import path from "path";
import fs from "fs";
import { spawnSync } from "child_process";
import downloadEXE from "./download.js";

const __dirname = path.resolve(); // stupid mjs limits
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

if (!fs.existsSync(filePath)) {
  await downloadEXE(appInfo.url, filePath, appInfo.hash, appInfo.name);
}

const packagePath = filePath;

spawnSync(packagePath, input, {
  stdio: "inherit",
  shell: true,
});
