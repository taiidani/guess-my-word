/// <reference types="Cypress" />

const today = Date.UTC(2022, 11, 8, 0, 0, 0)
const tomorrow = Date.UTC(2022, 11, 9, 0, 0, 0)
const defaultToday = "rough"
const defaultYesterday = "alive"
const hardToday = "oatcake"
const hardYesterday = "nonjoiner"

describe('guess spec', () => {
    it('throws away old states', () => {
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

        // Same day
        cy.clock(today)
        cy.visit('/', {
            onBeforeLoad: (contentWindow) => {
                contentWindow.sessionStorage.setItem('state-default', JSON.stringify(sess))
            },
        })
        cy.get('footer .col:nth-child(1)').contains(defaultYesterday)
        cy.get('.before li:nth-child(1)').contains('apple')
        cy.get('.after').contains('No guesses after the word')

        // Next day
        cy.clock(tomorrow)
        cy.visit('/', {
            onBeforeLoad: (contentWindow) => {
                contentWindow.sessionStorage.setItem('state-default', JSON.stringify(sess))
            },
        })
        cy.get('footer .col:nth-child(1)').contains(defaultToday)
        cy.get('.before li:nth-child(1)').contains('No guesses before the word')
        cy.get('.after').contains('No guesses after the word')
    })
})
