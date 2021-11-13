# Guess My Word

![Go](https://github.com/taiidani/guess-my-word/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/taiidani/guess-my-word/branch/master/graph/badge.svg)](https://codecov.io/gh/taiidani/guess-my-word)

This is a pet project game around guessing a word based on its location in the alphabet. Users will provide a word, and the application will tell them if that word falls before or after the target in the alphabet. This process will then be repeated until the target is found. HINT: Think binary search trees!

Currently a work in progress. Feature requests are welcome using GitHub issues!

## Attributions

Inspiration was taken heavily from https://hryanjones.com/guess-my-word/.

The Scrabble word list was obtained from https://sourceforge.net/projects/scrabbledict/.

The backend web framework in use is Gin: https://gin-gonic.com

The frontend frameworks in use are:

* Vue.js: https://vuejs.org/
* Bootstrap: https://getbootstrap.com/
* jQuery: https://jquery.com/

Persistence is being enabled using a Redis backend.

## Contributing

The application requires the following to be configured:

* Go 1.17+ installed

Start the application in development mode with:

```
make && ./bin/guess-my-word
```

Then view the website at http://127.0.0.1:3000.

Dev away!

### Persistence

By default the application runs in "Local Mode" and will not persist any of the data (such as generated words). If you need to test the persistence options you may point the application at a local Redis instance:

```sh
docker-compose up -d redis
make && REDIS_URL=127.0.0.1:6379 ./bin/guess-my-word
```
