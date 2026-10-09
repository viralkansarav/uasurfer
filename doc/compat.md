# Compatibility with the BrightMountainMediaInc fork

bs-bidder was built against
[BrightMountainMediaInc/uasurfer](https://github.com/BrightMountainMediaInc/uasurfer),
a fork of uasurfer as it stood in 2024. This fork tracks the current upstream
and keeps unchanged what bs-bidder reads from a parse. bs-bidder reads two
things:

- the device type, which it maps onto an OpenRTB device type when a bid request
  states none;
- `OS.Name.StringTrimPrefix()`, which becomes the request's OS when the request
  states none it recognises. Campaigns target that string exactly, and bs-bidder
  compares it against `"RokuOS"`.

## What is kept

| | upstream | this fork, as the BrightMountainMediaInc fork |
|---|---|---|
| Apple TV OS name | `TvOS` | `tvOS` |
| Roku OS name | `Roku` | `RokuOS` |
| Apple TV platform name | `AppleTV` | `tvOS` |
| Roku platform | `PlatformLinux` | `PlatformRokuTV` |
| iPad OS | `OSiPadOS` | `OSiOS`, on `PlatformiPad` |
| iPod touch device | `DevicePhone` | `DeviceTablet` |
| constants | `OSTvOS`, `OSRoku`, `PlatformAppleTV` | also `OStvOS`, `OSRokuOS`, `PlatformtvOS`, as aliases |

`OSiPadOS` is still declared, so no constant changes value, but nothing reports
it: ad requests and campaign targeting treat iPadOS as iOS.

The detection that fork added is carried over: its television markers (Comcast
Firebolt, Sky OTT, Samsung's HTTP agent, Telly, Gemini, Bell Streamer, Philips
4K, STI6110, C02_2K, Mi projectors, Stick HD, AR2101 and the Echo Show 15),
Opera Presto as a phone, FFmpeg's `Lavf` as Android, Tizen without Linux, and
the `tvOS` token of SDK agents such as `IMAtvOS`.

The go line is 1.24.0, the floor bs-bidder builds on, with nothing required.

## What changes, on purpose

Where that fork was wrong and upstream is right, this fork follows upstream.
Over 6,349 agents, from the fixture sets, both test tables, bs-bidder's tests and
a sample of CTV traffic, these are the only differences in the two fields
bs-bidder reads:

- **Bots.** About 1,500 crawlers that fork did not recognise now report `OSBot`
  and `DeviceComputer`, as the bots it did recognise always have. A bid request
  carrying a crawler's agent and no device type now gets a PC device type rather
  than none; `IsBot()` is the way to tell.
- **Android tablets.** 349 tablets that read as phones now read as tablets.
- **Televisions.** About 100 televisions, sticks and boxes that read as phones or
  as unknown now read as `DeviceTV`. Opera's HbbTV browsers among them, which the
  Presto rule used to claim.
- **OS names on televisions.** Samsung's and LG's televisions report `Tizen` and
  `WebOS` rather than `Linux`.
- **Smaller corrections.** iOS apps that read as computers, the Fire Phone that
  read as a tablet, an iOS game read as a PlayStation, an Apple Watch read as
  tvOS, and a handful of agents that fork took for bots.

Fields bs-bidder does not read also change: OS versions are now parsed for
televisions, and in-app browsers report their app (`BrowserFacebook`) rather than
the engine underneath.

## How it is held

`testdata/compat.tsv` carries every agent from that fork's test table with its
expected device type and OS name, and `compat_test.go` asserts them alongside
the names and rules above. An upstream merge that changes what bs-bidder sees
fails there first.
