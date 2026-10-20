import { readFile, writeFile } from "node:fs/promises";

const basePath = "/docs/reference/";
const html = {
  index: await readFile("../reference/index.html", "utf8"),
  cmdWeb: await readFile("../reference/pkg-cmd-web.html", "utf8"),
  finder: await readFile("../reference/pkg-finder.html", "utf8"),
  internalRwc: await readFile("../reference/pkg-internal-rwc.html", "utf8"),
  keyboard: await readFile("../reference/pkg-keyboard.html", "utf8"),
  server: await readFile("../reference/pkg-server.html", "utf8"),
  serverService: await readFile("../reference/pkg-server-service.html", "utf8"),
};
let replacements = 0;

const rewritten = Object.fromEntries(
  Object.entries(html).map(([name, contents]) => [
    name,
    contents.replaceAll(
      /href="\/(index\.html|pkg-[^"]+\.html)"/g,
      (_match, target) => {
        replacements += 1;
        return `href="${basePath}${target}"`;
      },
    ),
  ]),
);

if (replacements === 0) {
  throw new Error("No Sourcey navigation links were rewritten");
}

await writeFile("../reference/index.html", rewritten.index);
await writeFile("../reference/pkg-cmd-web.html", rewritten.cmdWeb);
await writeFile("../reference/pkg-finder.html", rewritten.finder);
await writeFile("../reference/pkg-internal-rwc.html", rewritten.internalRwc);
await writeFile("../reference/pkg-keyboard.html", rewritten.keyboard);
await writeFile("../reference/pkg-server.html", rewritten.server);
await writeFile(
  "../reference/pkg-server-service.html",
  rewritten.serverService,
);

console.log(`Rewrote ${replacements} navigation links for ${basePath}`);
