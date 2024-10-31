import { GearData, GearRepository, PotentialRepository } from '@malib/gear';
import { InvertedIndex } from './InvertedIndex';

test('add로 추가된 아이템 1개을 get으로 가져올 수 있다.', () => {
  const index = new InvertedIndex<number>();
  index.add(1234, 'hello');
  expect(index.get('h')).toEqual([1234]);
});

test('add로 추가된 아이템 5개를 get으로 가져올 수 있다.', () => {
  const index = new InvertedIndex<string>();
  index.add('1004422', '앱솔랩스 나이트헬름');
  index.add('1004423', '앱솔랩스 메이지크라운');
  index.add('1004424', '앱솔랩스 아처후드');
  index.add('1004425', '앱솔랩스 시프캡');
  index.add('1004426', '앱솔랩스 파이렛페도라');
  index.add('1230001', '하이네스 워리어헬름');
  index.add('1230002', '하이네스 던위치햇');
  index.add('1230003', '하이네스 레인져베레');
  index.add('1230004', '하이네스 어새신보닛');
  index.add('1230005', '하이네스 원더러햇');

  expect(index.get('앱스')).toEqual(['1004422', '1004423', '1004424', '1004425', '1004426']);
  expect(index.get('앱솔랩스')).toEqual(['1004422', '1004423', '1004424', '1004425', '1004426']);
});

import gears from './data/gear.json';

const gearIndex = new InvertedIndex<[string, GearData]>();

Object.entries(gears).forEach((item) => {
  gearIndex.add(item, item[1].name);
});

test('test gearIndex', () => {
  expect(gearIndex.get('앱솔랩스')).not.empty;
});
