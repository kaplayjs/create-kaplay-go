import crypto from "crypto";
import axios from "axios";
import fs from "fs";
import ora from "ora";

export default async function downloadEXE(url, filePath, sha, name) {
  const spinner = ora(`Downloading create-kaplay for ${name}`);
  spinner.spinner = {
    interval: 80,
    frames: ["⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"],
  };
  spinner.start();
  return new Promise((resolve) => {
    axios
      .get(url, {
        responseType: "arraybuffer",
      })
      .then((response) => {
        const buffer = Buffer.from(response.data);

        const hash = crypto.createHash("sha256");
        hash.update(buffer);
        const digest = hash.digest("hex");
        if (digest !== sha) {
          spinner.fail(`SHA256 hash mismatch: expected ${sha}, got ${digest}`);
          process.exit(1);
        }

        fs.writeFileSync(filePath, buffer);
        spinner.stop();
        resolve();
      })
      .catch((err) => {
        spinner.fail("Download failed");
        console.error("Error downloading file:", err);
        process.exit(1);
      });
  });
}
