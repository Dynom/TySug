import { defineConfig, godoc } from "sourcey";

export default defineConfig({
  name: "TySug",
  siteUrl: "https://tysug.net",
  baseUrl: "/docs/reference",
  repo: "https://github.com/Dynom/TySug",
  editBranch: "master",
  theme: {
    preset: "api-first",
    colors: {
      primary: "#2f855a",
      light: "#48bb78",
      dark: "#22543d"
    }
  },
  navigation: {
    tabs: [
      {
        tab: "Go API",
        slug: "",
        source: godoc({
          module: "../..",
          packages: ["./..."],
          mode: "live",
          includeTests: true,
          includeUnexported: false,
          hideUndocumented: false
        })
      }
    ]
  },
  navbar: {
    links: [
      {
        type: "github",
        href: "https://github.com/Dynom/TySug"
      }
    ]
  },
  footer: {
    links: [
      {
        type: "github",
        href: "https://github.com/Dynom/TySug"
      }
    ]
  }
});
