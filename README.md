# Guess My Word

This is a pet project based on original work done at https://hryanjones.com/guess-my-word/, after addicting me for a couple of weeks. It is a variant that pushes more of the logic to the backend, preventing any JavaScript-savvy individuals from spoiling themselves of the fun of guessing :)

Currently a work in progress.

## Attributions

Inspiration was taken heavily from https://hryanjones.com/guess-my-word/.

The Scrabble word list was obtained from https://sourceforge.net/projects/scrabbledict/.

The backend web framework in use is Gin: https://gin-gonic.com

The static packaging library used for the compile is Pkger: https://github.com/markbates/pkger

The frontend frameworks in use are:

* Vue.js: https://vuejs.org/
* Bootstrap: https://getbootstrap.com/
* jQuery: https://jquery.com/

## Contributing

The application requires the following to be configured:

* Go 1.11+ installed

Start the application in development mode with:

```
make && ./bin/guess-my-word
```

Then view the website at http://127.0.0.1.

Dev away!
