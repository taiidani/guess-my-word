import { test, expect } from '@playwright/test'

const today = new Date(Date.UTC(2022, 11, 8, 0, 0, 0)).getTime()
const tomorrow = new Date(Date.UTC(2022, 11, 9, 0, 0, 0)).getTime()
const defaultToday = "rough"
const defaultYesterday = "alive"
const hardToday = "oatcake"
const hardYesterday = "nonjoiner"

test.use({
    locale: 'en-US',
    timezoneId: 'America/Los_Angeles',
})

test('guesses the default word', async ({ page }) => {
    await page.addInitScript(`{
        // Extend Date constructor to default to today
        Date = class extends Date {
          constructor(...args) {
            if (args.length === 0) {
              super(${today});
            } else {
              super(...args);
            }
          }
        }
        // Override Date.now() to start from today
        const __DateNowOffset = ${today} - Date.now();
        const __DateNow = Date.now;
        Date.now = () => __DateNow() + __DateNowOffset;
      }`)

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
    await lWordEntry.type('medium')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('medium')
    await expect(lGuessAfter).toHaveText('No guesses after the word')

    // Second guess
    await lWordEntry.type('tame')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('medium')
    await expect(lGuessAfter).toHaveText('tame')

    // Third guess
    await lWordEntry.type('right')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('mediumright')
    await expect(lGuessAfter).toHaveText('tame')

    // Fourth guess
    await lWordEntry.type('round')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('mediumright')
    await expect(lGuessAfter).toHaveText('roundtame')

    // Correct guess
    await lWordEntry.type(defaultToday)
    await lWordEntry.press('Enter')
    await expect(page.locator('#app')).toContainText('You guessed "' + defaultToday + '" correctly')
    await expect(lGuessBefore).toHaveText('mediumright')
    await expect(lGuessAfter).toHaveText('roundtame')
})

test('guesses the hard word', async ({ page }) => {
    await page.addInitScript(`{
        // Extend Date constructor to default to today
        Date = class extends Date {
          constructor(...args) {
            if (args.length === 0) {
              super(${today});
            } else {
              super(...args);
            }
          }
        }
        // Override Date.now() to start from today
        const __DateNowOffset = ${today} - Date.now();
        const __DateNow = Date.now;
        Date.now = () => __DateNow() + __DateNowOffset;
      }`)

    await page.goto('/')
    await expect(page).toHaveTitle(/Guess My Word/)

    const lWordEntry = page.getByPlaceholder('Enter a word here')
    const lStatsYesterday = page.locator("#stats-yesterday")
    const lGuessBefore = page.locator(".before")
    const lGuessAfter = page.locator(".after")

    // Switch to hard
    await page.locator('#mode').selectOption('Hard')
    await expect(lStatsYesterday).toContainText(hardYesterday)
    await expect(lGuessBefore).toHaveText('No guesses before the word')
    await expect(lGuessAfter).toHaveText('No guesses after the word')

    // First guess
    await lWordEntry.fill('medium')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('medium')
    await expect(lGuessAfter).toHaveText('No guesses after the word')

    // Second guess
    await lWordEntry.fill('tame')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('medium')
    await expect(lGuessAfter).toHaveText('tame')

    // Third guess
    await lWordEntry.fill('oat')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('mediumoat')
    await expect(lGuessAfter).toHaveText('tame')

    // Fourth guess
    await lWordEntry.fill('out')
    await lWordEntry.press('Enter')
    await expect(lGuessBefore).toHaveText('mediumoat')
    await expect(lGuessAfter).toHaveText('outtame')

    // Correct guess
    await lWordEntry.type(hardToday)
    await lWordEntry.press('Enter')
    await expect(page.locator('#app')).toContainText('You guessed "' + hardToday + '" correctly')
    await expect(lGuessBefore).toHaveText('mediumoat')
    await expect(lGuessAfter).toHaveText('outtame')
})
