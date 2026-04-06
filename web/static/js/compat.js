document.addEventListener('DOMContentLoaded', function() {
  if (/^((?!chrome|android).)*safari/i.test(navigator.userAgent)) {
    document.querySelector('link[href="/static/css/main.css"]').disabled = true;
    document.querySelector('link[href="/static/css/main.compat.css"]').disabled = false;
  }
});
