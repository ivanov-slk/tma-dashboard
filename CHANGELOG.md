# Changelog

## [1.1.0](https://github.com/ivanov-slk/tma-dashboard/compare/v1.0.0...v1.1.0) (2023-10-20)


### Features

* The HTTP handler now returns the first message on the NATS data generator subject. ([3050932](https://github.com/ivanov-slk/tma-dashboard/commit/30509328fc109d9796a4da19db574bf6befbf1d0))


### Refactoring/Restructuring

* Create a dedicated HTTP handler function to use in main.go. ([a0023e7](https://github.com/ivanov-slk/tma-dashboard/commit/a0023e7b8eb5856f9f3b8082e51179b85cf2df6d))
* Use the dedicated handler function in the main program. ([c8334ab](https://github.com/ivanov-slk/tma-dashboard/commit/c8334ab25e0f750e21e3fc75fb1b565f18a5a550))


### Testing

* Add NATS server for testcontainers. ([4b4bbb6](https://github.com/ivanov-slk/tma-dashboard/commit/4b4bbb61312958bd21ff0ba14f6001cb28efe669))

## [1.0.0](https://github.com/ivanov-slk/tma-dashboard/compare/v0.0.1...v1.0.0) (2023-09-29)


### Features

* Created a 'hello' service, returning just 'hello message'. ([334e878](https://github.com/ivanov-slk/tma-dashboard/commit/334e878512f3d13cd8b632900d8df870379f258b))

## 0.0.1 (2023-09-24)


### Continuous Integration & Continuous Delivery

* Add github actions. ([1f2b908](https://github.com/ivanov-slk/tma-dashboard/commit/1f2b908557579279f73392b8fc4a40518ad137d3))
