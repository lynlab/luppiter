# LYnLab Luppiter GraphQL

> GraphQL APIs for LYnLab Luppiter services.  
> For further information, see documentations on [LYnLab Luppiter console](https://luppiter.lynlab.co.kr/web).

## Specifications

### URL Paths

  - `/apis/graphql` : GET or POST requests for GraphQL APIs.
  - `/files/{bucketName}/{fileName}` : GET, POST or DELETE for file access. See [storage service](https://luppiter.lynlab.co.kr/web/services/storage) for more information.

## Development

### Prerequisites

  - Go 1.11
  - PostgreSQL

### Environment Variables

```sh
cp .envrc.example .envrc
vi .envrc
```

See `.envrc.example` for list of variables. You can use [direnv](http://direnv.net) to set these values.
