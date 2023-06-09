import { IRequest, Router, error } from 'itty-router';
import gearRouter from './gear/router';
import itemRouter from './item/router';
import { preflight } from './middlewares/cors';
import { CF } from './types';

const router = Router<IRequest, CF>();
router
  .all('*', preflight)
  .all('/gears/*', gearRouter.handle)
  .all('/items/*', itemRouter.handle)
  .all('*', () => error(404));

export default router;
