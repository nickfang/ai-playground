const express = require('express');
const health = require('./routes/health');
const welcome = require('./routes/welcome');

const app = express.Router();
health(app);
welcome(app);

module.exports = app;
