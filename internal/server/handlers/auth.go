package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/throskam/ki"
	"github.com/throskam/kix/auth"
	"github.com/throskam/kix/htmx"
	"github.com/throskam/kix/i18n"
	"github.com/throskam/kix/sess"
	"github.com/throskam/memo/internal/lib"
	"github.com/throskam/memo/internal/orm"
	"github.com/throskam/memo/internal/views/pages"
)

type AuthController struct {
	jwks    auth.JWKS
	us      *lib.UserService
	smtp    string
	baseURL *url.URL
	from    string
	apiKey  string
}

func NewAuthController(us *lib.UserService, jwks auth.JWKS, smtp string, baseURL *url.URL, form string, apiKey string) *AuthController {
	c := &AuthController{
		jwks:    jwks,
		us:      us,
		smtp:    smtp,
		baseURL: baseURL,
		from:    form,
		apiKey:  apiKey,
	}

	return c
}

func (c *AuthController) PageGet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		RenderProblem(w, r, NewProblem(fmt.Errorf("failed to parse form: %w", err)))
		return
	}

	Render(w, r, pages.AuthPage(pages.AuthPageProps{
		State: r.Form,
	}))
}

func (c *AuthController) PasswordlessSend(w http.ResponseWriter, r *http.Request) {
	form := htmx.NewFormFromRequest(r, &pages.AuthPasswordlessForm{})

	if !form.OK() {
		Render(w, r, pages.AuthPasswordless(pages.AuthPasswordlessProps{
			State:   r.URL.Query(),
			Form:    form,
			Success: false,
		}))
		return
	}

	jwt, err2 := auth.GenerateJWT(c.jwks, auth.Claims{
		"email":        form.Data.Email,
		"redirect_url": r.FormValue("redirect_url"),
	}, time.Minute*15)
	if err2 != nil {
		RenderProblem(w, r, NewProblem(fmt.Errorf("failed to generate JWT: %w", err2)))
		return
	}

	link := c.baseURL.ResolveReference(ki.GetLocation(r.Context(), "auth:passwordless:verify").WithQueryParam("token", jwt).URL())
	from := c.from
	to := []string{form.Data.Email}
	subject := i18n.T(r.Context(), "Connect")
	body := i18n.T(r.Context(), "Follow the link to connect : %s", link.String())

	err3 := lib.SendEmail(from, to, subject, body, lib.WithAPIKey(c.apiKey), lib.WithSMTP(c.smtp))
	if err3 != nil {
		RenderProblem(w, r, NewProblem(fmt.Errorf("failed to send mail: %w", err3)))
		return
	}

	Render(w, r, pages.AuthPasswordless(pages.AuthPasswordlessProps{
		State:   r.URL.Query(),
		Form:    htmx.NewForm(&pages.AuthPasswordlessForm{}),
		Success: true,
	}))
}

func (c *AuthController) PasswordlessVerify(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ParseJWT(c.jwks, r.FormValue("token"), time.Minute*5)
	if err != nil {
		RenderProblem(w, r, NewProblem(fmt.Errorf("failed to parse JWT: %w", err)))
		return
	}

	email := token["email"].(string)

	authenticationMethod := lib.NewAuthenticationMethod(orm.AuthenticationProviderPasswordless, email)

	user, err := c.us.GetByAuthenticationMethodOrCreate(r.Context(), authenticationMethod)
	if err != nil {
		RenderProblem(w, r, NewProblem(fmt.Errorf("failed to get or create user: %w", err)))
		return
	}

	claims := map[string]any{}

	claims["sub"] = user.ID

	bearer, err := auth.GenerateJWT(c.jwks, claims, time.Hour*30*24)
	if err != nil {
		RenderProblem(w, r, NewProblem(fmt.Errorf("failed to generate jwt: %w", err)))
		return
	}

	session := sess.MustGetSession(r.Context())
	session.Add("bearer", bearer)

	redirectURL, err := url.QueryUnescape(token["redirect_url"].(string))
	if err != nil || redirectURL == "" {
		redirectURL = "/"
	}

	htmx.Redirect(w, r, redirectURL)
}