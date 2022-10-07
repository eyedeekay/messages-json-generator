# messages-json-generator

If you knew what you were doing when you started translating your WebExtension then you
probably don't need this. If you didn't, this might help you find and fix your mistakes
more quickly.

This is a very messy but effective program for generating translatable `messages.json`
files from the HTML files used to build pages and user-interfaces for WebExtensions in
Firefox and Chrome. It scans a directory for HTML files, then scans the files for elements
which have an `id`. After performing a few checks, It gathers the default `textData` from the
HTML tags and places it into a properly-formatted `messages.ex.json` file. Then it generates
a `messages.js` file which automatically enables the translations on pages where you include
it with:

```HTML
<script src="messages.js"></script>
```

If you do not have any existing translations, `messages.ex.json` simply be copied to
`_locales/$lang/messages.json`. If you do have existing translations, there are additional steps
you need to take to migrate away incrementally.

1. Continue using your existing translation fields
2. **Combine** your existing `messages.json` file with the new `messages.ex.json` file

`messages-json-combine` is a supplementary program to help you combine 2 `messages.json` files
into one.

Once you're done, review your new translated text(things have probably been changed or overwritten
with default values determined by the HTML `textData`).

You now have translation coverage.

See also: https://github.com/eyedeekay/webext-translator.
