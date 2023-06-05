import { createGearFromId, createGearFromNode, createPotentialFromCode } from '@malib/create-gear';
import { Gear, GearPropType, gearToPlain, isGearLike, migrate, plainToGear, validateParseGear } from '@malib/gear';
import { StatusError, json, png } from 'itty-router';
import { etag } from '../middlewares/cache';
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
  const gear = gearRepository.findById(id);
  if (!gear) {
    throw new StatusError(404);
  }
  return json(toDto(gear));
}

async function getIcon(id: number, bucket: R2Bucket) {
  const icon = await bucket.get(`gears/icon/${id}.png`);
  if (icon === null) {
    throw new StatusError(404);
  }
  return etag(png(icon.body), icon.etag);
}

function getIconOrigin(id: number) {
  const origin = gearRepository.findOriginById(id);
  if (!origin) {
    throw new StatusError(404);
  }
  return json(origin);
}

async function getIconRaw(id: number, bucket: R2Bucket) {
  const icon = await bucket.get(`gears/iconRaw/${id}.png`);
  if (icon === null) {
    throw new StatusError(404);
  }
  return etag(png(icon.body), icon.etag);
}

function getIconRawOrigin(id: number) {
  const origin = gearRepository.findRawOriginById(id);
  if (!origin) {
    throw new StatusError(404);
  }
  return json(origin);
}

function getMigratedGear(content: unknown) {
  if (!isGearLike(content)) {
    throw new StatusError(400);
  }
  const gear = plainToGear(content);
  const newGear = createGearFromId(gear.itemID)!;
  migrate(gear, newGear, {
    ignorePropTypes: [GearPropType.equipTradeBlock],
    getPotentialFunc: createPotentialFromCode,
  });
  return gearToPlain(newGear);
}

function toDto(entity: GearEntity): GearDto {
  const gear = createGearFromNode(entity, entity.id, createPotentialFromCode);
  return gearToPlain(gear);
}

export { search, get, getIcon, getIconOrigin, getIconRaw, getIconRawOrigin, getMigratedGear };
