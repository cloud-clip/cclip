# Cloud Clip

> Server and command line interface, which provides and accesses self-hosted clipboards in the cloud, written in [Go](https://golang.org/).

## Build

```bash
# clone repo
git clone https://github.com/cloud-clip/cclip

# change into repository
cd cclip

# build
go build .
```

## Usage

### Server

```bash
# run server on port 50979
cclip
```

Environment variables:

| Name | Description | Example |
|------|-------------|----------|
| `CCLIP_DIR` | The directory where all clips should be / are stored. Default: `./clips` | `/var/cclip/clips` |
| `CCLIP_MAX_SIZE` | The maximum size of a clip, in bytes. Default: `134217728` | `0` (unlimited) |
| `CCLIP_PASSWORD` | The password to use for all API calls. Default: none | `MySecretP@ssword123!` |
| `CCLIP_PORT` | The TCP port, the server should run on. Default: `50979` | `23979` |

### Docker

#### Build and run

```bash
# build image
docker build . -t cclip/cclip --no-cache

# mount docker directory '/app/clips' to local '/path/to/local/clips/dir'
# expose TCP port 59079
# run with password 'test' (CCLIP_PASSWORD)
docker run -v /path/to/local/clips/dir:/app/clips -p 50979:50979 -e CCLIP_PASSWORD=test -it cclip/cclip
```

#### Pull and run

```bash
# pull from Docker Hub
docker pull cclip/cclip

# run
docker run -v /path/to/local/clips/dir:/app/clips -p 50979:50979 -e CCLIP_PASSWORD=test -it cclip/cclip
```

### API

#### [GET] /api/v1

Returns (status) information about the server.

Request:

```http
GET http://localhost:50979/api/v1
Authorization: Bearer <YOUR-PASSWORD-HERE>

```

Response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 05 Sep 1979 21:09:00 GMT
Content-Length: 59
Connection: close

{
  "ip": "127.0.0.1:60000",
  "time": "2020-09-05T23:09:00+02:00"
}
```

#### [GET] /api/v1/clips

Returns the current list of clips.

Request:

```http
GET http://localhost:50979/api/v1/clips
Authorization: Bearer <YOUR-PASSWORD-HERE>

```

Response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 05 Sep 1979 21:09:00 GMT
Content-Length: 352
Connection: close

[
  {
    "id": "01234567890123456789012345678901",
    "name": "A HTML file",
    "mime": "text/html",
    "ctime": 1596200000,
    "mtime": 1596200000,
    "size": 23979,
    "resource": "/api/v1/clips/01234567890123456789012345678901",
    "share": "/api/v1/clips/01234567890123456789012345678901/share"
  },
  {
    "id": "01234567890123456789012345678902",
    "name": "A text file",
    "mime": "text/plain",
    "ctime": 1596200001,
    "mtime": 1596200001,
    "size": 5979,
    "resource": "/api/v1/clips/01234567890123456789012345678902",
    "share": "/api/v1/clips/01234567890123456789012345678902/share"
  }
]
```

#### [HEAD] /api/v1/clips

Returns short information about the current clip list.

Request:

```http
HEAD http://localhost:50979/api/v1/clips
Authorization: Bearer <YOUR-PASSWORD-HERE>

```

Response:

```http
HTTP/1.1 204 OK
Date: Wed, 05 Sep 1979 21:09:00 GMT
Content-Length: 0
X-Cclip-Count: 2
Connection: close

```

| Header | Description |
|------|-------------|
| `Date` | The timestamp of the newest clip. |
| `X-Cclip-Count` | The total number of clips. |

#### [DELETE] /api/v1/clips

Deltes ALL clips.

Request:

```http
DELETE http://localhost:50979/api/v1/clips
Authorization: Bearer <YOUR-PASSWORD-HERE>

```

Response:

```http
HTTP/1.1 204 OK
Date: Wed, 05 Sep 1979 21:09:00 GMT
Content-Length: 0
Connection: close

```

#### [GET] /api/v1/clips/{id}

Gets the data of a clip.

Request:

```http
GET http://localhost:50979/api/v1/clips/01234567890123456789012345678901
Authorization: Bearer <YOUR-PASSWORD-HERE>

```

Response:

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Wed, 05 Sep 1979 21:09:00 GMT
Content-Length: 1238
Connection: close

Gallia est omnis divisa in partes tres quarum unam incolunt Belgae aliam Aquitani tertiam qui ipsorum lingua Celtae nostra Galli appellantur
Hi omnes lingua institutis legibus inter se differunt
Gallos ab Aquitanis Garunna flumen a Belgis Matrona et Sequana dividit
Horum omnium fortissimi sunt Belgae propterea quod a cultu atque humanitate provinciae longissime absunt minimeque ad eos mercatores saepe commeant atque ea quae ad effeminandos animos pertinent important proximique sunt Germanis qui trans Rhenum incolunt quibuscum continenter bellum gerunt
Qua de causa Helvetii quoque reliquos Gallos virtute praecedunt quod fere cotidianis proeliis cum Germanis contendunt cum aut suis finibus eos prohibent aut ipsi in eorum finibus bellum gerunt
Eorum una pars quam Gallos obtinere dictum est initium capit a flumine Rhodano continetur Garunna flumine Oceano finibus Belgarum attingit etiam ab Sequanis et Helvetiis flumen Rhenum vergit ad septentriones
Belgae ab extremis Galliae finibus oriuntur pertinent ad inferiorem partem fluminis Rheni spectant in septentrionem et orientem solem
Aquitania a Garunna flumine ad Pyrenaeos montes et eam partem Oceani quae est ad Hispaniam pertinet spectat inter occasum solis et septentriones
```

#### [DELETE] /api/v1/clips/{id}

Deltes a clip.

Request:

```http
DELETE http://localhost:50979/api/v1/clips/01234567890123456789012345678901
Authorization: Bearer <YOUR-PASSWORD-HERE>

```

Response:

```http
HTTP/1.1 204 OK
Date: Wed, 05 Sep 1979 21:09:00 GMT
Content-Length: 0
Connection: close

```

#### [POST] /api/v1/clips

Uploads the data for a new clip.

Request:

```http
POST http://localhost:50979/api/v1/clips
Authorization: Bearer <YOUR-PASSWORD-HERE>
Content-Type: text/plain; charset=utf-8

Gallia est omnis divisa in partes tres quarum unam incolunt Belgae aliam Aquitani tertiam qui ipsorum lingua Celtae nostra Galli appellantur
Hi omnes lingua institutis legibus inter se differunt
Gallos ab Aquitanis Garunna flumen a Belgis Matrona et Sequana dividit
Horum omnium fortissimi sunt Belgae propterea quod a cultu atque humanitate provinciae longissime absunt minimeque ad eos mercatores saepe commeant atque ea quae ad effeminandos animos pertinent important proximique sunt Germanis qui trans Rhenum incolunt quibuscum continenter bellum gerunt
Qua de causa Helvetii quoque reliquos Gallos virtute praecedunt quod fere cotidianis proeliis cum Germanis contendunt cum aut suis finibus eos prohibent aut ipsi in eorum finibus bellum gerunt
Eorum una pars quam Gallos obtinere dictum est initium capit a flumine Rhodano continetur Garunna flumine Oceano finibus Belgarum attingit etiam ab Sequanis et Helvetiis flumen Rhenum vergit ad septentriones
Belgae ab extremis Galliae finibus oriuntur pertinent ad inferiorem partem fluminis Rheni spectant in septentrionem et orientem solem
Aquitania a Garunna flumine ad Pyrenaeos montes et eam partem Oceani quae est ad Hispaniam pertinet spectat inter occasum solis et septentriones
```

Response:

```http
HTTP/1.1 201 OK
Content-Type: application/json; charset=utf-8
Date: Wed, 05 Sep 1979 21:09:00 GMT
Content-Length: 352
Connection: close

{
  "id": "01234567890123456789012345678901",
  "name": "A HTML file",
  "mime": "text/plain; charset=utf-8",
  "ctime": 1596200000,
  "mtime": 1596200000,
  "size": 23979,
  "resource": "/api/v1/clips/01234567890123456789012345678901",
  "share": "/api/v1/clips/01234567890123456789012345678901/share"
}
```
