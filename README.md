# localgoat

## Example Config

```json
{
    "address": "127.0.0.1",
    "port": 8080,
    "proxies": {
        "default": {
            "host": "http://localhost:1234/"
        }
    },
    "routes": [
        {
            "prefix": "/editor/",
            "static": {
                "path": "editor/public/index.htm"
            }
        },
        {
            "prefix": "/static/",
            "static": {
                "path": "editor/public/static",
                "stripPrefix": true
            }
        },
        {
            "prefix": "/",
            "proxy": {
                "target": "default"
            }
        }
    ]
}
```

## TODO

  - new CLI arg parser
  - new CLI invocation scheme
    - `lg .` - static server only
    - `lg /app/ http://localhost:4000` - reverse proxy only
    - `lg . /app/ http://localhost:4000` - static server and reverse proxy
