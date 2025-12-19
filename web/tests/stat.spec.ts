import { test, expect } from "@playwright/test";

const defaultToday = "website";
const defaultYesterday = "worst";
const hardToday = "gemshorn";
const hardYesterday = "gabbroid";

test.use({
  locale: "en-US",
});

test("guesses the default word", async ({ page }) => {
  await page.goto("/api/seed");
  await page.waitForLoadState("load");

  await page.goto("/stats");
  await page.waitForLoadState("load");
  await expect(page).toHaveTitle(/Guess My Word/);

  const lStatsDefaultYesterday = page.locator(".list-default .yesterday");
  await expect(lStatsDefaultYesterday).toContainText(defaultYesterday);

  const lStatsHardYesterday = page.locator(".list-hard .yesterday");
  await expect(lStatsHardYesterday).toContainText(hardYesterday);
});
