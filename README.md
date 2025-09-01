# BBAUM Steak Parser

A Golang application that monitors the BBAUM steaks catalog and sends Telegram notifications when specified steaks are available.

## Features

- Hourly monitoring of https://www.bbaum.ru/catalog/steak/
- Configurable steak tracking
- Telegram bot notifications
- Cron-based scheduling

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Configure the application by editing `config.yaml`:
```yaml
telegram:
  token: "YOUR_BOT_TOKEN"     # Get from @BotFather
  chat_id: "YOUR_CHAT_ID"     # Your Telegram chat ID

tracking:
  url: "https://www.bbaum.ru/catalog/steak/"
  interval: "0 * * * *"       # Cron format (every hour)
  steaks_to_track:
    - "Рибай"
    - "Стриплойн"
    # Add more steaks to track
```

3. Get your Telegram bot token:
   - Message @BotFather on Telegram
   - Create a new bot with `/newbot`
   - Copy the token to config.yaml

4. Get your chat ID:
   - Message your bot
   - Visit: `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`
   - Find your chat ID in the response

5. Run the application:
```bash
go run .
```

## Configuration

### Tracking Interval
The `interval` field uses cron format:
- `"0 * * * *"` - Every hour
- `"0 */2 * * *"` - Every 2 hours
- `"0 9-17 * * *"` - Every hour from 9 AM to 5 PM

### Steaks to Track
Add steak names to the `steaks_to_track` list. The parser will match steaks containing these names (case-insensitive).

## How it Works

1. The application starts a cron scheduler
2. Every hour (or as configured), it fetches the BBAUM steaks page
3. Parses HTML to extract steak information
4. Filters results based on configured steak names
5. Sends Telegram notifications for matching steaks

## Dependencies

- `github.com/PuerkitoBio/goquery` - HTML parsing
- `github.com/robfig/cron/v3` - Cron scheduling
- `gopkg.in/yaml.v2` - YAML configuration

## Docker Deployment

### Build and run with Docker:

1. **Build the Docker image**:
   ```bash
   docker build -t bbaum-parser .
   ```

2. **Run the container**:
   ```bash
   docker run -d --name bbaum-parser bbaum-parser
   ```

3. **View logs**:
   ```bash
   docker logs -f bbaum-parser
   ```

4. **Stop the container**:
   ```bash
   docker stop bbaum-parser
   ```

### Docker files:
- `Dockerfile` - Container configuration
- `.dockerignore` - Files to exclude from Docker build

### Notes:
- Container will run the parser continuously with cron scheduling
- All logs are output to stdout for easy monitoring
- Container will restart automatically on crashes if run with `--restart=always`