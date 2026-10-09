package uasurfer

import (
	"strings"
)

// Markers that are a substring of another marker are omitted: "tv" already
// subsumes googletv/appletv/smarttv/smart-tv/hbbtv/dtv/"tv box", "telly"
// subsumes "stableavb_telly", "bravia" "bravia 4k", and " x96" the
// whole X96 box family - with its leading space, because "540x960" is a screen
// resolution and X96 boxes state their model as a field of its own. Amazon's Fire TV models are not listed at all; isFireTV
// matches the whole family.
var tvMarkers = []string{
	"tv", "roku", "crkey", "chromecast", "stb", "tuner", "vizio", "viera", "aquos", "bravia",
	"netcast", "youview", "adt-", "swisscom-ip", "mibox", "ott-g1", "ottera",
	// operator set-top boxes. The French ISP boxes are worth naming: every
	// reference parser reads them as a desktop or a phone.
	"freebox", "sfrwpebrowser", "mfi_airplay",
	"tpm191e", "tpm171e", "nokia streaming box", "lxbox51",
	" x96", "canal plus box", "vectra 4k box",
	"diw377", "diw380", "dv8555", "dctiw362", "gd1 4k", "ai pont", "b-stream",
	// Android boxes, sticks and projectors whose model names no television: MEO's
	// operator box, the boxes built on H313, Allwinner's TV box chip, Magcubic's
	// HY300 projectors and the M98-Y10 stick. Each starts with the "; " Android
	// puts before the model, which anchors it to that field as " x96" is anchored,
	// and costs the scan far less: ";" turns up a few times in an agent, a space
	// at every word.
	"; meobox", "; h313", "; hy300", "; m98-y10",
	// Carried over from the BrightMountainMediaInc fork, whose televisions they
	// are: Comcast's Firebolt and Sky's WPE boxes, Samsung's own HTTP agent,
	// Telly's free TV, the Gemini and Bell streamers, a Philips TV, the
	// STI6110, C02_2K and Mi projector boxes, the AR2101 build some Android
	// TVs ship, and "aeoh", the Echo Show 15 with Fire TV built in.
	"firebolt", "sky_ott", "samsung-agent", "telly", "gemini", "bell streamer",
	"philips 4k", "sti6110", "c02_2k", "miproj", "stick hd", "ar2101", "aeoh",
	// The media framework of Hisense VIDAA, LG and Samsung players, fetching on
	// their behalf. Desktop Linux ships it too, but a desktop's ads fire from
	// a browser, not from its media framework.
	"gstreamer",
}

// tvMarkersByFirstByte buckets tvMarkers so isTV can scan the agent once
// instead of running a full strings.Contains pass per marker. This is the
// hottest check in the parser: every non-desktop, non-iOS agent walks the
// whole set before being classified.
var tvMarkersByFirstByte = func() (buckets [256][]string) {
	for _, m := range tvMarkers {
		buckets[m[0]] = append(buckets[m[0]], m)
	}
	return
}()

func isTV(ua string) bool {
	for i := 0; i < len(ua); i++ {
		for _, m := range tvMarkersByFirstByte[ua[i]] {
			if strings.HasPrefix(ua[i:], m) {
				return true
			}
		}
	}

	if strings.Contains(ua, "mbox") && !strings.Contains(ua, "xbox") {
		return true
	}

	// Last, because they are the only checks that have to look at where in the
	// agent the match sits: every other marker here is distinctive enough to
	// match anywhere.
	agentPlatform := platformGroup(ua)
	return isAndroidTVBuild(agentPlatform) || isFireTV(agentPlatform)
}

// isAndroidTVBuild reports whether agentPlatform carries the build ID of an
// Android TV image: "Build/PTT1.190222.001", "Build/STTL.240206.002". That is
// the release letter (P is Android 9, U is 14), "TT", one more letter or digit,
// a dot and the six digit date the build was cut. Phone and tablet builds put
// something else after the release letter - "TP1A", "UQ1A", "SKQ1" - and none
// of those in testdata/devices.tsv carries "TT".
//
// It is what recognises a television whose model is an opaque code no marker
// could name, such as TD Systems' "K55DLG12US". The whole shape is required, so
// a custom ROM's "Build/Attila" is not a TV. Only the platform group is
// searched, because that is where Android states its build, once: scanning the
// whole agent cost every Chrome on Android a pass that never finds one.
func isAndroidTVBuild(agentPlatform string) bool {
	_, id, _ := strings.Cut(agentPlatform, "build/")
	if len(id) < 11 || !isLower(id[0]) || id[1] != 't' || id[2] != 't' ||
		!isLower(id[3]) && !isDigit(id[3]) || id[4] != '.' {
		return false
	}
	for i := 5; i < 11; i++ {
		if !isDigit(id[i]) {
			return false
		}
	}
	return true
}

// isFireTV reports whether specs carries an Amazon Fire TV model token: "aft"
// and up to ten more letters or digits, as its own field.
//
// Amazon has shipped dozens of these - aftb, aftmm, aftka, aftkm, aftdck,
// aftkmst12, aftbamr311 - and enumerating them meant every new stick read as a
// phone until someone noticed. The family prefix is stable, so match on that.
//
// Only the platform group is searched, never the whole agent: unanchored, "aft"
// also begins "after", which turns any agent quoting a URL like
// "afterice.se" into a television.
func isFireTV(specs string) bool {
	for i := 0; i+3 <= len(specs); i++ {
		if i > 0 && !isFieldEdge(specs[i-1]) {
			continue
		}
		if specs[i] != 'a' || specs[i+1] != 'f' || specs[i+2] != 't' {
			continue
		}
		n := i + 3
		for n < len(specs) && n-i <= 12 && (isLower(specs[n]) || isDigit(specs[n])) {
			n++
		}
		if n > i+3 && (n == len(specs) || isFieldEdge(specs[n])) {
			return true
		}
	}
	return false
}

// isFieldEdge reports whether c separates two fields of a platform group.
func isFieldEdge(c byte) bool { return isSpace(c) || c == ';' || c == ',' }

// isBrowserAgent reports whether ua has the shape a browser announces itself
// with: the Mozilla preamble every one of them still carries, and a rendering
// engine.
//
// Naming an engine is not enough on its own. A decade of vendor software put
// "AppleWebKit" in agents of its own invention -
// "Android 4.0.4;AppleWebKit/534.30;Build/HuaweiU8836D;U8836D" - and those are
// phones, and app agents besides. Requiring the preamble is what separates a
// browser from software that merely mentions one.
func isBrowserAgent(ua string) bool {
	return strings.HasPrefix(ua, "mozilla/") &&
		(strings.Contains(ua, "applewebkit") || strings.Contains(ua, "gecko"))
}

// Android tablet models, for the tablets that state "mobile" like a phone -
// which they should not, and often do:
// http://android-developers.blogspot.com/2010/12/android-browser-user-agent-issues.html
//
// Samsung has three tablet lines and no phone shares their prefixes: "SM-X" is
// the current Tab A9 and S10, "SM-T" the older Tabs, "SM-P" the ones with a pen
// (Tab S6 Lite, S7). Its phones are SM-G, SM-A, SM-S, SM-N and SM-F.
//
// "; t1-" is Huawei's MediaPad T1 and needs its dash: "; t1" alone also matches
// a "T100" handset. Lenovo's bare Tab codes need their digit for the same reason:
// "; tb" alone also matches "TB-7000", a handset. The dashed Lenovo models -
// "Lenovo TB-X606F" - all carry the brand, so "lenovo tb" reaches those.
//
// "-w09" is the one marker here that is not exact. It is Huawei and Honor's
// WiFi-only variant suffix, which is a tablet 300 times in a 200k sample of real
// traffic and a phone 21 times - "JMS-W09" is a handset. It is in because 300
// tablets read as phones is the larger error, but it is the first thing to remove
// if those phones matter more; nothing else in this list has a known exception.
var tabletMarkers = []string{
	"tablet", "nexus 7", "nexus 9", "nexus 10", "xoom",
	"sm-t", "sm-x", "sm-p", "; kf", "; t1-",
	"lenovo tab", "lenovo tb", "; tb1", "; tb3", // Lenovo writes its Tabs both ways
	"rmp2", // Xiaomi's Redmi Pad
	"-w09", // Huawei and Honor's WiFi-only variants
}

// isAndroidTablet reports whether ua names a model known to be a tablet.
//
// A plain loop rather than the bucketed scan isTV uses: at this length the
// Contains calls are cheaper than a pass over the agent, and the list is a table
// so that a new model line is one line and a test can enumerate it.
func isAndroidTablet(ua string) bool {
	for _, m := range tabletMarkers {
		if strings.Contains(ua, m) {
			return true
		}
	}
	return false
}

// isAndroidWearable reports whether an Android agent belongs on a wrist.
//
// There is no form factor field to read: Wear OS states "Mobile" like a phone
// does, so the model name is all there is. "sm-r" is Samsung's Galaxy Watch
// line, whose phones and tablets are "sm-g" and "sm-t".
func isAndroidWearable(ua string) bool {
	return strings.Contains(ua, "watch") ||
		strings.Contains(ua, "wear os") ||
		strings.Contains(ua, "sm-r") ||
		strings.Contains(ua, "glass")
}

func (u *UserAgent) parseDevice(ua string) {
	switch {

	case u.OS.Platform == PlatformWindows || u.OS.Platform == PlatformMac || u.OS.Name == OSChromeOS:
		if strings.Contains(ua, "mobile") || strings.Contains(ua, "touch") {
			u.DeviceType = DeviceTablet // windows rt, linux haxor tablets
			return
		}
		u.DeviceType = DeviceComputer

	// long list of smarttv and tv dongle identifiers - above "phone" and "tablet" check to prevent TVs from being detected as phones/tablets
	case u.OS.Platform != PlatformiPhone && u.OS.Platform != PlatformiPad && isTV(ua):
		u.DeviceType = DeviceTV

	// The iPod touch reports as a tablet, as it did in the BrightMountainMediaInc
	// fork and in uasurfer before it. Upstream now files it with the phones;
	// this fork keeps the old answer so callers see no change. See
	// doc/compat.md.
	case u.OS.Platform == PlatformiPad || u.OS.Platform == PlatformiPod || strings.Contains(ua, "tablet") || strings.Contains(ua, "kindle/") || strings.Contains(ua, "playbook"):
		u.DeviceType = DeviceTablet

	case u.OS.Platform == PlatformiPhone || u.OS.Platform == PlatformBlackberry || strings.Contains(ua, "phone"):
		u.DeviceType = DevicePhone

	case u.OS.Name == OSAndroid:
		// Wear OS names the watch in the model field and says nothing else about
		// itself, so it has to be read before the phone default below claims it.
		if isAndroidWearable(ua) {
			u.DeviceType = DeviceWearable
			return
		}

		// A model known to be a tablet settles it, and is read before the token
		// below rather than after: the whole point of the list is the tablets
		// that state "mobile" anyway.
		if isAndroidTablet(ua) {
			u.DeviceType = DeviceTablet
			return
		}

		if strings.Contains(ua, "mobile") {
			u.DeviceType = DevicePhone
			return
		}

		// No "mobile" token, so the version decides at both ends of the history.
		// Honeycomb shipped on tablets and nothing else; before it, the token
		// was not yet a convention anyone kept, and the devices of that era are
		// overwhelmingly phones.
		switch {
		case u.OS.Version.Major == 3:
			u.DeviceType = DeviceTablet
			return
		case u.OS.Version.Major > 0 && u.OS.Version.Major < 3:
			u.DeviceType = DevicePhone
			return
		}

		// From Android 4 on, every browser omits "mobile" on a tablet and states
		// it on a phone. Since Chrome's user agent reduction - "(Linux; Android
		// 10; K)", no model, no real version - that token is the only signal
		// left.
		//
		// Browsers only: an app HTTP stack states no form factor either way, and
		// a Dalvik or a download manager on a phone would otherwise read as a
		// tablet.
		if isBrowserAgent(ua) {
			u.DeviceType = DeviceTablet
			return
		}

		u.DeviceType = DevicePhone // default to phone

	case u.OS.Platform == PlatformPlaystation || u.OS.Platform == PlatformXbox || u.OS.Platform == PlatformNintendo:
		u.DeviceType = DeviceConsole

	case strings.Contains(ua, "glass") || strings.Contains(ua, "watch") || strings.Contains(ua, "sm-v"):
		u.DeviceType = DeviceWearable

	// Above the "mobile" check, because Kindle Fire tablets report as "mobile".
	// The exclusion has to cover the Silk arm as well: the Fire Phone runs Silk
	// too, and reading "|| OSKindle && !sd4930ur" the way Go groups it left the
	// phone matching on Silk alone.
	case (u.Browser.Name == BrowserSilk || u.OS.Name == OSKindle) && !strings.Contains(ua, "sd4930ur"):
		u.DeviceType = DeviceTablet

	// Opera's Presto engine outside a desktop is Opera Mini and Opera Mobile on
	// feature phones, which state nothing else about the device. Carried over
	// from the BrightMountainMediaInc fork, and placed after the Android case,
	// which settles Opera Mobile on Android itself, so Android agents are not
	// scanned for it. Desktop Linux is excluded: Opera 12 ran Presto there too,
	// and that falls through to the computer case below.
	case strings.Contains(ua, "presto") && u.OS.Name != OSLinux:
		u.DeviceType = DevicePhone

	case strings.Contains(ua, "mobile") || strings.Contains(ua, "touch") || strings.Contains(ua, " mobi") || strings.Contains(ua, "webos"): //anything "mobile"/"touch" that didn't get captured as tablet, console or wearable is presumed a phone
		u.DeviceType = DevicePhone

	case u.OS.Name == OSLinux: // linux goes last since it's in so many other device types (tvs, wearables, android-based stuff)
		u.DeviceType = DeviceComputer

	default:
		u.DeviceType = DeviceUnknown
	}
}
