import { expect, test } from "@playwright/test";

test("renders the home page with its metadata", async ({ page }) => {
  await page.goto("/");

  await expect(page).toHaveTitle("Threadly");
  await expect(page.getByRole("heading", { name: "Threadly" })).toBeVisible();
});

test("serves the requested robots policy", async ({ request }) => {
  const response = await request.get("/robots.txt");

  expect(response.ok()).toBeTruthy();
  expect(await response.text()).toBe(
    "Sitemap: https://precedent.dev/sitemap.xml\n\nUser-agent: *\nDisallow: /api\n",
  );
});
