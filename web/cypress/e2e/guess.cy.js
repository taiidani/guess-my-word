/// <reference types="Cypress" />

const today = new Date(2022, 11, 8, 0, 0, 0)
const tomorrow = new Date(2022, 11, 9, 0, 0, 0)
const defaultToday = "rough"
const defaultYesterday = "alive"
const hardToday = "oatcake"
const hardYesterday = "nonjoiner"

describe('guess spec', () => {
  it('guesses the default word', () => {
    cy.clock(today)
    cy.visit('/')

    cy.contains('Guess')
    cy.get('footer .col:nth-child(1)').contains(defaultYesterday)
    cy.get('.before li:nth-child(1)').contains('No guesses before the word')
    cy.get('.after').contains('No guesses after the word')

    // First guess
    cy.get('form#guesser input').type('medium')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.after').contains('No guesses after the word')

    // Second guess
    cy.get('form#guesser input').type('tame')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.after li:nth-child(1)').contains('tame')

    // Third guess
    cy.get('form#guesser input').type('right')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.before li:nth-child(2)').contains('right')
    cy.get('.after li:nth-child(1)').contains('tame')

    // Fourth guess
    cy.get('form#guesser input').type('round')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.before li:nth-child(2)').contains('right')
    cy.get('.after li:nth-child(1)').contains('round')
    cy.get('.after li:nth-child(2)').contains('tame')

    // Correct guess
    cy.get('form#guesser input').type(defaultToday)
    cy.get('form#guesser').submit()
    cy.get('#app').contains('You guessed "' + defaultToday + '" correctly')
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.after li:nth-child(1)').contains('round')
    cy.get('.after li:nth-child(2)').contains('tame')
    cy.get('.after').contains('round')
  })

  it('guesses the hard word', () => {
    cy.clock(today)
    cy.visit('/')
    cy.contains('Guess')

    // Switch to hard
    cy.get('#mode').select('Hard')
    cy.get('footer .col:nth-child(1)').contains(hardYesterday)
    cy.get('.before li:nth-child(1)').contains('No guesses before the word')
    cy.get('.after').contains('No guesses after the word')

    // First guess
    cy.get('form#guesser input').type('medium')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.after').contains('No guesses after the word')

    // Second guess
    cy.get('form#guesser input').type('tame')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.after li:nth-child(1)').contains('tame')

    // Third guess
    cy.get('form#guesser input').type('oat')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.before li:nth-child(2)').contains('oat')
    cy.get('.after li:nth-child(1)').contains('tame')

    // Fourth guess
    cy.get('form#guesser input').type('out')
    cy.get('form#guesser').submit()
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.before li:nth-child(2)').contains('oat')
    cy.get('.after li:nth-child(1)').contains('out')
    cy.get('.after li:nth-child(2)').contains('tame')

    // Correct guess
    cy.get('form#guesser input').type(hardToday)
    cy.get('form#guesser').submit()
    cy.get('#app').contains('You guessed "' + hardToday + '" correctly')
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.before li:nth-child(2)').contains('oat')
    cy.get('.after li:nth-child(1)').contains('out')
    cy.get('.after li:nth-child(2)').contains('tame')
  })
})
