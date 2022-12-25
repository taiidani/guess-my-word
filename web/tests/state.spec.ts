import { test, expect } from '@playwright/test';
import { padEnd } from 'core-js/core/string';

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

const sess = {
    "version": 0.9,
    "before": ["apple"],
    "after": [],
    "answer": "",
    "guesses": 0,
    "idleTime": 0,
    "start": today,
    "end": null,
}

test('reads current state', async ({ page }) => {
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

    await page.context().addInitScript(value => {
        window.sessionStorage.setItem('state-default', JSON.stringify(value));
    }, sess);

    await page.goto('/');
    await expect(page).toHaveTitle(/Guess My Word/);

    await expect(page.locator('footer .col:nth-child(1)')).toContainText(defaultYesterday)
    await expect(page.locator('.before li:nth-child(1)')).toContainText('apple')
    await expect(page.locator('.after')).toContainText('No guesses after the word')
})

test('ignores prior state', async ({ page }) => {
    await page.addInitScript(`{
        // Extend Date constructor to default to tomorrow
        Date = class extends Date {
          constructor(...args) {
            if (args.length === 0) {
              super(${tomorrow});
            } else {
              super(...args);
            }
          }
        }
        // Override Date.now() to start from tomorrow
        const __DateNowOffset = ${tomorrow} - Date.now();
        const __DateNow = Date.now;
        Date.now = () => __DateNow() + __DateNowOffset;
      }`);

    await page.context().addInitScript(value => {
        window.sessionStorage.setItem('state-default', JSON.stringify(value));
    }, sess);

    await page.goto('/');
    await expect(page).toHaveTitle(/Guess My Word/);

    await expect(page.locator('footer .col:nth-child(1)')).toContainText(defaultToday)
    await expect(page.locator('.before li:nth-child(1)')).toContainText('No guesses before the word')
    await expect(page.locator('.after')).toContainText('No guesses after the word')
})
