const data = {
  default: {
    success: true,
    errorMessage: 'fiba',
    data: { id: 'vihmen', name: 'je', nickName: 'gunof', email: 'ficingus', gender: 'MALE' },
  },
  '自动生成 mock': {
    success: false,
    errorMessage: 'AJ0G',
    data: { id: 'F]J', name: '#O]', nickName: '*Yn3', email: 'EbC', gender: 'FEMALE' },
  },
};

module.exports = {
  'POST /api/v1/user': (req, res) => {
    res.send(data.default);
  },
};
