// 运行时配置
import services from '@/services/afs2demo';
import {
  createGlobalStyle,
  RequestConfig,
  RunTimeLayoutConfig,
} from '@alipay/bigfish';
import { HeatMapOutlined } from "@ant-design/icons"

// 全局初始化数据配置，用于 Layout 用户信息和权限初始化
// 更多信息见文档：https://bigfish.antfin-inc.com/doc/console-initial-state
export async function getInitialState(): Promise<API.UserInfo> {
  const { data } = await services.UserController.getUserDetail({
    userId: '1',
  });
  return { ...data };
}

// request 配置参考文档：https://bigfish.antfin-inc.com/doc/console-request
export const request: RequestConfig = {
  requestInterceptors: [
    // oneapi mock server
    // (config: any) => {
    //   if (process.env.NODE_ENV !== 'development') return config;
    //   if (!config.url.startsWith('/')) return config;
    //   const appName = 'afs2demo';
    //   const tag = 'master';
    //   const source = 'ZAPPINFO';
    //   const scene = 'default';
    //   // &mode=static
    //   config.url = `https://oneapitwa.alipay.com/api/mock/proxy?appName=${appName}&tag=${tag}&source=${source}&scene=${scene}&path=${
    //     config.url.split('?')[0]
    //   }&method=${config.method}`;
    //   return config;
    // },
  ],
};

export const styledComponents = {
  GlobalStyle: createGlobalStyle`
html,
body {
  height: 100%;
}
// qiankun 会修改主应用根节点为 root-master
// 若不需要撑开高度，可以自行删除这个样式
#root,
#root-master {
  height: 100%;
}
  `,
};

export const reactQuery = {
  queryClient: {
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  },
};

//  initialState
export const layout: RunTimeLayoutConfig = () => {
  return {
    // 常用属性
    title: 'Karbour',
    logo: <HeatMapOutlined style={{color: '#000'}} />,

    // 默认布局调整
    // rightContentRender: () => <RightContent />,
    // footerRender: () => <Footer />,
    menuHeaderRender: undefined,
    layout: 'top',

    // 其他属性见：https://procomponents.ant.design/components/layout#prolayout
  };
};
