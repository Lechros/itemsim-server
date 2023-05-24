import { createGearFromNode, createPotentialFromCode } from '@malib/create-gear';
import { gearToPlain } from '@malib/gear';
import { StatusError } from 'itty-router';
import { GearDto, GearEntity, GearIconOrigin } from './gear';
import * as gearRepository from './repository';

function search(query: string): GearDto[] {
  query = query.trim();
  if (query.length === 0) {
    return [];
  }
  return gearRepository.findByName(query).map(toDto);
}

function get(id: number): GearDto {
  const result = gearRepository.findById(id);
  if (!result) {
    throw new StatusError(404);
  }
  return toDto(result);
}

async function getIcon(id: number, bucket: R2Bucket) {
  const icon = await bucket.get(`gear/${id}.png`);
  if (icon === null) {
    throw new StatusError(404);
  }
  return icon.body;
}

function getOrigin(id: number): GearIconOrigin {
  const result = gearRepository.findById(id);
  if (!result) {
    throw new StatusError(404);
  }
  return result.origin;
}

function toDto(entity: GearEntity): GearDto {
  const gear = createGearFromNode(entity, entity.id, createPotentialFromCode);
  return gearToPlain(gear);
}

export { search, get, getIcon, getOrigin };
