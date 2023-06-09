import { error, json } from 'itty-router';
import { cacheOk } from './middlewares/cache';
import { corsify } from './middlewares/cors';
import router from './router';

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    if (request.method === 'GET') {
      const cacheUrl = new URL(request.url);

      const cacheKey = new Request(cacheUrl.toString(), request);
      const cfCache = caches.default;

      let response = await cfCache.match(cacheKey);

      if (!response) {
        response = await router.handle(request, env, ctx).then(json).then(cacheOk).catch(error).then(corsify);

        ctx.waitUntil(cfCache.put(cacheKey, response!.clone()));
      }

      return response!;
    } else {
      return await router.handle(request, env, ctx).then(json).then(cacheOk).catch(error).then(corsify);
    }
  },
};
