const request = require('supertest');
const app = require('../../app');

describe('GET /', () => {
  it('should respond with a 200 status code and a welcome message', async () => {
    const response = await request(app).get('/');
    expect(response.statusCode).toBe(200);
    expect(response.text).toBe('Welcome to the Service!');
  });
});
