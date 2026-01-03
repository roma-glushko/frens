# Frens

A friendship management & journaling application for introverts and not only.
Build relationships with people that last.

`frens` is a command-line application that helps you to keep track of your relationships 
with other people you care about.

`frens` gives you:

- An organized and systematic way around staying in touch with friends and family.
- A low-effort way to record and remember big moments in your life.
- A way to track how your relationships develop over time.

> [!NOTE]
>
> 🚧 **This project is under active development, so some CLI or behaviors may change.** 🚧

## Features

- Record your relationships with friends, family, colleagues, and acquaintances using a simple set of concepts like `Friends`, `Locations`, `Activities`, and `Notes`.
- A simple journaling language (called `frentxt`) to simplify the process of recording your thoughts and activities.

## Philosophy

- **Simplicity**: Should be quick and easy to use. As few concepts as possible to keep in mind. Little to no manual record maintenance.
- **Journaling First**: Should focus you on journaling and jotting down your thoughts.
- **Intelligence**: Guessing your friends' names, understanding relative dates (e.g. "yesterday", "tomorrow").
- **Privacy & Transparency**: All data is stored locally on your machine in TOML file format. You can optionally share it across your laptops via Git.
- **Hackable**: Should be possible to use the collected data in automations and scripts.

## Installation

### MacOS

```bash
brew tap roma-glushko/frens https://github.com/roma-glushko/frens
brew install frens
```

### Download Binaries

For other platforms and architectures, you can download `frens`' binaries right from [Github Releases](https://github.com/roma-glushko/frens/releases).

## Main Concepts

![Diagram](./docs/friens-data-model.png?raw=true)

- **Friends**: People you know and care about. Can be family, colleagues, or acquaintances.
- **Locations**: Places where you and your friends live, work, or spend time together.
- **Activities**: Things you do with your friends, like going to the movies, having dinner, or attending events.
- **Notes**: Insights, preferences, deep meaning information with long-term value about your friends, activities, or locations.

## Language

One of the major `fren`' features is the ability to input all data as a free-form text
using a simple and straightforward syntax.

### Tags

Tags are a common element across all entities like `Friends`, `Locations`, `Activities`, `Notes`, `Contacts`, and `Wishlist` items.

Tags can be specified via the `#tag` syntax:

```text
#family #university #work:engineering #hobby-photography
```

Tags support namespaces with colons and hyphens for multi-word tags.

### Properties

Properties allow you to set specific attributes on entities using the `$property:value` syntax:

```text
$id:custom-id $price:$50 $cal:hebrew
```

### Locations

`Locations` represent places where you and your friends live, work, or spend time together:

```text
Scranton, USA (aka "The Electric City") :: a great place to live and work #office $id:scranton
```

```text
Berlin, Germany :: vibrant tech hub with amazing coffee culture #europe #tech
```

```text
Tokyo (aka "東京") :: visited during cherry blossom season #travel #asia
```

Reference locations in other entities via the `@location` syntax:

```text
@NewYork @Berlin @Tokyo
```

### Friends

Basic information about a `Friend` can be added using a similar syntax:

```text
Michael Harry Scott (aka "The World's Best Boss") :: my Dunder Mifflin boss #office @Scranton $id:mscott
```

```text
Sarah Chen (aka "Saz", "SC") :: college roommate, now works at Google #college #tech @SanFrancisco
```

```text
Hans Mueller :: met at a conference in Berlin, shares my love for hiking #conference #outdoor @Berlin
```

### Contacts

`Contacts` store contact information for your friends with support for various platforms:

```text
email:sarah@example.com phone:+1234567890 #work
```

```text
tg:@telegram_user gh:roma-glushko ig:@photography_account #personal
```

Supported contact types and their aliases:
- `email` / `mail` - email addresses
- `phone` / `tel` - phone numbers
- `telegram` / `tg` - Telegram handles
- `whatsapp` / `wa` - WhatsApp contacts
- `twitter` / `x` - Twitter/X handles
- `linkedin` / `li` - LinkedIn profiles
- `github` / `gh` - GitHub profiles
- `instagram` / `ig` - Instagram handles
- `facebook` / `fb` - Facebook profiles
- `discord`, `slack`, `signal` - messaging platforms

Auto-detection works for common formats:

```text
john@company.com +48123456789 #work
```

### Dates

`Dates` track important dates for your friends like birthdays and anniversaries:

```text
January 15 :: Birthday #birthday
```

```text
March 10, 2020 :: Wedding anniversary #anniversary #celebration
```

```text
15 Nisan :: Passover celebration #holiday $cal:hebrew
```

Supported calendars: `gregorian` (default), `hebrew`.

### Wishlist

`Wishlist` items help you track gift ideas for your friends:

```text
Kindle Paperwhite https://amazon.com/kindle $price:$140 #reading #gift-ideas
```

```text
Specialty coffee subscription #coffee #monthly
```

```text
Concert tickets for Coldplay @Berlin $price:$120 #music #experience
```

### Activities & Notes

`Activities` and `Notes` are events that share the same syntax:

```text
yesterday :: Jim put my stuff in jello #office @Scranton
```

```text
2024-12-25 :: Christmas dinner with the whole family #family #holiday @Home
```

```text
last week :: Helped Sarah move to her new apartment #friends @SanFrancisco
```

You can omit the date, and it will default to the current date:

```text
Grabbed coffee and caught up on life #casual @Berlin
```

Relative dates like "yesterday", "last week", "2 days ago" are supported.

## Credits

Inspired by awesome [JacobEvelyn/friends](https://github.com/JacobEvelyn/friends).

Made with ❤️ by [Roman Glushko](https://github.com/roma-glushko).