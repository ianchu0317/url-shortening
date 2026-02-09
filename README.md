# url-shortening
Simple URL Shortening Service in Go

The project is deployed on https://nanolinq.ianchenn.com 


## Getting Started

To run and deploy the server run the following commands in order

```bash
# Clone the repo
git clone https://github.com/ianchu0317/url-shortening.git
# Enter the folder
cd url-shortening
# Deploy server
docker compose up -d
```
Then the services will be open in:
- `8080` -> API / Shortener Backend
- `80`   -> HTTP / Frontend Server 
- `5432` -> DB 


## Project To Do List
- [X] Basic shortening features (`roadmap.sh`)
- [X] Endpoints Testing
- [X] Endpoints Documentation
- [X] Create Docker Compose
- [X] Basic Deploy: CI/CD
    - Buy Domain and VPS
- [X] Create Frontend
- Add new features

## Endpoints requests and response format

<details>

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
GET /abc123
```

The endpoint should return a `200 OK` status code and redirect to original URL.

If not found (short code doesn't exist), it will return `404 Not Found`.


### Update Short URL

Update an existing short URL using the `PUT` method

```
PUT /abc123
{
  "url": "https://www.example.com/some/updated/url"
}
```

The endpoint validate the request body and return a `200 OK` status code with the updated short URL i.e.

```json
{
    "id": 1,
    "url": "https://www.example.com/some/updated/url",
    "shortCode": "813a43f95e",
    "createdAt": "2025-12-18T13:00:54.837205Z",
    "updatedAt": "2025-12-18T13:01:37.485254Z",
    "accessed": 3
}
```

It returns `400 Bad Request` if have bad requests or `404 Not Found` if status code not in server.


### Delete Short URL

Delete an existing short URL using the `DELETE` method.

```
DELETE /abc123
```

The endpoint should return a `204 No Content` status code if the short URL was successfully deleted or a `404 Not Found` status code if the short URL was not found.


### Get URL Statistics

Get statistics for a short URL using the `GET` method

```
GET /abc123/stats
```

The endpoint should return a 200 OK status code with the statistics i.e.

```json
{
  "id": "1",
  "url": "https://www.example.com/some/long/url",
  "shortCode": "abc123",
  "createdAt": "2021-09-01T12:00:00Z",
  "updatedAt": "2021-09-01T12:00:00Z",
  "accessCount": 10
}
```

or a `404 Not Found` status code if the short URL was not found.

</details>