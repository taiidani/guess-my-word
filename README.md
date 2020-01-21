# Guess My Word

This is a pet project based on original work done at https://hryanjones.com/guess-my-word/, after addicting me for a couple of weeks. It is a variant that pushes more of the logic to the backend, preventing any JavaScript-savvy individuals from spoiling themselves of the fun of guessing :)

Currently a work in progress.

## Attributions

Inspiration was taken heavily from https://hryanjones.com/guess-my-word/.

The Scrabble word list was obtained from https://sourceforge.net/projects/scrabbledict/.

The web framework in use is Buffalo: https://gobuffalo.io/en

## Contributing

The application requires the following to be configured:

* Go 1.11+ installed
* The [Buffalo](https://gobuffalo.io/en) CLI installed.
* If using the leaderboard functionality, Docker and Docker Compose.

If you are going to access the leaderboard functionality, first start its Postgres database with:

```sh
docker-compose up -d

# First time use only -- run migrations
buffalo pop create -a
```

Once up and running, start the application in development mode with:

```
buffalo dev
```

Then view the website at http://127.0.0.1:3000.

Dev away!