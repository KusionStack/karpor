const data = {
  default: {
    success: true,
    errorMessage: 'midimeza',
    data: { id: 'ge', name: 'feopa', nickName: 'hubucu', email: 'wig', gender: 'MALE' },
  },
  '自动生成 mock': {
    success: true,
    errorMessage: 'dRA',
    data: { id: 'tey', name: '1uV', nickName: 'hA0V', email: 'q^q', gender: 'MALE' },
  },
};

module.exports = {
  'GET /api/v1/user/:userId': (req, res) => {
    res.send(data.default);
  },
};
