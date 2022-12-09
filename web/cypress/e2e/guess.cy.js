/// <reference types="Cypress" />

describe('guess spec', () => {
  it('guesses the default word', () => {
    cy.visit('/', {
      onBeforeLoad: (contentWindow) => {
        contentWindow.sessionStorage.setItem('state-default', JSON.stringify({
          "version": 0.9,
          "before": [],
          "after": [],
          "answer": "",
          "guesses": 0,
          "idleTime": 0,
          "start": "2022-12-07T17:05:37.121Z",
          "end": null,
        }))
      },
    })

    cy.contains('Guess')
    cy.get('footer .col:nth-child(1)').contains("alive")
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
    cy.get('form#guesser input').type('rough')
    cy.get('form#guesser').submit()
    cy.get('#app').contains('You guessed "rough" correctly')
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.after li:nth-child(1)').contains('round')
    cy.get('.after li:nth-child(2)').contains('tame')
    cy.get('.after').contains('round')
  })

  it('guesses the hard word', () => {
    cy.visit('/', {
      onBeforeLoad: (contentWindow) => {
        contentWindow.sessionStorage.setItem('state-hard', JSON.stringify({
          "version": 0.9,
          "before": [],
          "after": [],
          "answer": "",
          "guesses": 0,
          "idleTime": 0,
          "start": "2022-12-07T17:05:37.121Z",
          "end": null,
        }))
      },
    })
    cy.contains('Guess')

    // Switch to hard
    cy.get('#mode').select('Hard')
    cy.get('footer .col:nth-child(1)').contains("nonjoiner")
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
    cy.get('form#guesser input').type('oatcake')
    cy.get('form#guesser').submit()
    cy.get('#app').contains('You guessed "oatcake" correctly')
    cy.get('.before li:nth-child(1)').contains('medium')
    cy.get('.before li:nth-child(2)').contains('oat')
    cy.get('.after li:nth-child(1)').contains('out')
    cy.get('.after li:nth-child(2)').contains('tame')
  })
})
