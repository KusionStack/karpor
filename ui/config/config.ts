import { defineConfig } from '@alipay/bigfish';
import { routes } from './routes';

const testServer = 'http://miaosang.alipay.net:7443';
const stableServer = 'http://sitebuilder-5.gz00b.stable.alipay.net:8090';
const preServer = 'https://sitebuilder-pre.alipay.com';
const prodServer = 'https://sitebuilder.alipay.com';

const headers = {};

export default defineConfig({
  layout: {
    locale: false, // 默认开启，如无需菜单国际化可关闭
  },
  appType: 'console',
  deployMode: 'tern',
  favicons: ['https://i.alipayobjects.com/common/favicon/favicon.ico'],
  oneApi: {
    apps: [
      {
        name: 'afs2demo', // OneAPI 应用名
      },
    ],
  },
  ctoken: {},
  qiankun: {
    master: {
      enable: false,
    },
  },
  history: {
    type: 'browser'
  },
  routes,
  styledComponents: {},
  reactQuery: {},
  icons: {},
  codeSplitting: { jsStrategy: 'granularChunks' },
  npmClient: 'tnpm',
  antd: { import: false },
  tern: {
    proxy: {
      mode: 'cors',
      DEV: {
        '/apis': {
          // target: devServer,
          target: testServer,
          headers,
        },
      },
      STABLE: {
        '/apis': {
          target: stableServer,
          headers,
        },
      },
      TEST: {
        '/apis': {
          target: testServer,
          headers,
        },
      },
      PRE: {
        '/apis': {
          target: preServer,
          headers,
        },
      },
      PROD: {
        '/apis': {
          target: prodServer,
          headers,
        },
      },
    },
  },
});
