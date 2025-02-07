/* eslint-disable @typescript-eslint/no-var-requires */
const CracoLessPlugin = require('craco-less')
const CompressionPlugin = require('compression-webpack-plugin')
const path = require('path')
const { whenProd } = require('@craco/craco')

module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
    },
  ],
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
    configure: (webpackConfig, { paths }) => {
      paths.appBuild = path.resolve(__dirname, 'build')

      whenProd(() => {
        const TerserPlugin = webpackConfig.optimization.minimizer.find(
          i => i.constructor.name === 'TerserPlugin',
        )
        if (TerserPlugin) {
          TerserPlugin.options.minimizer.options.compress['drop_debugger'] =
            true
          TerserPlugin.options.minimizer.options.compress['pure_funcs'] = [
            'console.log',
          ]
        }

        webpackConfig.plugins.push(
          new CompressionPlugin({
            algorithm: 'gzip',
            test: /\.(js|css)$/,
            threshold: 10240,
            minRatio: 0.8,
            deleteOriginalAssets: false,
          }),
        )
      })
      return webpackConfig
    },
  },
}
