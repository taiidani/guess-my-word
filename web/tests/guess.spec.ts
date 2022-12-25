import { test, expect } from '@playwright/test';

const today = new Date(Date.UTC(2022, 11, 8, 0, 0, 0)).getTime()
const tomorrow = new Date(Date.UTC(2022, 11, 9, 0, 0, 0)).getTime()
const defaultToday = "rough"
const defaultYesterday = "alive"
const hardToday = "oatcake"
const hardYesterday = "nonjoiner"

test.use({
    locale: 'en-US',
    timezoneId: 'America/Los_Angeles',
});

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
      }`);

    await page.goto('/');
    await expect(page).toHaveTitle(/Guess My Word/);

    await expect(page.locator('footer .col:nth-child(1)')).toContainText(defaultYesterday)
    await expect(page.locator('.before li:nth-child(1)')).toContainText('No guesses before the word')
    await expect(page.locator('.after')).toContainText('No guesses after the word')

    // First guess
    await page.locator('form#guesser input').type('medium')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.after')).toContainText('No guesses after the word')

    // Second guess
    await page.locator('form#guesser input').type('tame')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('tame')

    // Third guess
    await page.locator('form#guesser input').type('right')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.before li:nth-child(2)')).toContainText('right')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('tame')

    // Fourth guess
    await page.locator('form#guesser input').type('round')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.before li:nth-child(2)')).toContainText('right')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('round')
    await expect(page.locator('.after li:nth-child(2)')).toContainText('tame')

    // Correct guess
    await page.locator('form#guesser input').type(defaultToday)
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('#app')).toContainText('You guessed "' + defaultToday + '" correctly')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('round')
    await expect(page.locator('.after li:nth-child(2)')).toContainText('tame')
    await expect(page.locator('.after')).toContainText('round')
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
      }`);

    await page.goto('/');
    await expect(page).toHaveTitle(/Guess My Word/);

    // Switch to hard
    await page.locator('#mode').selectOption('Hard')
    await expect(page.locator('footer .col:nth-child(1)')).toContainText(hardYesterday)
    await expect(page.locator('.before li:nth-child(1)')).toContainText('No guesses before the word')
    await expect(page.locator('.after')).toContainText('No guesses after the word')

    // First guess
    await page.locator('form#guesser input').type('medium')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.after')).toContainText('No guesses after the word')

    // Second guess
    await page.locator('form#guesser input').type('tame')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('tame')

    // Third guess
    await page.locator('form#guesser input').type('oat')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.before li:nth-child(2)')).toContainText('oat')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('tame')

    // Fourth guess
    await page.locator('form#guesser input').type('out')
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.before li:nth-child(2)')).toContainText('oat')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('out')
    await expect(page.locator('.after li:nth-child(2)')).toContainText('tame')

    // Correct guess
    await page.locator('form#guesser input').type(hardToday)
    await page.locator('form#guesser input').press('Enter')
    await expect(page.locator('#app')).toContainText('You guessed "' + hardToday + '" correctly')
    await expect(page.locator('.before li:nth-child(1)')).toContainText('medium')
    await expect(page.locator('.before li:nth-child(2)')).toContainText('oat')
    await expect(page.locator('.after li:nth-child(1)')).toContainText('out')
    await expect(page.locator('.after li:nth-child(2)')).toContainText('tame')
})
