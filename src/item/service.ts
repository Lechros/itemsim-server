import { StatusError } from 'itty-router';
import { ItemIconOrigin } from './item';
import * as itemRepository from './repository';

async function getIcon(id: number, bucket: R2Bucket) {
  const icon = await bucket.get(`item/${id}.png`);
  if (icon === null) {
    throw new StatusError(404);
  }
  return icon.body;
}

function getOrigin(id: number): ItemIconOrigin {
  const result = itemRepository.findById(id);
  if (!result) {
    throw new StatusError(404);
  }
  return result;
}

export { getIcon, getOrigin };
