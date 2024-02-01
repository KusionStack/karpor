/* eslint-disable @typescript-eslint/no-var-requires */
const CracoLessPlugin = require('craco-less')
const CompressionPlugin = require('compression-webpack-plugin') //引入gzip压缩插件
// 多线程打包
const HappyPack = require('happypack')
const os = require('os')
const path = require('path')
// 开辟一个线程池，拿到系统CPU的核数，happypack 将编译工作利用所有线程
const happyThreadPool = HappyPack.ThreadPool({ size: os.cpus().length })

const { whenProd } = require('@craco/craco')

// 具体配置见官网：https://craco.js.org/docs/
module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
    },
  ],
  // webpack 配置
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
    configure: (webpackConfig, { paths }) => {
      // 修改打包输出文件目录
      paths.appBuild = path.resolve(__dirname, 'build')
      webpackConfig.output = {
        ...webpackConfig.output,
        clean: true, // 自动将上次打包目录资源清空
        path: path.resolve(__dirname, 'build'),
        publicPath: '/', //资源名
      }

      // 生产环境 才会下面配置
      whenProd(() => {
        // 删除log
        const TerserPlugin = webpackConfig.optimization.minimizer.find(
          i => i.constructor.name === 'TerserPlugin',
        )
        if (TerserPlugin) {
          // TerserPlugin.options.minimizer.options.compress['drop_console'] = true // 删除所有console语句
          TerserPlugin.options.minimizer.options.compress['drop_debugger'] =
            true
          TerserPlugin.options.minimizer.options.compress['pure_funcs'] = [
            'console.log',
          ] //删除打印语句
        }

        // webpack添加插件
        webpackConfig.plugins.push(
          // 配置完以后，暂时还不能使用，还需要后端做一下配置，这里后端以nginx为例
          // 使用gzip压缩超过1M的js和css文件
          new CompressionPlugin({
            // filename: "[path][base].gz", // 这种方式是默认的，多个文件压缩就有多个.gz文件
            algorithm: 'gzip', // 官方默认压缩算法也是gzip
            test: /\.(js|css)$/, // 使用正则给匹配到的文件做压缩，这里是给css、js
            threshold: 10240, //以字节为单位压缩超过此大小的文件，小于10KB就不进行压缩
            minRatio: 0.8, // 最小压缩比率，官方默认0.8
            //是否删除原有静态资源文件，即只保留压缩后的.gz文件，建议这个置为false，还保留源文件。以防：假如出现访问.gz文件访问不到的时候，还可以访问源文件双重保障
            deleteOriginalAssets: false,
          }),

          // 使用多线程打包
          new HappyPack({
            // id标识happyPack处理那一类文件
            id: 'babel',
            loaders: ['babel-loader'],
            // 共享进程池
            threadPool: happyThreadPool,
          }),
        )
      })
      return webpackConfig
    },
  },
}
