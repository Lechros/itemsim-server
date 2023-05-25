export interface CacheOption {
  maxAge?: number;
  cachePublic?: boolean;
}

export function cache(response: Response, options?: CacheOption): Response {
  const { maxAge, cachePublic } = { maxAge: 86400, ...options };
  const value = [`max-age=${maxAge}`];
  if (cachePublic !== undefined) {
    value.push(getHttpPublic(cachePublic));
  }

  response.headers.set('Cache-Control', value.join(', '));
  return response;
}

function getHttpPublic(cachePublic?: boolean) {
  if (cachePublic === undefined) {
    return '';
  }
  if (cachePublic === true) {
    return 'public';
  }
  return 'private';
}