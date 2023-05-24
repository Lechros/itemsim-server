import itemOrigin from '../data/item-origin.json';
import { ItemIconOrigin } from './item';

const db = itemOrigin as { [id: number]: (typeof itemOrigin)[keyof typeof itemOrigin] };

function findById(id: number): ItemIconOrigin | undefined {
  if (!db.hasOwnProperty(id)) {
    return undefined;
  }
  return db[id] as [number, number];
}

export { findById };
