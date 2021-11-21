# Spam-fighter-bot

![pcherednichenko](https://circleci.com/gh/pcherednichenko/spam_fighter_bot.svg?style=svg)

## Overview

This bot is used to verify that the user who entered the chat is not a bot
and knows Russian language. The reason for creating the bot was that many 
spammers solve English captcha manually and then write ads. At the same 
time, other bots contained too much unnecessary information and filled chat
with it. So I decided to create this bot

**Link to the bot**: https://t.me/spam_fighter_bot

## Example

![example](/.github/example.png)

![example_en](/.github/example-en.jpg)

### Project Navigation

- [CMD start file](./cmd/spam_fighter_bot/spam_fighter_bot.go)
- [Welcome message logic](./internal/app/handler/userJoined.go)
- [Checking correct answer](./internal/app/handler/text.go)

### Article about this bot

https://blogpavel.com/2021/04/23/spam-fighter-bot/

### Did you like it?

If you like this bot and you use it put a star :star: , I will be pleased.

The bot is free, no ads, open source available for everyone :heart: