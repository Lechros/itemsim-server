import { IRequest, Router, error, json, png } from 'itty-router';
import { withId } from '../middlewares/id';
import { CF } from '../types';
import * as gearService from './service';

const router = Router<IRequest, CF>({ base: '/gears' })
  .get('/search', (req) => {
    const { query } = req.query;
    if (typeof query !== 'string') {
      return error(400);
    }
    return json(gearService.search(query));
  })
  .get('/:id', withId, (req) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    return json(gearService.get(id));
  })
  .get<IRequest, CF>('/:id/icon', withId, async (req, env) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    const bucket = env.MY_BUCKET;
    if (!bucket) {
      return error(500);
    }
    return png(await gearService.getIcon(id, bucket));
  })
  .get('/:id/origin', withId, (req) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    return json(gearService.getOrigin(id));
  });

export default router;
