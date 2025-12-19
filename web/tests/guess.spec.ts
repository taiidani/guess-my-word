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

  await page.goto("/");
  await expect(page).toHaveTitle(/Guess My Word/);

  const lWordEntry = page.getByPlaceholder("Enter a word here");
  const lGuessBefore = page.locator(".before");
  const lGuessAfter = page.locator(".after");

  await expect(lGuessBefore).toHaveText("No guesses before the word");
  await expect(lGuessAfter).toHaveText("No guesses after the word");

  // First guess
  await lWordEntry.type("apple");
  await lWordEntry.press("Enter");
  await expect(lGuessBefore).toHaveText("keyboard_arrow_down apple");
  await expect(lGuessAfter).toHaveText("No guesses after the word");

  // Second guess
  await lWordEntry.type("yam");
  await lWordEntry.press("Enter");
  await expect(lGuessBefore).toHaveText("keyboard_arrow_down apple");
  await expect(lGuessAfter).toHaveText("keyboard_arrow_up yam");

  // Third guess
  await lWordEntry.type("ham");
  await lWordEntry.press("Enter");
  await expect(lGuessBefore).toHaveText(
    "keyboard_arrow_down apple keyboard_arrow_down ham",
  );
  await expect(lGuessAfter).toHaveText("keyboard_arrow_up yam");

  // Fourth guess
  await lWordEntry.type("zoo");
  await lWordEntry.press("Enter");
  await expect(lGuessBefore).toHaveText(
    "keyboard_arrow_down apple keyboard_arrow_down ham",
  );
  await expect(lGuessAfter).toHaveText(
    "keyboard_arrow_up yam keyboard_arrow_up zoo",
  );

  // Correct guess
  await lWordEntry.type(defaultToday);
  await lWordEntry.press("Enter");
  await expect(page.locator("#guess-form").first()).toContainText(
    'You guessed "' + defaultToday + '" correctly',
  );
  await expect(lGuessBefore).toHaveText("apple ham");
  await expect(lGuessAfter).toHaveText("yam zoo");
});

test("guesses the hard word", async ({ page }) => {
  await page.goto("/api/seed");
  await page.waitForLoadState("load");

  await page.goto("/mode/hard");
  await expect(page).toHaveTitle(/Guess My Word/);

  const lWordEntry = page.getByPlaceholder("Enter a word here");
  const lGuessBefore = page.locator(".before");
  const lGuessAfter = page.locator(".after");

  await expect(lGuessBefore).toHaveText("No guesses before the word");
  await expect(lGuessAfter).toHaveText("No guesses after the word");

  // First guess
  await lWordEntry.fill("apple");
  await lWordEntry.press("Enter");
  await expect(lGuessBefore).toHaveText("keyboard_arrow_down apple");
  await expect(lGuessAfter).toHaveText("No guesses after the word");

  // Second guess
  await lWordEntry.fill("tree");
  await lWordEntry.press("Enter");
  await page.waitForLoadState("load");
  await expect(lGuessBefore).toHaveText("keyboard_arrow_down apple");
  await expect(lGuessAfter).toHaveText("keyboard_arrow_up tree");

  // Third guess
  await lWordEntry.fill("cherry");
  await lWordEntry.press("Enter");
  await page.waitForLoadState("load");
  await expect(lGuessBefore).toHaveText(
    "keyboard_arrow_down apple keyboard_arrow_down cherry",
  );
  await expect(lGuessAfter).toHaveText("keyboard_arrow_up tree");

  // Fourth guess
  await lWordEntry.fill("trunk");
  await lWordEntry.press("Enter");
  await page.waitForLoadState("load");
  await expect(lGuessBefore).toHaveText(
    "keyboard_arrow_down apple keyboard_arrow_down cherry",
  );
  await expect(lGuessAfter).toHaveText(
    "keyboard_arrow_uptree keyboard_arrow_up trunk",
  );

  // Correct guess
  await lWordEntry.type(hardToday);
  await lWordEntry.press("Enter");
  await expect(page.locator("#guess-form").first()).toContainText(
    'You guessed "' + hardToday + '" correctly',
  );
  await expect(lGuessBefore).toHaveText("apple cherry");
  await expect(lGuessAfter).toHaveText("tree trunk");
});
