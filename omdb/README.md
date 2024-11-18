# omdb

A simple api crawler for OMDB API.

## Usage

Create an API key on [OMDB API](https://www.omdbapi.com/apikey.aspx) and put it
on `.env` file.

Start the server running `go run ./main.go` or use docker.

Make a request with the desired movie name. For example:

```sh
curl http://localhost:8008/v1/movies\?s\=Blade%20Runner | json_pp
```
