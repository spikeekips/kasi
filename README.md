# kasi

`kasi` is the transparent gateway(or bridge, or proxy) for the privately managed or public services, which are based on `HTTP`(`HTTPS`) . It can be in the middle of between your clients and services and be possible to make a single point of service.

> In Kiswahili, kasi means speedy or speed. (http://en.wiktionary.org/wiki/kasi)


## Feature

- [ ] support HTTP, HTTPS service [12] [14]
- [x] virtualhost [1] [9] [16]
- [x] reshape the endpoint [7] [8]
- [x] user-defined middlewares using javascript [13]
- [ ] cache control [15] [17]
- [ ] timeout
- [ ] KeepAlive
- [ ] in-place response in config
- [ ] error back
- [ ] support `CORS` (Cross-Origin Resource Sharing) [2]
- [ ] statistic
- [ ] service discovery
- [x] configuration by YAML [3]


### In-Place Response

If your managed source has some trouble and you want to send redirect response to the client for safety, you can easily make custom response within configuration  e.g.

```yaml
response: |
  HTTP/1.1 302 Found
  Location: http://www.iana.org/domains/example/
```

### Error Back

When failed to get the expected response from source, the custom response will be sent to the client, that response can be made in configuration easily. As you guess, this error back also can be possible using middleware.

```yaml
when-error:
  status:
    - 400 - 499
    - 500
  response: |
    HTTP/1.1 302 Found
    Location: http://www.iana.org/
```


### User-Defined Middlewares

kasi has it's own middlewares and support the user-defined middlewares. You can write your own middlewares using `javascript`, so you easily manipulate the request and response.

> The current middleware behaviors are higly affected by the Django Web Frameworks.


```javascript
/*
`process_request` will be called before sending request to the target service.
*/
var process_request = function (endpointSetting, request) {
    // remove `If-Modified-Since`
    if (request.header["If-Modified-Since"]) {
        delete request.header["If-Modified-Since"];
    }
    
      // or, just return the response.
    var response = new Response();
    response.StatusCode = 304;
    response.Body = "";
    
    return response;
}

/*
`process_response` will be called after receiving the response from target service.
*/
var process_response = function (endpointSetting, reqeust, response) {
    response.Header["Expires"] = 10000;
    response.Body += "\n;";
    
    return response;
}
```


### CORS

If you make the web application with ajax and your target service does not provide the `CORS`, you can make kasi can do it.

With `cors` and it's children settings, you can easily support `CORS` by services, endpoints or globally.


### Service Discovery

kasi provides it's own API, so you can add or remove the services without any break or reload of kasi.


### Statistic

kasi provides the special dashboard, it has the special statistic page. It will show you the current statistics for services.

- number of requests by services and it's endpoints
- statistics of status code by services and it's endpoints
- etc.

> You also apply the "Measurement Protocole" of Google Analytics using the user defined middlewares by services. [6]


### Cache Control

Usually for performance reason, most of the web applications and servers support cache with their own ways, e.g. `ETAG`[4], or using headers like `If-Modified-Since` or `Expires`.

kasi respects their own cache mechanisms, and further more, explicitly supports cache.

With using cache of kasi, kasi will cache the data from the services and produce the response with it.


### Support Virtual Domain

Like `nginx`, kasi supports the virtualhost with SSL.


### Reshape Endpoints

If you connect to the API of github, you can rename the API endpoint like this,

`https://api.github.com/users/spikeekips`
to
`https://api.kasi.org/github/u/spikeekips/`

Naturally the original endpoint will be hided. With the custom middleware, which you can write, you can
manipulate the request and response.

The interesting feature in reshaping is the using the regular expression to reshape. e.g.

`https://api.github.com/users/(?P<username>)`
to
`https://api.kasi/{username}`


### timeout

In config, the `timeout` will set the tomeout globally, services or endpoints.


### YAML configuration Example

This is example configuration.

```yaml
%YAML 1.1

- default:
    cache:
        expire: 10m
        backend: [memory, memcache://127.0.0.1:11211, redis://127.0.0.1:6379]
    timeout: 5s

- service:
    bind: :8000
    hostname:
        - my0.github.com
        - my1.github.com
    ssl:
        cert: /secret/kasi-github.cert
        key: /secret/kasi-github.key
        pem: /secret/kasi-github.pem

    source: https://github.com/api/v1
    timeout: 10s
    endpoints:
        - endpoint:
            open: No
            to: /find/{username}
            from: /users/(?P<username>.*)
            source: [https://api0.github.com/v2, https://api1.github.com/v2]
        - endpoint:
            open: Yes
            source: https://api0.github.com/v2
            to: /github/\1
            from: (.*)
            timeout: 60s
            cors:
                allow-origin: http://localhost:9090
                # or allow-origin: [http://localhost:9090, http://localhost:9091]
                allow-methods: [GET, POST, PUT, OPTIONS, HEAD]
                allow-headers: [X-My-Header]
                allow-credentials: true
                max-age: 178000
...
```


## Todo

- [ ] more testing code
- [ ] refactoring the entire code
- [ ] clean up the monkey patch


[1]: https://gist.github.com/camoles/523dac8cc0fe40d52f66 "VirtualHost in Golang"
[2]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS "CORS"
[3]: http://www.yaml.org/start.html "YAML"
[4]: http://en.m.wikipedia.org/wiki/HTTP_ETag "ETAG"
[5]: http://en.m.wikipedia.org/wiki/List_of_HTTP_header_fields "Cache Control By Header"
[6]: https://developers.google.com/analytics/devguides/collection/protocol/v1/devguide "Measurement Analytics Protocol of Google Analytics"
[7]: https://github.com/StefanSchroeder/Golang-Regex-Tutorial "Golang-Regex-Tutorial"
[8]: https://regex-golang.appspot.com/assets/html/index.html "Regex Tester - Golang"
[9]: http://stackoverflow.com/questions/14170799/how-to-get-virtualhost-functionality-in-go "How to get “virtualhost” functionality in Go?"
[11]: https://github.com/pquerna/ffjson "ffjson"
[12]: https://github.com/epio/mantrid "Python based load-balancer"
[13]: https://github.com/robertkrimen/otto "A JavaScript interpreter in Go (golang) http://godoc.org/github.com/robertkrimen/otto"
[14]: http://fastah.blackbuck.mobi/blog/securing-https-in-go/ "The easy guide to securing HTTP + TLS with Go"
[15]: https://github.com/coocood/freecache "freecache"
[16]: http://www.reddit.com/r/golang/comments/34bem7/socket_master_a_zeroconfig_reverse_proxy/
[17]: https://github.com/coocood/freecache/blob/master/README.md "FreeCache"

