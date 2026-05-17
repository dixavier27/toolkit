import cors from "@fastify/cors";
import helmet from "@fastify/helmet";
import Fastify from "fastify";
import {
  type ZodTypeProvider,
  serializerCompiler,
  validatorCompiler,
} from "fastify-type-provider-zod";
import { env } from "./env.ts";
import { errorHandlerPlugin } from "./plugins/error-handler.ts";
import { helloRoutes } from "./routes/hello.ts";

export async function buildServer() {
  const app = Fastify({
    logger: {
      level: env.LOG_LEVEL,
      transport:
        env.NODE_ENV === "development"
          ? { target: "pino-pretty" }
          : undefined,
    },
  }).withTypeProvider<ZodTypeProvider>();

  app.setValidatorCompiler(validatorCompiler);
  app.setSerializerCompiler(serializerCompiler);

  await app.register(helmet);
  await app.register(cors, { origin: true });
  await app.register(errorHandlerPlugin);

  app.get("/health", () => ({ status: "ok", uptime: process.uptime() }));

  await app.register(helloRoutes, { prefix: "/hello" });

  return app;
}
