import { error, json } from 'itty-router';
import { corsify } from './middlewares/cors';
import router from './router';

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    const cacheUrl = new URL(request.url);

    const cacheKey = new Request(cacheUrl.toString(), request);
    const cache = caches.default;

    let response = await cache.match(cacheKey);

    if (!response) {
      response = await router.handle(request, env, ctx).then(json).catch(error).then(corsify);

      response = new Response(response!.body, response);

      ctx.waitUntil(cache.put(cacheKey, response.clone()));
    }

    return response;
  },
};
