# Multiple languages support - i18n

The library used for localizing the messages is [go-i18n](https://github.com/nicksnyder/go-i18n).

## Setting the language

You can change the language in which the bot replies. The language is set per chat/user basis and can be changed by
issuing the `/lang` or `/language` command. The command accepts one parameter - the language tag, which must suffice the
**BCP 47** standard.

Example command:

```text 
/lang sv   #switches to swedish 
/lang en   #switches to english 
```

If the message set does not have the desired language translations in the `i18n/translations` folder, the message will
be in English. When first registering a user, the default language is English.

## Contributing to translations

Every translation contribution is welcome! Issue a PR with the translation(s). You only need to add a translation yaml
file - everything else is handled.
