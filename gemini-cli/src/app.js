const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const routes = require('./api');

const app = express();

app.use(helmet());
app.use(cors());
app.use(express.json());

app.use('/', routes);

module.exports = app; // test
