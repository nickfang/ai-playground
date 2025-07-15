// Set the NODE_ENV to 'development' by default
process.env.NODE_ENV = process.env.NODE_ENV || 'development';

module.exports = {
  /**
   * Your favorite port
   */
  port: parseInt(process.env.PORT, 10) || 3000,

  /**
   * That long string that is supposed to be secret
   */
  jwtSecret: process.env.JWT_SECRET,

  /**
   * Your secret jwt algorithm
   */
  jwtAlgorithm: process.env.JWT_ALGORITHM || 'RS256',

  /**
   * Used by winston logger
   */
  logs: {
    level: process.env.LOG_LEVEL || 'silly',
  },

  /**
   * API configs
   */
  api: {
    prefix: '/api',
  },
};
