const data = {
  default: { success: true, errorMessage: 'munkid', data: 'mogpag' },
  '自动生成 mock': { success: true, errorMessage: 't!pu23H', data: 'Y9[Npe' },
};

module.exports = {
  'DELETE /api/v1/user/:userId': (req, res) => {
    res.send(data.default);
  },
};
