# Frens

A friendship management & journaling application for introverts and not only.
Build relationships with people that lasts.

## Get Started

TBU

## Main Concepts

![Diagram](./docs/friens-data-model.png?raw=true)

## Language

One of the major Fren's features is the ability to input all data as a free-form text using a simple and straightforward syntax.

### Tags

Tags are one of the common parts of other pieces of information like `Activities`, `Notes`, `Friends` or `Locations`.

Tags can be specified via the `#tag` syntax like this:

```text
#scool #family #university #school:math #school:physics #family-extended
```

### Locations

```text
Scranton, USA (aka "The Electric City") :: a great place to live and work #office @Scranton $id:scranton
```

Then, you can set location for your Friends via the `@location` syntax like this:

```text
@NewYork @LosAngeles @Scranton
```

### Friends

The basic Friend information can be inputted like this:

```text
Michael Harry Scott (aka "The World's Best Boss"), my Dunder Mifflin boss, is a great friend of mine #office @Scranton
```

You can also specify ID of the friend via `$id:nickname` syntax like this:

```text
Michael Harry Scott (aka "The World's Best Boss") $id:mscott, ...
```

#### Contacts

```text

```

#### Wishlist

You can add a wishlist to your Friend like this:

```text
```

### Activities & Notes

```text
yesterday :: Jim put my stuff in jello #office @Scranton
```

## Credits

Inspired by awesome [JacobEvelyn/friends](https://github.com/JacobEvelyn/friends).

Made with ❤️ by [Roman Glushko]().