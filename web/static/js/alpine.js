document.addEventListener("alpine:init", () => {
  // Keep focus between htmx swaps if the active element has an id.
  let lastActiveElementId = null;

  document.body.addEventListener("htmx:beforeSwap", () => {
    lastActiveElementId = document.activeElement?.id || null;
  });

  document.body.addEventListener("htmx:afterSwap", () => {
    document.getElementById(lastActiveElementId)?.focus?.();
    lastActiveElementId = null;
  });

  // Alpine components are not properly destroyed on htmx boosted requests.
  // This is a workaround to force components destruction.
  // @see https://github.com/alpinejs/alpine/discussions/4485
  document.body.addEventListener("htmx:beforeSwap", (evt) => {
    if (!evt.detail?.requestConfig?.boosted) {
      return;
    }

    Alpine.destroyTree(evt.detail.target);
  });

  // Format time element based on locale.
  Alpine.data("time", () => {
    const options = {
      year: "numeric",
      month: "long",
      day: "numeric",
    };

    return {
      init() {
        const formatted = new Intl.DateTimeFormat(
          document.documentElement.lang,
          options,
        ).format(new Date(this.$el.getAttribute("datetime")));

        this.$el.textContent = formatted;
      },
    };
  });

  // Define smart focus magic.
  Alpine.magic("smartFocus", (el) => {
    const focusMagic = Alpine.$data(el).$focus;

    return {
      ...focusMagic,
      focusInput() {
        const length = el.value?.length ?? 0;

        el.focus();

        if (length > 0) {
          el.setSelectionRange(length, length);
        }
      },
      focusedItem(selector) {
        const focus = focusMagic.within([...el.querySelectorAll(selector)]);

        focus.wrap();

        if (!focus.focusables().some((el) => document.activeElement === el)) {
          focus.first();
          return focus;
        }

        if (!el.contains(focusMagic.focused())) {
          return null;
        }

        return focus;
      },
      focusedItemNode(selector) {
        this.focusedItem(selector);

        return this.focused();
      },
    };
  });

  Alpine.directive("autofocus", (el) => {
    const length = el.value.length;

    el.focus();
    el.setSelectionRange(length, length);
  });

  const TYPABLE_INPUT_TYPES = new Set([
    "text",
    "search",
    "email",
    "url",
    "tel",
    "password",
    "number",
    "date",
    "datetime-local",
    "month",
    "time",
    "week",
  ]);

  const isWriting = () => {
    const active = document.activeElement;
    if (!(active instanceof HTMLElement)) {
      return false;
    }

    if (active.isContentEditable || active.getAttribute("role") === "textbox") {
      return true;
    }

    if (active.nodeName.toLowerCase() === "textarea") {
      return true;
    }

    if (active.nodeName.toLowerCase() !== "input") {
      return false;
    }

    const type = (active.getAttribute("type") || "text").toLowerCase();
    return TYPABLE_INPUT_TYPES.has(type);
  };

  window.tinykeys.tinykeys(window, {
    Escape: (event) => {
      if (isWriting()) {
        event.preventDefault();

        document.activeElement.blur();
      }
    },
  });

  Alpine.directive(
    "shortcut",
    (el, { modifiers, expression }, { evaluateLater, cleanup }) => {
      const binding = modifiers.reduce((acc, modifier) => {
        switch (modifier) {
          case "shift":
            return `${acc}Shift`;
          case "and":
            return `${acc}+`;
          case "then":
            return `${acc} `;
          case "space":
            return `${acc}Space`;
          default:
            return `${acc}${modifier}`;
        }
      }, "");

      const evaluate = evaluateLater(expression);

      const unsubscribe = window.tinykeys.tinykeys(window, {
        [binding]: (e) => {
          if (isWriting()) {
            return;
          }

          e.preventDefault();

          evaluate();
        },
      });

      cleanup(() => {
        unsubscribe();
      });
    },
  );
});
