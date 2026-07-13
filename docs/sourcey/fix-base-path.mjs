import { readdir, readFile, writeFile } from "node:fs/promises";
import { resolve } from "node:path";

const outputDir = resolve("../reference");
const basePath = "/docs/reference/";
const htmlFiles = (await readdir(outputDir)).filter((file) => file.endsWith(".html"));
let replacements = 0;

for (const file of htmlFiles) {
  const path = resolve(outputDir, file);
  const original = await readFile(path, "utf8");
  const updated = original.replaceAll(
    /href="\/(index\.html|pkg-[^"]+\.html)"/g,
    (_match, target) => {
      replacements += 1;
      return `href="${basePath}${target}"`;
    }
  );

  await writeFile(path, updated);
}

if (replacements === 0) {
  throw new Error("No Sourcey navigation links were rewritten");
}

console.log(`Rewrote ${replacements} navigation links for ${basePath}`);
