const request = require('supertest');
const app = require('../../app');

describe('GET /health', () => {
  it('should respond with a 200 status code and a JSON object with status UP', async () => {
    const response = await request(app).get('/health');
    expect(response.statusCode).toBe(200);
    expect(response.body).toEqual({ status: 'UP' });
  });
});
