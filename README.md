# `gaze-web` - a web application for serving `gaze` records

This web app can capture the JSON payloads produced by `gaze` (see [here](https://github.com/AstromechZA/gaze)) and display
them while allowing nice paginated and searchable lists of the results.

|![Landing](https://i.imgur.com/NZU6IRt.jpg)|![Reports](https://i.imgur.com/gtwlRNR.jpg)|![Report Display](https://i.imgur.com/U5XKzg4.jpg)|
|:---:|:---:|:---:|
|![Report Filter](https://i.imgur.com/oc5wRic.jpg)|![Reports Graph](https://i.imgur.com/9G5ZsOU.jpg)||

## Installation

The archive for your platform, found on the releases page, contains a file structure like:

```
gaze-web/gaze-web (binary)
gaze-web/static/..
gaze-web/templates/..
```

You can deploy this directory anywhere on your host and link to the binary to run it.

## Database

All reports received by the server are stored in a BoltDB file on disk in the project directory.
