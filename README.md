# Spam-fighter-bot

## Overview

This bot is used to verify that the user who entered the chat is not a bot
and knows Russian language. The reason for creating the bot was that many 
spammers solve English captcha manually and then write ads. At the same 
time, other bots contained too much unnecessary information and filled chat
with it. So I decided to create this bot

**Link to the bot**: https://t.me/spam_fighter_bot

## Example

![example](/.github/example.png)

### Project Navigation

- [CMD start file](./cmd/spam_fighter_bot/spam_fighter_bot)
- [Welcome message logic](./internal/handler/userJoined.go)
- [Checking correct answer](./internal/handler/text.go)
