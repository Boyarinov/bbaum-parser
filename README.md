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

## Railway Deployment

### Setup Steps

1. **Create Railway account**: Go to [railway.app](https://railway.app) and sign up

2. **Push code to GitHub**: Make sure your code is in a GitHub repository

3. **Deploy to Railway**:
   - Login to Railway dashboard
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository
   - Railway will auto-detect the Dockerfile and build

4. **Set environment variables** in Railway dashboard:
   ```
   TELEGRAM_TOKEN=your_bot_token
   TELEGRAM_CHAT_ID=your_chat_id
   ```

5. **Update config.yaml** to use environment variables (optional):
   ```yaml
   telegram:
     token: "${TELEGRAM_TOKEN}"
     chat_id: "${TELEGRAM_CHAT_ID}"
   
   tracking:
     url: "https://www.bbaum.ru/catalog/steak/"
     interval: "0 * * * *"
     steaks_to_track:
       - "Рибай"
       - "Стриплойн"
   ```

### Files for Railway deployment:
- `Dockerfile` - Container configuration
- `railway.json` - Railway-specific settings
- `.dockerignore` - Files to exclude from Docker build

### Notes:
- Railway provides free tier with 500 hours/month
- App will automatically restart on crashes
- Logs are available in Railway dashboard
- No port configuration needed - Railway handles networking