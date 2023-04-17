import { Icon, Link, styled } from '@alipay/bigfish';
import { Button, Layout, Row, Space, Typography } from '@alipay/bigfish/antd';
import React from '@alipay/bigfish/react';
import { WelcomeCard } from '@alipay/tech-ui';

interface Props {
  name: string;
}

const StyledLayout = styled(Layout)`
  background-color: transparent;
  .title {
    margin: 0 auto 3rem;
    font-weight: 200;
  }
`;

// 脚手架示例组件
const Guide: React.FC<Props> = (props) => {
  const { name } = props;
  return (
    <StyledLayout>
      <Row>
        <Typography.Title level={3} className="title">
          欢迎使用 <strong>{name}</strong> ！
        </Typography.Title>
      </Row>
      <Row>
        <WelcomeCard.Operation
          style={{ width: 960 }}
          icon="https://gw-office.alipayobjects.com/basement_prod/c83c53ab-515e-43e2-85d0-4d0da16f11ef.svg"
          leftTitle="关于 Bigfish"
          leftContent="Bigfish 是一个企业级前端研发框架，并提供了一套完整的最佳实践。写前端，用 Bigfish 就够了！"
          leftActions={[
            <Space key="0" style={{ marginBottom: 8 }}>
              <Button
                href="https://bigfish.antgroup-inc.cn/"
                target="_blank"
                type="primary"
                icon={
                  <Icon
                    icon="local:home"
                    hover=""
                    fill="white"
                    style={{
                      width: '1rem',
                      height: '1rem',
                      verticalAlign: 'text-top',
                      marginRight: '4px',
                    }}
                  />
                }
              >
                Bigfish 官网
              </Button>
            </Space>,
            <Space key="1">
              <Link to="/access">
                <Button
                  icon={
                    <Icon
                      icon="local:home"
                      hover=""
                      style={{
                        width: '1rem',
                        height: '1rem',
                        verticalAlign: 'text-top',
                        marginRight: '4px',
                      }}
                    />
                  }
                >
                  权限演示
                </Button>
              </Link>
              <Link to="/table">
                <Button>CRUD 示例</Button>
              </Link>
              <Link to="/editor">
                <Button>Editor 示例</Button>
              </Link>
            </Space>,
            <p key="1" style={{ marginTop: 16, fontSize: 13 }}>
              <i>
                注：新增燕鸥应用默认不开启 ProLayout 布局，如需开启请参考
                <a
                  href="https://bigfish.antgroup-inc.cn/docs/guides/layout-menu"
                  target="_blank"
                  rel="noreferrer"
                >
                  官网文档
                </a>
                。
              </i>
            </p>,
          ]}
          leftSpan={11}
          rightSpan={13}
          rightTitle="相关帮助"
          rightContent={
            <WelcomeCard.QuickLinks
              links={[
                {
                  text: 'Ant Design - 企业级产品设计体系，创造高效愉悦的工作体验',
                  link: 'https://ant.design',
                },
                {
                  text: 'TechUI - 基于 Ant Design 的蚂蚁集团企业级 UI 资产库',
                  link: 'https://techui.alipay.com',
                },
                {
                  text: 'AntV - 企业级可视分析前端类库',
                  link: 'https://antv.alipay.com',
                },
                {
                  text: 'OneAPI - 前后端联调标准方案',
                  link: 'https://oneapi.alipay.com',
                },
                {
                  text: '雨燕 - 你的大前端工作台',
                  link: 'https://yuyan.antfin-inc.com/',
                },
              ]}
            />
          }
        />
      </Row>
    </StyledLayout>
  );
};

export default Guide;
