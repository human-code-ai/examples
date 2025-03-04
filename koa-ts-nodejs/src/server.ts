import * as Router from 'koa-router';
import ApiController from './controllers/api.controller';

const router = new Router();

router.get('/', ApiController.getIndex);
router.get('/session_id', ApiController.getSessionId);
router.get('/registration_url', ApiController.genRegistrationUrl);
router.get('/verification_url', ApiController.genVerificationUrl);
router.get('/verify', ApiController.verify);

export default router;
