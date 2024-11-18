# url-shortener

A simple URL Shortener built with go.

## Usage

Start the server running `go run ./main.go` or use the docker container.

Make a request with the desired URL to short. Example using `google.com`:

```sh
curl -X POST http://localhost:8008/v1/urls -d "{\"url\":\"https://google.com\"}" | json_pp
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

Open a browser in the following url:

`localhost:8008/lBYEK5Cz`

You will be redirected to `google.com`!
