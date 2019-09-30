# Supermarket API

Included is a [REST-ful](https://restfulapi.net/) API for a Supermarket written in Golang. Upon startup, the API is initialized with the following data:

```json
[
  {
      "ProduceCode": "A12T-4GH7-QPL9-3N4M",
      "Name": "Lettuce",
      "UnitPrice": "3.64"
  },
  {
      "ProduceCode": "E5T6-9UI3-TH15-QR88",
      "Name": "Peach",
      "UnitPrice": "2.99"
  },
  {
      "ProduceCode": "YRT6-72AS-K736-L4AR",
      "Name": "Green Pepper",
      "UnitPrice": "0.79"
  },
  {
      "ProduceCode": "TQ4C-VV6T-75ZX-1RMR",
      "Name": "Gala Apple",
      "UnitPrice": "3.59"
  }
]
```

Each Food requires the following data:
* ProduceCode - a 16 (kabob-cased) character string which only contains alphanumeric characters in each 4 character rune segment.
* Name - name of the food. Must be included as a string in any request to create a food.
* UnitPrice - unit price of the food, included as a stringified 2 digit precision float.

## Supported Routes

### Show
Index

`GET /food`

Example Request:
```
curl --request GET \
  --url http://localhost:9000/food
```

Example Response:

`200 - OK`
```
[
  {
      "ProduceCode": "A12T-4GH7-QPL9-3N4M",
      "Name": "Lettuce",
      "UnitPrice": "3.64"
  },
  {
      "ProduceCode": "E5T6-9UI3-TH15-QR88",
      "Name": "Peach",
      "UnitPrice": "2.99"
  },
  {
      "ProduceCode": "YRT6-72AS-K736-L4AR",
      "Name": "Green Pepper",
      "UnitPrice": "0.79"
  },
  {
      "ProduceCode": "TQ4C-VV6T-75ZX-1RMR",
      "Name": "Gala Apple",
      "UnitPrice": "3.59"
  }
]
```
Detailed Show

Example Request: 

`GET /food/{code}`

Example Request:
```
curl --request GET \
  --url http://localhost:9000/food/A12T-4GH7-QPL9-3N4M
```

Example Response:

`200 - OK`
```
[
  {
      "ProduceCode": "A12T-4GH7-QPL9-3N4M",
      "Name": "Lettuce",
      "UnitPrice": "3.64"
  },
  {
      "ProduceCode": "E5T6-9UI3-TH15-QR88",
      "Name": "Peach",
      "UnitPrice": "2.99"
  },
  {
      "ProduceCode": "YRT6-72AS-K736-L4AR",
      "Name": "Green Pepper",
      "UnitPrice": "0.79"
  },
  {
      "ProduceCode": "TQ4C-VV6T-75ZX-1RMR",
      "Name": "Gala Apple",
      "UnitPrice": "3.59"
  }
]
```

Example Response:

### Create
`POST /foods`

Example Request:

```
curl --request POST \
  --url http://localhost:9000/foods \
  --header 'content-type: application/json' \
  --data '{
	"ProduceCode": "ABCD-1234-DEFG-7890",
	"Name": "Orange",
	"UnitPrice": "1.22"
}'
```

Example Response:

`201 - Created`
```
[
  {
      "ProduceCode": "ABCD-1234-DEFG-7890",
      "Name": "Orange",
      "UnitPrice": "1.22"
  }
]
```

Example Request:
```
curl --request POST \
  --url http://localhost:9000/foods \
  --header 'content-type: application/json' \
  --data '[
	{
		"ProduceCode": "ABCD-1234-DEFG-7890",
		"Name": "Orange",
		"UnitPrice": "1.22"
	},
	{
		"ProduceCode": "DCBA-4321-GFED-9876",
		"Name": "Grapes",
		"UnitPrice": "3.45"
	}
]
'
```

Example Response:
`201 - Created`

Example Response:
```
[
  {
	"ProduceCode": "ABCD-1234-DEFG-7890",
	"Name": "Orange",
	"UnitPrice": "1.22"
  },
  {
	"ProduceCode": "DCBA-4321-GFED-9876",
	"Name": "Grapes",
	"UnitPrice": "3.45"
  }
]
```

### Delete
`DELETE /foods/{code}`

Example Request:

```
curl --request DELETE \
  --url http://localhost:9000/foods/A12T-4GH7-QPL9-3N4M
```

Example Response:

`204 - No Content`
```
No Body
```

## Travis
A Travis build configuration is also included here for the purpose of linting, testing and publishing the image to the [DockerHub repository](https://hub.docker.com/r/danbelling/supermarket-api).

## Building the Docker Image
The included Dockerfile implements a [multi-stage process](https://docs.docker.com/develop/develop-images/multistage-build/) to produce a smaller image size. This uses the build flags `CGO_ENABLED` and `GOOS` to produce a minimal binary compatible with the Docker [scratch image](https://hub.docker.com/_/scratch).

```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .
docker build -t danbelling/supermarket-api .
```