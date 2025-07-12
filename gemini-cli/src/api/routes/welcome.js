const express = require('express');
const route = express.Router();

module.exports = (app) => {
  app.use('/', route);

  route.get('/', (req, res) => {
    res.status(200).send('Welcome to the Service!');
  });
};
