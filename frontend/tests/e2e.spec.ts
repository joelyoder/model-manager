import { test, expect } from "@playwright/test";

// Mock external API calls to keep tests deterministic
test.beforeEach(async ({ page }) => {
  await page.route("**/civitai.com/**", (route) => {
    return route.fulfill({ status: 200, body: "{}" });
  });
});

test("navigate between routes and SPA fallback", async ({ page }) => {
  await page.goto("/");
  await page.click('a[href="#/utilities"]');
  await expect(page.locator("h2")).toHaveText("Utilities");

  await page.click('a[href="#/settings"]');
  await expect(page.locator("h2")).toHaveText("Settings");

  await page.click('a[href="#/"]');
  await expect(
    page.locator('input[placeholder="Search models..."]'),
  ).toBeVisible();

  await page.goto("/non-existent");
  await expect(
    page.locator('input[placeholder="Search models..."]'),
  ).toBeVisible();
});

test("add model via UI and verify in list", async ({ page }) => {
  await page.goto("/");
  await page.click("text=Add Models");
  await page.click("text=Add Model");

  await page
    .locator('label:has-text("Name") + input')
    .first()
    .fill("Test Model");
  await page.locator('label:has-text("Version Name") + input').fill("v1");
  await page.click("text=Save");
  await page.click("text=Back");

  await expect(
    page.locator(".card").filter({ hasText: "Test Model" }),
  ).toBeVisible();
});

test("update settings and ensure persistence", async ({ page }) => {
  await page.goto("/");
  await page.click('a[href="#/settings"]');

  const apiKeyInput = page.locator('input[type="text"]');
  await apiKeyInput.fill("my-test-key");
  await page.click("text=Save");

  await page.click('a[href="#/"]');
  await page.click('a[href="#/settings"]');

  await expect(page.locator('input[type="text"]')).toHaveValue("my-test-key");
});

test("find duplicate file paths shows message when none found", async ({
  page,
}) => {
  await page.goto("/");
  await page.click('a[href="#/utilities"]');
  await page.click("text=Search Duplicates");
  await expect(
    page.locator("text=No duplicate file paths found"),
  ).toBeVisible();
});
