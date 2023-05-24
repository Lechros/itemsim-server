import { error, json } from 'itty-router';
import { corsify } from './middlewares/cors';
import router from './router';

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    return router.handle(request, env, ctx).then(json).catch(error).then(corsify);
  },
};
