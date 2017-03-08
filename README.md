# `gaze-web` - a web application for serving `gaze` records

This web app can capture the JSON payloads produced by `gaze` (see [here](https://github.com/AstromechZA/gaze)) and display
them while allowing nice paginated and searchable lists of the results.

## Installation

The archive for your platform, found on the releases page, contains a file structure like:

```
gaze-web
gaze-web/gaze-web
gaze-web/static/..
gaze-web/templates/..
```

You can deploy it anywhere on your host.

## Database

It will store and use a sqlite database `gaze-web.db` in the current directory from which it is run.
