# Changelog

## [1.3.0](https://github.com/ivanov-slk/tma-dashboard/compare/v1.2.6...v1.3.0) (2024-01-18)

### Features

- Messages that produce error are logged to `stdout` instead of being returned to the user. Instead, the last valid message is returned by the handler. ([d56c62d](https://github.com/ivanov-slk/tma-dashboard/commit/d56c62d89c0786d9ca28828060d4a51528cfadd9))

## [1.2.6](https://github.com/ivanov-slk/tma-dashboard/compare/v1.2.5...v1.2.6) (2024-01-13)

### Maintenance

- Add more logging to the server. ([2466e1a](https://github.com/ivanov-slk/tma-dashboard/commit/2466e1a5c8d9d8a1ce2716716eee72eae05c0186))

## [1.2.5](https://github.com/ivanov-slk/tma-dashboard/compare/v1.2.4...v1.2.5) (2024-01-13)

### Maintenance

- In case invalid message has been fetched, log the error and the message as string, instead of failing silently. ([6067695](https://github.com/ivanov-slk/tma-dashboard/commit/60676955af028fbcea8ab76f0996fa1c30fd9469))

### Bug Fixes

- Explicitly use jetstream.MemoryStorage during stream creation. This is sufficient for the current purposes of the service. ([14e6d2a](https://github.com/ivanov-slk/tma-dashboard/commit/14e6d2ad66542eecc65afef4d84c867565aa0dd0))

## [1.2.4](https://github.com/ivanov-slk/tma-dashboard/compare/v1.2.3...v1.2.4) (2024-01-13)

### Maintenance

- Service creates a stream if it has not been previously created. ([d58b3cf](https://github.com/ivanov-slk/tma-dashboard/commit/d58b3cfbce3daf6bcbe277b4ec768ac8ba3e16fc))
- Service now exits with error immediately in case of errors during connecting to the NATS broker. ([134c4ef](https://github.com/ivanov-slk/tma-dashboard/commit/134c4efce09448772355fbdde0fdd83f2d8ad066))

## [1.2.3](https://github.com/ivanov-slk/tma-dashboard/compare/v1.2.2...v1.2.3) (2024-01-13)

### Maintenance

- Add more logs during connecting to NATS. ([0bdf4fc](https://github.com/ivanov-slk/tma-dashboard/commit/0bdf4fc01ea45feae7bbe30b9295abeb0379d7f7))

## [1.2.2](https://github.com/ivanov-slk/tma-dashboard/compare/v1.2.1...v1.2.2) (2024-01-13)

### Maintenance

- Add more logs during connecting to NATS. ([608de51](https://github.com/ivanov-slk/tma-dashboard/commit/608de5161be8e0cf51f1d6d25545e22278572b89))

## [1.2.1](https://github.com/ivanov-slk/tma-dashboard/compare/v1.2.0...v1.2.1) (2024-01-03)

### Bug Fixes

- Expect that the NATS stream is already created. The dashboard should not be responsible for creating the stream. ([cc231fc](https://github.com/ivanov-slk/tma-dashboard/commit/cc231fc0dd1cdbc301fe3f6dd776df26763a6f95))

## [1.2.0](https://github.com/ivanov-slk/tma-dashboard/compare/v1.1.1...v1.2.0) (2024-01-03)

### Documentation

- Add docstrings for the important objects and packages. ([5935804](https://github.com/ivanov-slk/tma-dashboard/commit/59358045ac677065c193c047ca9db5e0ef2b07e2))
- Document important behavioral details of the HTTP server. ([64ac0ef](https://github.com/ivanov-slk/tma-dashboard/commit/64ac0ef7786660dd9d0c87511bc493d4a3786bea))

### Features

- Make the NATS URI convigurable via the environment variable . ([f6f0530](https://github.com/ivanov-slk/tma-dashboard/commit/f6f05301f0030ed6874b21736d29a6bd99bfc841))
- The service now parses an input message in the proper format and outputs the supplied temperature. ([c6c57e4](https://github.com/ivanov-slk/tma-dashboard/commit/c6c57e4f3273fcf4974fab91d0c025c4dd27c646))

### Refactoring/Restructuring

- Convert the internal channel to array of bytes to avoid unnecessary back-and-forth string conversion. ([f8116ce](https://github.com/ivanov-slk/tma-dashboard/commit/f8116ce727a09c28d8224e765cabff284413d60a))
- Decouple the HTTP server from the input data source, use channels for internal data transfer. ([f9ccd65](https://github.com/ivanov-slk/tma-dashboard/commit/f9ccd654ed0aee7acb946af7fb942f1ed85e1f1b))
- Remove unneeded stubs used in the HTTP server tests. ([d863d44](https://github.com/ivanov-slk/tma-dashboard/commit/d863d44980fde0687c30406a142ff49cf3b8513f))

## [1.1.1](https://github.com/ivanov-slk/tma-dashboard/compare/v1.1.0...v1.1.1) (2023-11-30)

### Refactoring/Restructuring

- Convert the HTTP handler and NATS client to structs, allowing for better modularity and testability. ([dece9ee](https://github.com/ivanov-slk/tma-dashboard/commit/dece9ee93fa5ea15df7a1b6ad3c3b25fdfd76d2d))
- Move all NATS connection code to a dedicated function, and the connection setup in a dedicated type. ([514382c](https://github.com/ivanov-slk/tma-dashboard/commit/514382c225bec9571a3df8ecce84fd29d19df8e5))
- Move the NATS part into a dedicated package. ([a55626a](https://github.com/ivanov-slk/tma-dashboard/commit/a55626afff6ae23fb814293d0645e2c34f436b81))

## [1.1.0](https://github.com/ivanov-slk/tma-dashboard/compare/v1.0.0...v1.1.0) (2023-10-20)

### Features

- The HTTP handler now returns the first message on the NATS data generator subject. ([3050932](https://github.com/ivanov-slk/tma-dashboard/commit/30509328fc109d9796a4da19db574bf6befbf1d0))

### Refactoring/Restructuring

- Create a dedicated HTTP handler function to use in main.go. ([a0023e7](https://github.com/ivanov-slk/tma-dashboard/commit/a0023e7b8eb5856f9f3b8082e51179b85cf2df6d))
- Use the dedicated handler function in the main program. ([c8334ab](https://github.com/ivanov-slk/tma-dashboard/commit/c8334ab25e0f750e21e3fc75fb1b565f18a5a550))

### Testing

- Add NATS server for testcontainers. ([4b4bbb6](https://github.com/ivanov-slk/tma-dashboard/commit/4b4bbb61312958bd21ff0ba14f6001cb28efe669))

## [1.0.0](https://github.com/ivanov-slk/tma-dashboard/compare/v0.0.1...v1.0.0) (2023-09-29)

### Features

- Created a 'hello' service, returning just 'hello message'. ([334e878](https://github.com/ivanov-slk/tma-dashboard/commit/334e878512f3d13cd8b632900d8df870379f258b))

## 0.0.1 (2023-09-24)

### Continuous Integration & Continuous Delivery

- Add github actions. ([1f2b908](https://github.com/ivanov-slk/tma-dashboard/commit/1f2b908557579279f73392b8fc4a40518ad137d3))
