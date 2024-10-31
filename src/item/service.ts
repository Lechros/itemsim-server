import { StatusError, json, png } from 'itty-router';
import * as itemRepository from './repository';

function getIconRawOrigin(id: number) {
  const origin = itemRepository.findRawOriginById(id);
  if (!origin) {
    throw new StatusError(404);
  }
  return json(origin);
}

export { getIconRawOrigin };
