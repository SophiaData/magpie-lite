//go:build darwin && cgo && !nogui

package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// A closed NSWindow keeps its WKWebView's document visible unless magpie
// says otherwise. The page here records what it saw as it was hidden, and
// reports that record once it is shown again; without hideWebview's signal
// it remains a visible page with running animations.
//
// The page reads document.hidden as app.js does, from magpie's signal.
// WebKit does turn an ordered-out window's document.hidden on, but in its
// own time and without a visibilitychange of its own, so a page reading
// WebKit's value saw itself visible whenever magpie's signal arrived first,
// recorded nothing, and the window was shown again after a fixed pause with
// nothing to report (#4924). The window is now shown again only once the
// page has said it saw itself hidden.
//
// On mote (macOS 27.2) WebKit's own document.hidden stayed true from load
// to the end, the window shown or not, with no visibilitychange of its own.
// The shown signal's own visibilitychange then read hidden and unmarked,
// and was taken for the next hiding: its record and its /hidden reached
// the host before the panel's Escape was told, and the panel's second
// round failed with marked=false play="running" (#5322). Only a
// visibilitychange outside the shown signal counts as a hiding now.
func TestClosingWindowMarksPageOutOfSight(t *testing.T) {
	out, err := runAppKit(t, 30*time.Second, []string{"MAGPIE_TEST_WINDOW_VISIBILITY=1"})
	if err != nil || !strings.Contains(string(out), "window visibility: ok") {
		t.Fatalf("window visibility: %v\n%s", err, out)
	}
}

func init() {
	if os.Getenv("MAGPIE_TEST_WINDOW_VISIBILITY") != "1" {
		return
	}
	ready := make(chan struct{}, 1)
	reports := make(chan url.Values, 1)
	app := application.New(application.Options{
		Name: "magpie window visibility regression",
		Mac:  application.MacOptions{ActivationPolicy: application.ActivationPolicyProhibited},
		Assets: application.AssetOptions{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprint(w, `<!doctype html><script type="module">import "/wails/runtime.js";</script><style>
#probe{animation:spin 1s linear infinite}@keyframes spin{to{transform:rotate(1turn)}}
html.magpie-out-of-sight #probe{animation-play-state:paused!important}
</style><i id="probe"></i><script>
const probe = document.querySelector("#probe");
// A Mac window hidden by orderOut leaves its page drawing: the app says so
// and reports document.hidden itself (#4293), and so does this page, or the
// report below could never be written (the window never tells the page).
let outOfSight = false;
for (const [k, off] of [["hidden", true], ["visibilityState", "hidden"]]) {
  const own = Object.getOwnPropertyDescriptor(Document.prototype, k)?.get;
  own && Object.defineProperty(Document.prototype, k, { configurable: true, get() { return outOfSight ? off : own.call(this); } });
}
let hiddenTold = false, showing = false;
window.magpieWindowVisibility = (visible) => {
  outOfSight = !visible;
  document.documentElement.classList.toggle("magpie-out-of-sight", !visible);
  showing = visible;
  document.dispatchEvent(new Event("visibilitychange"));
  showing = false;
};
// a hiding is what follows a hide: the shown signal's own event is not one,
// though WebKit may still say hidden then
document.addEventListener("visibilitychange", () => {
  if (document.hidden && !showing) {
    localStorage.setItem("hiddenProbe", JSON.stringify({
      hidden: document.hidden,
      state: document.visibilityState,
      marked: document.documentElement.classList.contains("magpie-out-of-sight"),
      play: getComputedStyle(probe).animationPlayState,
    }));
  }
});
setInterval(() => {
  const probe = localStorage.getItem("hiddenProbe");
  if (probe) fetch("/report?probe=" + encodeURIComponent(probe));
}, 100);
fetch("/ready");
</script>`)
			case "/ready":
				select {
				case ready <- struct{}{}:
				default:
				}
				w.WriteHeader(http.StatusNoContent)
			case "/report":
				select {
				case reports <- r.URL.Query():
				default:
				}
				w.WriteHeader(http.StatusNoContent)
			default:
				http.NotFound(w, r)
			}
		})},
	})
	h := &host{app: app}
	fail := func(format string, a ...any) {
		fmt.Printf("window visibility: "+format+"\n", a...)
		os.Exit(1)
	}
	wait := func(what string, ch <-chan url.Values) url.Values {
		select {
		case v := <-ch:
			return v
		case <-time.After(8 * time.Second):
			fail("timed out waiting for %s", what)
			return nil
		}
	}
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		go func() {
			var win *application.WebviewWindow
			application.InvokeSync(func() {
				win = app.Window.NewWithOptions(application.WebviewWindowOptions{
					Name:   "visibility",
					Title:  "visibility",
					URL:    "/",
					Width:  480,
					Height: 320,
				})
				h.main = win
			})
			select {
			case <-ready:
			case <-time.After(8 * time.Second):
				fail("timed out waiting for the page")
			}

			application.InvokeSync(h.hideMain)
			time.Sleep(500 * time.Millisecond)
			application.InvokeSync(func() { showWebview(win).Focus() })
			report := wait("the page's hidden report", reports)
			probe := struct {
				Hidden bool   `json:"hidden"`
				State  string `json:"state"`
				Marked bool   `json:"marked"`
				Play   string `json:"play"`
			}{}
			if err := json.Unmarshal([]byte(report.Get("probe")), &probe); err != nil {
				fail("the page's hidden event is %q: %v", report.Get("probe"), err)
			}
			if !probe.Hidden || probe.State != "hidden" || !probe.Marked || probe.Play != "paused" {
				fail("while hidden the page saw hidden=%v state=%q marked=%v play=%q", probe.Hidden, probe.State, probe.Marked, probe.Play)
			}
			fmt.Println("window visibility: ok")
			app.Quit()
		}()
	})
	if err := app.Run(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
