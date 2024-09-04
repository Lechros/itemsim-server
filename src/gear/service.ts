import { Gear, gearToPlain } from '@malib/gear';
import { StatusError, json, png } from 'itty-router';
import { etag } from '../middlewares/cache';
import { GearRes } from './data';
import * as gearRepository from './repository';

function search(query: string) {
  query = query.trim();
  if (query.length === 0) {
    return json([]);
  }
  return json(gearRepository.findByName(query).map(toGearRes));
}

function get(id: number) {
  const gear = gearRepository.findById(id);
  if (!gear) {
    throw new StatusError(404);
  }
  return json(toGearRes(gear));
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

// function getMigratedGear(content: unknown) {
//   if (!isGearLike(content)) {
//     throw new StatusError(400);
//   }
//   const gear = plainToGear(content);
//   const newGear = gearRepository.findById(gear.itemID)!;
//   migrate(gear, newGear, {
//     ignorePropTypes: [GearPropType.equipTradeBlock],
//     getPotentialFunc: createPotentialFromCode,
//   });
//   return gearToPlain(newGear);
// }

function toGearRes(gear: Gear): GearRes {
  return gearToPlain(gear);
}

export { get, getIcon, getIconOrigin, getIconRaw, getIconRawOrigin, search };
