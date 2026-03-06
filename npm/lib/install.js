#!/usr/bin/env node

"use strict";

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");

const PACKAGE = require("../package.json");
const VERSION = PACKAGE.version;
const BINARY_NAME = "terminalcrypto";
const GITHUB_REPO = "Carpe-Wang/terminalCrypto";

const PLATFORM_MAPPING = {
  darwin: "darwin",
  linux: "linux",
};

const ARCH_MAPPING = {
  x64: "amd64",
  arm64: "arm64",
};

function getPlatformInfo() {
  const platform = PLATFORM_MAPPING[process.platform];
  const arch = ARCH_MAPPING[process.arch];

  if (!platform || !arch) {
    console.error(
      `Unsupported platform: ${process.platform}-${process.arch}`
    );
    console.error(
      "terminalcrypto supports: darwin-amd64, darwin-arm64, linux-amd64, linux-arm64"
    );
    process.exit(1);
  }

  return { platform, arch };
}

function getDownloadUrl(platform, arch) {
  const filename = `${BINARY_NAME}-v${VERSION}-${platform}-${arch}.tar.gz`;
  return `https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${filename}`;
}

function download(url) {
  return new Promise((resolve, reject) => {
    const follow = (url) => {
      https
        .get(url, (res) => {
          if (res.statusCode === 301 || res.statusCode === 302) {
            follow(res.headers.location);
            return;
          }
          if (res.statusCode !== 200) {
            reject(new Error(`Download failed (HTTP ${res.statusCode}): ${url}`));
            return;
          }
          const chunks = [];
          res.on("data", (chunk) => chunks.push(chunk));
          res.on("end", () => resolve(Buffer.concat(chunks)));
          res.on("error", reject);
        })
        .on("error", reject);
    };
    follow(url);
  });
}

async function main() {
  const { platform, arch } = getPlatformInfo();
  const url = getDownloadUrl(platform, arch);

  const binDir = path.join(__dirname, "..", "bin");
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const binaryDest = path.join(binDir, `${BINARY_NAME}-${platform}-${arch}`);

  if (fs.existsSync(binaryDest)) {
    console.log(`${BINARY_NAME} binary already exists, skipping download.`);
    return;
  }

  console.log(
    `Downloading ${BINARY_NAME} v${VERSION} for ${platform}-${arch}...`
  );

  try {
    const buffer = await download(url);

    const tmpFile = path.join(binDir, "_tmp.tar.gz");
    fs.writeFileSync(tmpFile, buffer);

    try {
      execSync(`tar -xzf "${tmpFile}" -C "${binDir}"`, { stdio: "pipe" });
    } finally {
      if (fs.existsSync(tmpFile)) {
        fs.unlinkSync(tmpFile);
      }
    }

    // The archive contains a binary named "terminalcrypto", rename to include platform suffix
    const extracted = path.join(binDir, BINARY_NAME);
    if (fs.existsSync(extracted) && extracted !== binaryDest) {
      fs.renameSync(extracted, binaryDest);
    }

    fs.chmodSync(binaryDest, 0o755);
    console.log(`Successfully installed ${BINARY_NAME} v${VERSION}`);
  } catch (err) {
    console.error(`Failed to download ${BINARY_NAME}:`, err.message);
    console.error("");
    console.error("You can download the binary manually from:");
    console.error(`  https://github.com/${GITHUB_REPO}/releases`);
    process.exit(1);
  }
}

main();
