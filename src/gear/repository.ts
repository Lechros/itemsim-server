import { Gear, GearData, GearDataMap, GearRepository, PotentialRepository } from '@malib/gear';
import gearOrigin from '../data/gear-origin.json';
import gearRawOrigin from '../data/gear-raw-origin.json';
import gears from '../data/gear.json';
import itemOptions from '../data/item-option.json';
import { InvertedIndex } from '../InvertedIndex';
import { GearIconOrigin } from './data';

const db = gears as GearDataMap;

const originDb = gearOrigin as { [id: number]: (typeof gearOrigin)[keyof typeof gearOrigin] };

const rawOriginDb = gearRawOrigin as { [id: number]: (typeof gearRawOrigin)[keyof typeof gearRawOrigin] };

const gearRepository = new GearRepository(gears, new PotentialRepository(itemOptions));

const gearIndex = new InvertedIndex<[string, GearData]>();

Object.entries(db).forEach((item) => {
  gearIndex.add(item, item[1].name);
});

function findById(id: number): Gear | undefined {
  return gearRepository.createGearFromId(id);
}

function findByName(keyword: string): Gear[] {
  return gearIndex
    .get(keyword)
    .filter((item) => match(item[1].name, keyword))
    .sort(([, d1], [, d2]) => compare(d1.name, d2.name, keyword))
    .map(([id]) => gearRepository.createGearFromId(Number(id)) as Gear);
}

function findOriginById(id: number): GearIconOrigin | undefined {
  if (!originDb.hasOwnProperty(id)) {
    return undefined;
  }
  return originDb[id] as GearIconOrigin;
}

function findRawOriginById(id: number): GearIconOrigin | undefined {
  if (!rawOriginDb.hasOwnProperty(id)) {
    return undefined;
  }
  return rawOriginDb[id] as GearIconOrigin;
}

function match(name: string, keyword: string) {
  let j = 0;
  for (let i = 0; i < name.length && j < keyword.length; i++) {
    if (name[i] === keyword[j]) {
      j++;
    }
  }
  return j === keyword.length;
}

function compare(name1: string, name2: string, keyword: string) {
  // both names are same as keyword, preserve itemID ordering
  if (name1 == keyword && name2 == keyword) return 0;
  // return exact match first
  if (name1 == keyword) return -1;
  if (name2 == keyword) return 1;
  const index1 = name1.indexOf(keyword);
  const index2 = name2.indexOf(keyword);
  // both names contain word
  if (index1 >= 0 && index2 >= 0) {
    // both have same word position
    if (index1 == index2) return 0;
    // name with keyword appearing earlier comes front
    return index1 - index2;
  }
  // only one name contains word
  if (index1 >= 0) return -1;
  if (index2 >= 0) return 1;
  // both are only fuzzily-matched, preserve itemID ordering
  return 0;
}

export { findById, findByName, findOriginById, findRawOriginById };
