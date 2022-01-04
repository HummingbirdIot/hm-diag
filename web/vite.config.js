import vue from "@vitejs/plugin-vue";
import styleImport from "vite-plugin-style-import";

export default {
  resolve: {
    alias: {
      "vue": "vue/dist/vue.esm-bundler.js"
    },
  },
  plugins: [
    vue(),
    styleImport({
      libs: [
        {
          libraryName: "vant",
          esModule: true,
          resolveStyle: (name) => `vant/es/${name}/style`,
        },
      ],
    }),
  ],
  server: {
    host: '0.0.0.0',
    proxy: {
      '/api/': {
        target: 'http://http://192.168.89.45',
        changeOrigin: true
      },

      '/state': {
        target: 'http://192.168.89.45',
        changeOrigin: true
      }
    }
  }
};
