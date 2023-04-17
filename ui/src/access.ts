export default (initialState: API.UserInfo) => {
  // 在这里按照初始化数据定义项目中的权限，统一管理
  // 参考文档 https://bigfish.antgroup-inc.cn/docs/guides/access
  const canSeeAdmin = !!(initialState && initialState.name !== 'nothasaccess');
  return {
    canSeeAdmin,
  };
};
