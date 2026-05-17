import { expect, test } from "@playwright/test";

test("home page exibe título e contador funcional", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { level: 1 })).toContainText("{{name}}");

  const counter = page.getByTestId("counter-btn");
  await expect(counter).toContainText("Cliques: 0");

  await counter.click();
  await expect(counter).toContainText("Cliques: 1");

  await counter.click();
  await expect(counter).toContainText("Cliques: 2");
});
