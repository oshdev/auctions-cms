# go hotwire

A demo of using [https://hotwire.dev](Hotwire) with a Go web server.

## Flamebait propaganda notes

### High level points

- Go developers value simplicity so why do so many of them enjoy Go and it's ease of use for the backend and then throw those values out when writing on the frontend by using something like React
    - This isn't to say React is "bad". It's just that it's an incredibly powerful tool with a very complex ecosystem behind it that most websites don't need.
- [You shouldn't need a frontend framework for most websites](https://dev.to/gypsydave5/why-you-shouldnt-use-a-web-framework-3g24). Instead, a lot of websites can be built on top of basic, battle-tested web technologies and take an approach of [graceful degredation](https://developer.mozilla.org/en-US/docs/Glossary/Graceful_degradation) to improve the experience for users of modern web browsers.
- Hotwire lets you keep your interesting domain logic on the server, written in your favourite programming language without having to work with `X-language to JavaScript` transpilers which add their own level of complexity.