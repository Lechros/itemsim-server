import { IRequest, Router, error, withContent } from 'itty-router';
import { withId } from '../middlewares/id';
import { CF } from '../types';
import * as gearService from './service';

const router = Router<IRequest, CF>({ base: '/gears' })
  .get('/search', (req) => {
    const { query } = req.query;
    if (typeof query !== 'string') {
      return error(400);
    }
    return gearService.search(query);
  })
  .get('/:id', withId, (req) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    return gearService.get(id);
  })
  .get('/:id/icon/origin', withId, (req) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    return gearService.getIconOrigin(id);
  })
  .get('/:id/iconRaw/origin', withId, (req) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    return gearService.getIconOrigin(id);
  });
// .post('/migrate', withContent, (req) => {
//   const { content } = req;
//   return gearService.getMigratedGear(content);
// });

export default router;
