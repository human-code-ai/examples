import { randomUUID } from 'crypto';
import { Context } from 'koa';

export default class ApiController {

  public static async getIndex(ctx: Context) {
    ctx.status = 200;
  }

  public static async getSessionId(ctx: Context) {
    const result = await ctx.state.humanCode.getSessionId(randomUUID());
    ctx.body = {
      sessionId: result.session_id
    };
  }

  public static async genRegistrationUrl(ctx: Context) {
    const result = await ctx.state.humanCode.getSessionId(randomUUID());
    const registrationUrl = ctx.state.humanCode.genRegistrationUrl(result.session_id, 'http://192.168.110.24:3000/verify');
    ctx.body = {
      registrationUrl
    };
  }
  
  public static async genVerificationUrl(ctx: Context) {
    const result = await ctx.state.humanCode.getSessionId(randomUUID());
    const registrationUrl = ctx.state.humanCode.genVerificationUrl(result.session_id, 'humanid...', 'http://192.168.110.24:3000/verify');
    ctx.body = {
      registrationUrl
    };
  }
  
  public static async verify(ctx: Context) {
    const { session_id, vcode, error_code } = ctx.request.query;

    if (Number(error_code) != 0) {
      ctx.status = 400;
    }
    const result = await ctx.state.humanCode.verify(session_id as string, vcode as string, randomUUID());
    ctx.body = {
      human_id: result.human_id
    };
  }
}
