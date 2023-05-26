import { createCors } from 'itty-router';

export const { preflight, corsify } = createCors({
  methods: ['GET', 'OPTIONS'],
  origins: ['https://itemsim.com', 'https://itemsim.pages.dev', 'http://localhost:5173'],
  headers: {
    'Access-Control-Allow-Credentials': 'true',
    'Access-Control-Allow-Headers': 'Cache-Control,Content-Type,Range,Authorization',
  },
  maxAge: 86400,
});
