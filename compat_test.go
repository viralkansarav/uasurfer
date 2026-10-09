package uasurfer

import "testing"

// TestCompatNames pins the names the BrightMountainMediaInc fork exported, and
// the strings they print: callers compare against the strings and target
// campaigns by them.
func TestCompatNames(t *testing.T) {
	if OStvOS != OSTvOS || OSRokuOS != OSRoku || PlatformtvOS != PlatformAppleTV {
		t.Error("compat aliases no longer alias upstream's constants")
	}
	for _, c := range []struct{ got, want string }{
		{OStvOS.StringTrimPrefix(), "tvOS"},
		{OSRokuOS.StringTrimPrefix(), "RokuOS"},
		{OSTizen.StringTrimPrefix(), "Tizen"},
		{PlatformtvOS.StringTrimPrefix(), "tvOS"},
		{PlatformRokuTV.StringTrimPrefix(), "RokuTV"},
	} {
		if c.got != c.want {
			t.Errorf("name = %q, want %q", c.got, c.want)
		}
	}
}

// TestCompatFixtures runs the BrightMountainMediaInc fork's agents and asserts
// the two fields bs-bidder reads.
func TestCompatFixtures(t *testing.T) {
	for _, row := range readFixtures(t, "compat.tsv", 3) {
		wantDevice, wantOS, agent := row[0], row[1], row[2]
		ua := Parse(agent)
		if got := ua.DeviceType.StringTrimPrefix(); got != wantDevice {
			t.Errorf("device = %s, want %s: %s", got, wantDevice, agent)
		}
		if got := ua.OS.Name.StringTrimPrefix(); got != wantOS {
			t.Errorf("os = %s, want %s: %s", got, wantOS, agent)
		}
	}
}

// TestCompatRules covers the rules carried over from the BrightMountainMediaInc
// fork that its test table does not: an iPad reports iOS on every path, Roku
// keeps its own platform, and Lavf, Presto and the tvOS SDK token each settle
// an agent that states nothing else.
func TestCompatRules(t *testing.T) {
	for _, c := range []struct {
		agent    string
		platform Platform
		os       OSName
		device   DeviceType
	}{
		{"Mozilla/5.0 (iPad; CPU OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1", PlatformiPad, OSiOS, DeviceTablet},
		{"TuneIn Radio/27.1.0 (iPadOS/16.6)", PlatformiPad, OSiOS, DeviceTablet},
		{"Mozilla/5.0 (iPod touch; CPU iPhone OS 15_8_8 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/19H384", PlatformiPod, OSiOS, DeviceTablet},
		{"Roku/DVP-12.5 (12.5.0.4174-46)", PlatformRokuTV, OSRoku, DeviceTV},
		{"IMAtvOS/t.3.8.0 tvOS/13.4", PlatformAppleTV, OSTvOS, DeviceTV},
		{"Lavf/58.22.100", PlatformLinux, OSAndroid, DevicePhone},
		{"Opera/9.80 (Tizen; Opera Mini/7.7.18/174.55; U; en) Presto/2.12.423 Version/12.16", PlatformLinux, OSTizen, DevicePhone},
		{"Opera/9.80 (BREW; Opera Mini/5.1/27.2213; U; en) Presto/2.8.119 240X400 LG VN27181", PlatformUnknown, OSUnknown, DevicePhone},
		// Desktop Opera ran Presto too, and stays a computer.
		{"Opera/9.80 (X11; Linux x86_64) Presto/2.12.388 Version/12.16", PlatformLinux, OSLinux, DeviceComputer},
		{"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.16", PlatformWindows, OSWindows, DeviceComputer},
	} {
		ua := Parse(c.agent)
		if ua.OS.Platform != c.platform || ua.OS.Name != c.os || ua.DeviceType != c.device {
			t.Errorf("Parse(%q) = %v %v %v, want %v %v %v", c.agent,
				ua.OS.Platform, ua.OS.Name, ua.DeviceType, c.platform, c.os, c.device)
		}
	}
}
