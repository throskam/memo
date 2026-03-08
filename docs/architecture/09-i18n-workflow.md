# Internationalization Workflow

## Translation Strategy

Use Go text extraction to keep translation keys near templates and handler messages.

## File Roles

- `internal/translations/locales/<locale>/out.gotext.json`: generated extraction output.
- `internal/translations/locales/<locale>/messages.gotext.json`: manually edited translation source.
- `internal/translations/catalog.go`: generated runtime catalog.

## Workflow

1. Add or update `i18n.T(...)` keys in code.
2. Run `make translations`.
3. Sync locale messages from output (example for `fr-FR`):

```bash
rm internal/translations/locales/fr-FR/messages.gotext.json
cp internal/translations/locales/fr-FR/out.gotext.json internal/translations/locales/fr-FR/messages.gotext.json
```

4. Fill missing `translation` values in `messages.gotext.json`.
5. Run `make translations` again.

Example command flow in this repository:

```bash
make translations
rm internal/translations/locales/fr-FR/messages.gotext.json
cp internal/translations/locales/fr-FR/out.gotext.json internal/translations/locales/fr-FR/messages.gotext.json
# edit internal/translations/locales/fr-FR/messages.gotext.json
make translations
```

## Notes

- Warnings like `fr-FR: Missing entry ...` after the first `make translations` are expected until translations are added in `messages.gotext.json`.
- Repeat the sync step whenever new keys are extracted.
- Treat locale files and generated catalog as versioned artifacts reviewed like code.
