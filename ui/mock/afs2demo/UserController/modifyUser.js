const data = {
  default: {
    success: true,
    errorMessage: 'ukozu',
    data: { id: 'cibso', name: 'uhvo', nickName: 'uge', email: 'wik', gender: 'MALE' },
  },
  '自动生成 mock': {
    success: true,
    errorMessage: '200T',
    data: { id: '7ERa', name: 't%Od5', nickName: 'AWJGZs', email: 'Uf#yvCI', gender: 'MALE' },
  },
};

module.exports = {
  'PUT /api/v1/user/:userId': (req, res) => {
    res.send(data.default);
  },
};
