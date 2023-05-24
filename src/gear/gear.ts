import { GearData } from '@malib/create-gear';
import { GearLike } from '@malib/gear';

export type GearEntity = GearData & { id: number };

export type GearDto = GearLike;

export type GearIconOrigin = [number, number];
