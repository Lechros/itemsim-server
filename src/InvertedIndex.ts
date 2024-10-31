type char = string;

export class InvertedIndex<T> {
  private map: Map<char, number[]>;
  private items: T[];

  constructor() {
    this.map = new Map();
    this.items = [];
  }

  add(item: T, text: string) {
    const index = this.items.length;
    this.items.push(item);

    for (const ch of text) {
      let indexes = this.map.get(ch);
      if (indexes === undefined) {
        indexes = [];
        this.map.set(ch, indexes);
      }
      indexes.push(index);
    }
  }

  get(text: string) {
    let result = this.map.get(text[0]);
    if (result === undefined || result.length === 0) {
      return [];
    }
    for (let i = 1; i < text.length; i++) {
      const indexes = this.map.get(text[i]);
      if (indexes === undefined) {
        return [];
      }
      result = intersect(result, indexes);
      if (result.length === 0) {
        return [];
      }
    }
    return result.map((index) => this.items[index]);
  }

  getItems() {
    return this.items;
  }
}

function intersect(a: number[], b: number[]): number[] {
  if (a.length === 0 || b.length === 0) {
    return [];
  }
  if (a[0] > b[b.length - 1]) {
    return [];
  }
  if (b[0] > a[a.length - 1]) {
    return [];
  }
  const result = [];
  for (let i = 0, j = 0; i < a.length && j < b.length; ) {
    if (a[i] === b[j]) {
      result.push(a[i]);
      i++;
      j++;
    } else if (a[i] < b[j]) {
      i++;
    } else {
      j++;
    }
  }
  return result;
}
