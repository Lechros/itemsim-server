import { createCors } from 'itty-router';

export const { preflight, corsify } = createCors({
  methods: ['GET', 'OPTIONS'],
  origins: ['https://itemsim.com', 'https://itemsim.pages.dev', 'http://localhost:5173'],
  headers: {
    'Access-Control-Allow-Credentials': 'true',
    'Access-Control-Allow-Headers':
      'DNT,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization',
  },
  maxAge: 86400,
});
