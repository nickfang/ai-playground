const app = require('./app');
const config = require('./config');
const logger = require('./loaders/logger');

app.listen(config.port, () => {
  logger.info(`
      ################################################
      🛡️  Server listening on port: ${config.port} 🛡️
      ################################################
    `);
}).on('error', (err) => {
  logger.error(err);
  process.exit(1);
});
