import { env } from "./env.ts";
import { buildServer } from "./server.ts";

declare const __VERSION__: string;

async function bootstrap() {
  const app = await buildServer();

  try {
    await app.listen({ port: env.PORT, host: env.HOST });
    app.log.info(
      `🚀 {{name}} v${__VERSION__} rodando em http://${env.HOST}:${env.PORT}`,
    );
  } catch (err) {
    app.log.error(err);
    process.exit(1);
  }
}

bootstrap();
