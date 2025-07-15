# Frens

A friendship management & journaling application for introverts and not only.
Build relationships with people that last.

`frens` is a command-line application that helps you to keep track of your relationships 
with other people you care about.

`frens` gives you:

- An organized and systematic way around staying in touch with friends and family.
- A low-effort way to record and remember big moments in your life.
- A way to track how your relationships develop over time.

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

One of the major `fren`'s features is the ability to input all data as a free-form text 
using a simple and straightforward syntax.

### Tags

Tags are one of the common parts of other pieces of information like `Activities`, `Notes`, `Friends` or `Locations`.

Tags can be specified via the `#tag` syntax like this:

```text
#scool #family #university #school:math #school:physics #family-extended
```

### Locations

`Locations` can be added like this:

```text
Scranton, USA (aka "The Electric City") :: a great place to live and work #office @Scranton $id:scranton
```

Then, you can set location for your `Friends`, `Activities` and `Notes` via the `@location` syntax like this:

```text
@NewYork @LosAngeles @Scranton
```

### Friends

Similarly to `Locations`, the basic information about `Friend` can be added using a similar syntax:

```text
Michael Harry Scott (aka "The World's Best Boss"), my Dunder Mifflin boss, is a great friend of mine #office @Scranton
```

You can also specify ID of the friend via `$id:friendid` syntax like this:

```text
Michael Harry Scott (aka "The World's Best Boss") $id:mscott, ...
```

### Activities & Notes

`Activities` and `Notes` are events that have the same syntax:

```text
yesterday :: Jim put my stuff in jello #office @Scranton
```

You can also completely omit the date, so it will be set to the current date:

```text
Dwight bought a new beet farm #office @Scranton
```

## Credits

Inspired by awesome [JacobEvelyn/friends](https://github.com/JacobEvelyn/friends).

Made with ❤️ by [Roman Glushko](https://github.com/roma-glushko).