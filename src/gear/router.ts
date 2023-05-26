import { IRequest, Router, error } from 'itty-router';
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
  .get<IRequest, CF>('/:id/icon', withId, (req, env, ctx) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    const bucket = env.MY_BUCKET;
    return gearService.getIcon(id, bucket);
  })
  .get('/:id/icon/origin', withId, (req) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    return gearService.getIconOrigin(id);
  })
  .get<IRequest, CF>('/:id/iconRaw', withId, (req, env, ctx) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    const bucket = env.MY_BUCKET;
    return gearService.getIconRaw(id, bucket);
  })
  .get('/:id/iconRaw/origin', withId, (req) => {
    const { id } = req;
    if (id === undefined) {
      return error(404);
    }
    return gearService.getIconOrigin(id);
  });

export default router;
