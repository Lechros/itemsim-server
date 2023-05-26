import { StatusError, json, png } from 'itty-router';
import { etag } from '../middlewares/cache';
import * as itemRepository from './repository';

async function getIconRaw(id: number, bucket: R2Bucket) {
  const icon = await bucket.get(`items/iconRaw/${id}.png`);
  if (icon === null) {
    throw new StatusError(404);
  }
  return etag(png(icon.body), icon.etag);
}

function getIconRawOrigin(id: number) {
  const origin = itemRepository.findRawOriginById(id);
  if (!origin) {
    throw new StatusError(404);
  }
  return json(origin);
}

export { getIconRaw, getIconRawOrigin };
