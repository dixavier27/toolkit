import { expect, test } from "bun:test";
import { greet } from "./index.ts";

test("greet em português por default", () => {
  expect(greet("mundo")).toBe("Olá, mundo!");
});

test("greet em inglês com locale", () => {
  expect(greet("world", { locale: "en" })).toBe("Hello, world!");
});
