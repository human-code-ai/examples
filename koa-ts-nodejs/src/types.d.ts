import HumanCodeProvider from "./provider/humancode.provider";

declare module 'koa' {
  interface Context {
    state: {
      humanCode: HumanCodeProvider;
    };
  }
}