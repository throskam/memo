# Passwordless Authentication Architecture

## Why Passwordless

Passwordless is a pragmatic default when you want low operational overhead:

- no password storage lifecycle
- no password reset flows
- reduced credential management surface

## Flow

1. User submits email.
2. Server generates short-lived signed token (magic link).
3. Email transport sends verification link.
4. Verification endpoint validates token and resolves/creates user identity.
5. Session/auth token is established.

## Operational Notes

- Keep magic-link token lifetime short.
- Bind redirects defensively.
- Ensure email sender reliability and domain reputation.
- Audit login events and suspicious activity as needed.

For many apps, this is secure enough with strong token signing, short TTLs, and HTTPS.
