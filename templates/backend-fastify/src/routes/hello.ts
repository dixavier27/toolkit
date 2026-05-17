import type { FastifyPluginAsyncZod } from "fastify-type-provider-zod";
import { z } from "zod";

export const helloRoutes: FastifyPluginAsyncZod = async (app) => {
  app.get(
    "/:name",
    {
      schema: {
        params: z.object({
          name: z.string().min(1).max(64),
        }),
        response: {
          200: z.object({
            greeting: z.string(),
          }),
        },
      },
    },
    async ({ params }) => {
      return { greeting: `Olá, ${params.name}!` };
    },
  );
};
