import { afterAll, beforeAll, expect, test } from "bun:test";
import type { FastifyInstance } from "fastify";
import { buildServer } from "../src/server.ts";

let app: FastifyInstance;

beforeAll(async () => {
  app = await buildServer();
  await app.ready();
});

afterAll(async () => {
  await app.close();
});

test("GET /health responde 200 com status ok", async () => {
  const response = await app.inject({ method: "GET", url: "/health" });
  expect(response.statusCode).toBe(200);
  expect(response.json()).toMatchObject({ status: "ok" });
});

test("GET /hello/:name responde com saudação", async () => {
  const response = await app.inject({ method: "GET", url: "/hello/mundo" });
  expect(response.statusCode).toBe(200);
  expect(response.json()).toEqual({ greeting: "Olá, mundo!" });
});

test("GET /hello/:name valida tamanho do nome", async () => {
  const response = await app.inject({ method: "GET", url: "/hello/" });
  expect(response.statusCode).toBe(404);
});
