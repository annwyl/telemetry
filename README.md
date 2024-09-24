# Telemetry Package

Package for telemetry/logging

## Features

- Multiple logging levels (Debug, Info, Warning, Error)
- Configurable logging backends (drivers)
- Transaction support for tracking related log entries
- JSON configuration
- Thread safe

## Basic Usage

There is an example in cmd/main.go

## Configuration

Configuration is done through the `config.json` like:

```json
{
  "driver": "console",
  "driver_config": "logs.txt",
  "log_level": 1,
  "default_tags": {
    "environment": "development",
    "go_version": "1.22",
    "app_version": "1.0.0"
  }
}
```

## Extending the Package

You can write your own driver by putting it into the drivers folder, and specifing it in the `config.json`. There are multiple drivers already, which can be used as an example or starting point.

### License

Use it as you wish