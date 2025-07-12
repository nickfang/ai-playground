const express = require('express');
const route = express.Router();

module.exports = (app) => {
  app.use('/health', route);

  route.get('/', (req, res) => {
    res.status(200).json({ status: 'UP' });
  });
};
