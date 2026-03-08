# UI and `Props` Conventions

## Rule

Every reusable component receives a `Props` struct.

## Why

- Explicit API between caller and component.
- Easier extension without breaking callsites.
- Better readability and testability.

## Conventions

- Use `ComponentNameProps` naming.
- Include only inputs the component needs.
- Prefer plain types in props; avoid hidden context coupling.
- Keep accessibility fields explicit (`ID`, labels, validation message bindings).
- Pass extra HTML attributes via an `Attrs` field when appropriate.

## Example Pattern

```go
type InputProps struct {
    ID                 string
    Name               string
    Kind               string
    Label              string
    Placeholder        string
    Value              string
    ValidationMessages []string
    Attrs              templ.Attributes
}
```

Treat `views/ui` as a small design system for server-driven pages.
