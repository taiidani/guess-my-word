# Guess My Word

![Go](https://github.com/taiidani/guess-my-word/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/taiidani/guess-my-word/branch/master/graph/badge.svg)](https://codecov.io/gh/taiidani/guess-my-word)

This is a pet project game around guessing a word based on its location in the alphabet. Users will provide a word, and the application will tell them if that word falls before or after the target in the alphabet. This process will then be repeated until the target is found. HINT: Think binary search trees!

Currently a work in progress. Feature requests are welcome using GitHub issues!

## Attributions

Inspiration was taken heavily from https://hryanjones.com/guess-my-word/.

The Scrabble word list was obtained from https://sourceforge.net/projects/scrabbledict/.

The backend web framework in use is Gin: https://gin-gonic.com

The static packaging library used for the compile is Pkger: https://github.com/markbates/pkger

The frontend frameworks in use are:

* Vue.js: https://vuejs.org/
* Bootstrap: https://getbootstrap.com/
* jQuery: https://jquery.com/

Persistence is being enabled using the Google Cloud Platform Firestore engine: https://console.cloud.google.com/firestore

## Contributing

The application requires the following to be configured:

* Go 1.11+ installed

Start the application in development mode with:

```
make && ./bin/guess-my-word
```

Then view the website at http://127.0.0.1.

Dev away!

### Persistence

By default the application runs in "Local Mode" and will not persist any of the data (such as generated words). If you need to test the persistence options:

* Generate a [Google Cloud Platform](https://console.cloud.google.com/firestore) project using Firestore Native Mode
* Create a service account for your project, then download its JSON credentials to an "auth.json" file at the root of this directory
* Set an environment variable in your shell of `GOOGLE_APPLICATION_PROJECT_ID` to the ID of your GCP project.

This should begin generating collections in your Firestore installation upon first pageload.
