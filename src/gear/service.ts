import { createGearFromNode, createPotentialFromCode } from '@malib/create-gear';
import { gearToPlain } from '@malib/gear';
import { StatusError, json, png } from 'itty-router';
import { cache } from '../middlewares/cache';
import { GearDto, GearEntity } from './gear';
import * as gearRepository from './repository';

function search(query: string) {
  query = query.trim();
  if (query.length === 0) {
    return json([]);
  }
  return json(gearRepository.findByName(query).map(toDto));
}

function get(id: number) {
  const result = gearRepository.findById(id);
  if (!result) {
    throw new StatusError(404);
  }
  return json(toDto(result));
}

async function getIcon(id: number, bucket: R2Bucket) {
  const icon = await bucket.get(`gear/${id}.png`);
  if (icon === null) {
    throw new StatusError(404);
  }
  const res = cache(png(icon.body), { maxAge: 86400, cachePublic: true });
  res.headers.set('Etag', icon.etag);
  return res;
}

function getOrigin(id: number) {
  const result = gearRepository.findById(id);
  if (!result) {
    throw new StatusError(404);
  }
  return json(result.origin);
}

function toDto(entity: GearEntity): GearDto {
  const gear = createGearFromNode(entity, entity.id, createPotentialFromCode);
  return gearToPlain(gear);
}

export { search, get, getIcon, getOrigin };
