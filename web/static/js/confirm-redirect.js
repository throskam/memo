document.addEventListener("DOMContentLoaded", () => {
  const handler = function (e) {
    const elems = document.querySelectorAll("[data-confirm-redirect]");

    for (const elem of elems) {
      const msg = elem.getAttribute("data-confirm-redirect");
      const ok = confirm(msg);

      if (!ok) {
        e.preventDefault();
        return;
      }
    }
  };

  window.addEventListener("beforeunload", handler);

  // HTMX boosted link.
  document.body.addEventListener("htmx:beforeRequest", function (e) {
    if (e.detail.boosted) {
      handler(e);
    }
  });
});
