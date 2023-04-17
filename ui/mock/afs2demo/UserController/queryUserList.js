const data = {
  default: {
    success: true,
    errorMessage: 'ebojuro',
    data: {
      current: 54,
      pageSize: 37,
      total: 9,
      list: [
        { id: 'igawusfar', name: 'naar', nickName: 'pebosbas', email: 'rusulkoc', gender: 'MALE' },
        { id: 'zigijmi', name: 'murez', nickName: 'pafaj', email: 'nubvizpe', gender: 'MALE' },
        { id: 'ufosubom', name: 'la', nickName: 'galebgom', email: 've', gender: 'MALE' },
      ],
    },
  },
};

module.exports = {
  'GET /api/v1/queryUserList': (req, res) => {
    res.send(data.default);
  },
};
