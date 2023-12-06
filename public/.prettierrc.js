const fabric = require('@umijs/fabric');

module.exports = {
  ...fabric.prettier,
  importOrder: ['^@formily/(.*)', '^@(.*)$', '^[./]'],
  importOrderSeparation: true,
};
