#!/usr/bin/env node

"use strict";

const path = require("path");
const { execFileSync } = require("child_process");
const fs = require("fs");

const PLATFORM_MAPPING = {
  darwin: "darwin",
  linux: "linux",
};

const ARCH_MAPPING = {
  x64: "amd64",
  arm64: "arm64",
};

function getBinaryPath() {
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

  const binaryName = `terminalcrypto-${platform}-${arch}`;
  const binaryPath = path.join(__dirname, "..", "bin", binaryName);

  if (!fs.existsSync(binaryPath)) {
    console.error(`Binary not found: ${binaryPath}`);
    console.error("Try reinstalling: npm install -g terminalcrypto");
    process.exit(1);
  }

  return binaryPath;
}

try {
  const binaryPath = getBinaryPath();
  execFileSync(binaryPath, process.argv.slice(2), {
    stdio: "inherit",
    env: process.env,
  });
} catch (err) {
  if (err.status !== null && err.status !== undefined) {
    process.exit(err.status);
  }
  console.error(err.message);
  process.exit(1);
}
