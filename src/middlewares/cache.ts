export interface CacheOption {
  maxAge?: number;
  cachePublic?: boolean | null;
}

export function cache(response: Response, options?: CacheOption): Response {
  const { maxAge, cachePublic } = { maxAge: 86400, cachePublic: true, ...options };
  let value = [`max-age=${maxAge}`];
  if (getHttpPublic(cachePublic) !== undefined) {
    value = [getHttpPublic(cachePublic)!, ...value];
  }
  response.headers.set('Cache-Control', value.join(', '));

  return response;
}

export function cacheOk(response: Response): Response {
  if (response.status < 400) {
    return cache(response, { maxAge: 86400, cachePublic: true });
  }
  return cache(response, { maxAge: 0, cachePublic: null });
}

export function etag(response: Response, etag: string) {
  response.headers.set('Etag', etag);

  return response;
}

function getHttpPublic(cachePublic: boolean | null) {
  if (cachePublic === null) {
    return undefined;
  }
  if (cachePublic === true) {
    return 'public';
  }
  return 'private';
}