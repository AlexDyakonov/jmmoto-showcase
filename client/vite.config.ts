import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    allowedHosts: ["tunnel.shamps.dev", "timur.tunnel.shamps.dev", "mrussy.tunnel.shamps.dev"],
  },
  define: {
    global: {},
  },
});
