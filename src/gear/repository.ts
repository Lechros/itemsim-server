import { gearJson } from '@malib/create-gear';
import { GearEntity } from './gear';

const db = Object.entries(gearJson);

function findById(id: number): GearEntity | undefined {
  if (!gearJson.hasOwnProperty(id)) {
    return undefined;
  }
  return { ...gearJson[id], id };
}

function findByName(keyword: string): GearEntity[] {
  return db
    .filter((e) => match(e[1].name, keyword))
    .sort((a, b) => compare(a[1].name, b[1].name, keyword))
    .map((g) => ({ ...g[1], id: Number(g[0]) }));
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

export { findById, findByName };
