import { IRequest, Router, error } from 'itty-router';
import { withId } from '../middlewares/id';
import { CF } from '../types';
import * as gearService from './service';

const router = Router<IRequest, CF>({ base: '/items' }).get('/:id/iconRaw/origin', withId, (req) => {
  const { id } = req;
  if (id === undefined) {
    return error(404);
  }
  return gearService.getIconRawOrigin(id);
});

export default router;
