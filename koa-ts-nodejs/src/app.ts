import * as Koa from 'koa';
import * as bodyParser from 'koa-bodyparser';
import * as cors from '@koa/cors';
import koaHelmet from 'koa-helmet';
import * as json from 'koa-json';
import * as logger from 'koa-logger';
import 'reflect-metadata';
import router from './server';
import HumanCodeProvider from './provider/humancode.provider';
import 'dotenv/config';

const app = new Koa();
const port = process.env.PORT || 3000;

app.use(koaHelmet());
app.use(cors());
app.use(json());
app.use(logger());
app.use(bodyParser());

app.use(async (ctx, next) => {
  ctx.state.humanCode = new HumanCodeProvider({
    baseUrl: process.env.BASE_URL,
    debug: process.env.DEBUG === 'true',
    appId: process.env.APP_ID,
    appKey: process.env.APP_KEY
  });
  await next();
})

app.use(router.routes()).use(router.allowedMethods());

app.listen(port, () => {
  console.log(`🚀 App listening on the port ${port}`);
});
