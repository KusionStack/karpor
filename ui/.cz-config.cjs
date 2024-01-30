module.exports = {
  types: [
    { value: 'feat', name: '新增: 新增功能、页面' },
    { value: 'fix', name: 'bug: 修复某个bug' },
    { value: 'docs', name: '文档: 修改增加文档、注释' },
    { value: 'style', name: '格式: 不影响代码运行的变动、空格、格式化等等' },
    { value: 'ui', name: 'ui修改: 布局、css样式等等' },
    { value: 'hotfix', name: 'bug: 修复线上紧急bug' },
    { value: 'build', name: 'build: 变更项目构建或外部依赖' },
    { value: 'refactor', name: '重构: 代码重构,未新增任何功能和修复任何bug' },
    { value: 'revert', name: '回滚: 代码回退到某个版本节点' },
    { value: 'perf', name: '优化: 提升性能、用户体验等' },
    { value: 'ci', name: '自动化构建: 对CI/CD配置文件和脚本的更改' },
    { value: 'chore', name: 'chore: 变更构建流程或辅助工具' },
    { value: 'test', name: '测试: 包括单元测试、集成测试' },
    { value: 'update', name: '更新: 普通更新' },
  ],
  scopes: [],
  allowTicketNumber: false,
  isTicketNumberRequired: false,
  ticketNumberPrefix: 'TICKET-',
  ticketNumberRegExp: '\\d{1,5}',
  // it needs to match the value for field type. Eg.: 'fix'
  /*
  scopeOverrides: {
    fix: [
      {name: 'merge'},
      {name: 'style'},
      {name: 'e2eTest'},
      {name: 'unitTest'}
    ]
  },
  */
  // override the messages, defaults are as follows
  messages: {
    type: '选择一种你的提交类型:',
    scope: '选择一个scope (可选):',
    // used if allowCustomScopes is true
    customScope: '表示此更改的范围:',
    subject: '简短说明(最多40个字):\n',
    body: '长说明,使用"|"换行(可选):\n',
    breaking: '非兼容性说明 (可选):\n',
    footer: '关联关闭的issue,例如:#12, #34(可选):\n',
    confirmCommit: '确定提交?',
  },
  allowCustomScopes: true,
  // 设置选择了那些type,才询问 breaking message
  allowBreakingChanges: ['feat', 'fix', 'ui', 'hotfix', 'update', 'perf'],
  // 想跳过的问题
  skipQuestions: ['scope', 'body', 'breaking'],
  // limit subject length
  subjectLimit: 100,
}
