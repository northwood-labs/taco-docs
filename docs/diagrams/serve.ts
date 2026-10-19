/**
 * Development server for the architecture diagram.
 * Serves static files and bundles TypeScript on-the-fly.
 *
 * Usage: bun run serve.ts
 */

const server = Bun.serve({
  port: 3333,
  async fetch(req) {
    const url = new URL(req.url);
    let pathname = url.pathname;

    if (pathname === "/") {
      pathname = "/index.html";
    }

    // Serve TypeScript files as bundled JavaScript
    if (pathname.endsWith(".ts")) {
      const filePath = import.meta.dir + pathname;
      const file = Bun.file(filePath);

      if (await file.exists()) {
        const result = await Bun.build({
          entrypoints: [filePath],
          target: "browser",
          format: "esm",
        });

        if (result.success && result.outputs.length > 0) {
          const text = await result.outputs[0].text();
          return new Response(text, {
            headers: { "Content-Type": "application/javascript" },
          });
        }

        return new Response("Build failed: " + result.logs.join("\n"), {
          status: 500,
        });
      }
    }

    // Serve static files
    const filePath = import.meta.dir + pathname;
    const file = Bun.file(filePath);

    if (await file.exists()) {
      return new Response(file);
    }

    return new Response("Not Found", { status: 404 });
  },
});

console.log(`Architecture diagram server running at http://localhost:${server.port}`);
