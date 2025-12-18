# url-shortening
Simple URL Shortening Service in Go

## Project To Do List
- [] Basic shortening features (`roadmap.sh`)
- [] Endpoints Testing
- [] Endpoints Documentation
- [] Create Docker Compose
- [] Basic Deploy: CI/CD
    - Buy Domain and VPS
- [] Create Frontend
- Add new features

## Endpoints requests and response format

### Create Short URL

Create short URL using POST method.

```
POST /shorten
{
  "url": "https://www.example.com/some/long/url"
}
```

The server responds with `201 Created` and the new shortend URL.

```json
{
  "id": "1",
  "url": "https://www.example.com/some/long/url",
  "shortCode": "abc123",
  "createdAt": "2021-09-01T12:00:00Z",
  "updatedAt": "2021-09-01T12:00:00Z"
}
```

It returns `400 Bad Request` in case of URL validation errors.


### Retrieve Original URL

Retrieve the original URL from a short URL using the `GET` method.

```
GET /shorten/abc123
```

The endpoint should return a `200 OK` status code and redirect to original URL.

If not found (short code doesn't exist), it will return `404 Not Found`.


