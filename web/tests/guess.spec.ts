import { test, expect } from '@playwright/test'

const defaultToday = "course"
const defaultYesterday = "worst"
const hardToday = "glissader"
const hardYesterday = "gabbroid"

test.use({
  locale: 'en-US',
  timezoneId: 'America/Los_Angeles',
})

test('guesses the default word', async ({ page }) => {
  await page.goto('/api/seed');
  await page.waitForLoadState("load");

  await page.goto('/')
  await expect(page).toHaveTitle(/Guess My Word/)

  const lWordEntry = page.getByPlaceholder('Enter a word here')
  const lStatsYesterday = page.locator("#stats-yesterday")
  const lGuessBefore = page.locator(".before")
  const lGuessAfter = page.locator(".after")

  await expect(lStatsYesterday).toContainText(defaultYesterday)
  await expect(lGuessBefore).toHaveText('No guesses before the word')
  await expect(lGuessAfter).toHaveText('No guesses after the word')

  // First guess
  await lWordEntry.type('apple')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple')
  await expect(lGuessAfter).toHaveText('No guesses after the word')

  // Second guess
  await lWordEntry.type('jeans')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple')
  await expect(lGuessAfter).toHaveText('jeans')

  // Third guess
  await lWordEntry.type('bottom')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple bottom')
  await expect(lGuessAfter).toHaveText('jeans')

  // Fourth guess
  await lWordEntry.type('hey')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple bottom')
  await expect(lGuessAfter).toHaveText('hey jeans')

  // Correct guess
  await lWordEntry.type(defaultToday)
  await lWordEntry.press('Enter')
  await expect(page.locator('#app')).toContainText('You guessed "' + defaultToday + '" correctly')
  await expect(lGuessBefore).toHaveText('apple bottom')
  await expect(lGuessAfter).toHaveText('hey jeans')
})

test('guesses the hard word', async ({ page }) => {
  await page.goto('/api/seed');
  await page.waitForLoadState("load");

  await page.goto('/?mode=hard')
  await expect(page).toHaveTitle(/Guess My Word/)

  const lWordEntry = page.getByPlaceholder('Enter a word here')
  const lStatsYesterday = page.locator("#stats-yesterday")
  const lGuessBefore = page.locator(".before")
  const lGuessAfter = page.locator(".after")

  await expect(lStatsYesterday).toContainText(hardYesterday)
  await expect(lGuessBefore).toHaveText('No guesses before the word')
  await expect(lGuessAfter).toHaveText('No guesses after the word')

  // First guess
  await lWordEntry.fill('apple')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple')
  await expect(lGuessAfter).toHaveText('No guesses after the word')

  // Second guess
  await lWordEntry.fill('tree')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple')
  await expect(lGuessAfter).toHaveText('tree')

  // Third guess
  await lWordEntry.fill('cherry')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple cherry')
  await expect(lGuessAfter).toHaveText('tree')

  // Fourth guess
  await lWordEntry.fill('trunk')
  await lWordEntry.press('Enter')
  await expect(lGuessBefore).toHaveText('apple cherry')
  await expect(lGuessAfter).toHaveText('tree trunk')

  // Correct guess
  await lWordEntry.type(hardToday)
  await lWordEntry.press('Enter')
  await expect(page.locator('#app')).toContainText('You guessed "' + hardToday + '" correctly')
  await expect(lGuessBefore).toHaveText('apple cherry')
  await expect(lGuessAfter).toHaveText('tree trunk')
})
