import itemRawOrigin from '../data/item-raw-origin.json';
import { ItemIconOrigin } from './item';

const rawOriginDb = itemRawOrigin as { [id: number]: (typeof itemRawOrigin)[keyof typeof itemRawOrigin] };

function findRawOriginById(id: number): ItemIconOrigin | undefined {
  if (!rawOriginDb.hasOwnProperty(id)) {
    return undefined;
  }
  return rawOriginDb[id] as ItemIconOrigin;
}

export { findRawOriginById };
