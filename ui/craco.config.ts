const CracoLessPlugin = require("craco-less");
const CracoCSSModules = require("craco-css-modules");
const path  = require("path");

module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
      options: {
        lessLoaderOptions: {
          lessOptions: {
            modifyVars: {},
            javascriptEnabled: true,
          },
        },
      },
    },
    {
      plugin: CracoCSSModules
    }
  ],
  webpack: {
    // 配置别名，设置别名是为了让后续引用的地方减少路径的复杂度
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
}
