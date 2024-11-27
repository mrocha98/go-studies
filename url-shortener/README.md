# url-shortener

A simple URL Shortener built with go.

## Usage

Start the server running `go run ./cmd/api/main.go` or use the docker container.

Make a request with the desired URL to short. Example using `google.com`:

```sh
curl -X POST http://localhost:8008/api/url/shorten -d "{\"url\":\"https://google.com\"}" | json_pp
```

The response will be something like this:

```json
{
   "data" : {
      "code" : "lBYEK5Cz"
   },
   "date" : "2024-11-18T00:18:07.965635045Z"
}
```

You can consult the url doing this:

```sh
curl http://localhost:8008/api/url/lBYEK5Cz | json_pp
```

Resulting in:

```json
{
   "data" : {
	  "url": "https://google.com",
   },
   "date" : "2024-11-18T00:18:07.965635045Z"
}
```

If you want to redirect to the url you can pass the query param `redirect=1` or `redirect=true`

```sh
# in your browser
http://localhost:8008/api/url/lBYEK5Cz?redirect=1
```
