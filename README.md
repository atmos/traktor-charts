# traktor-charts

Try to copy [Traktor DJ Charts Export](http://tomashg.com/?p=1132) in go.

[Example Output](https://gist.github.com/atmos/0ae724237f1ef859f25a).

# Development

    % brew install go hg
    % make boostrap
    % make

# Exit Status

| Exit Status | Description |
|---|-----------------------------------------------------------------|
| 0 | Successfully posted new archives to djcharts.io. |
| 1 | No unposted traktor archve files found. |
| 2 | Unable to pulse to djcharts.io. Credentials are probably wrong. |
| 3 | Unable to post new archives to djcharts.io. Credentials are probably wrong. |
