# Connected TV

Televisions, sticks, set top boxes and Android projectors report `DeviceTV`.
Where the platform states an OS of its own, that is reported too.

| platform | OS name | platform value | OS version |
|---|---|---|---|
| Samsung smart TVs | `OSTizen` | `PlatformLinux` | yes |
| LG smart TVs | `OSWebOS` | `PlatformLinux` | not stated by the device |
| Roku players | `OSRoku`, named `RokuOS` | `PlatformRokuTV` | yes |
| Apple TV | `OSTvOS`, named `tvOS` | `PlatformAppleTV`, named `tvOS` | yes |
| Amazon Fire TV | `OSAndroid` | `PlatformLinux` | the Android version |
| Android TV, Google TV, Chromecast | `OSAndroid` | `PlatformLinux` | the Android version |
| Vizio, Hisense VIDAA, HbbTV boxes, RDK, NetCast, Viera, Bravia | `OSLinux` or `OSUnknown` | `PlatformLinux` | not stated by the device |
| GStreamer, the media framework TV players fetch video with | `OSUnknown` | `PlatformUnknown` | not stated by the device |

The Roku and Apple TV names are the ones the BrightMountainMediaInc fork used,
kept for the callers built against it; see [compat.md](compat.md).

The last row is not an omission: those platforms are Linux and state no version
of their own, so a constant per brand would buy a name and nothing else. Use
`DeviceTV` to know it is a television, and the browser and its version for what
the device can actually render.

Android TV images are recognised by their build ID as well as their model.
Their builds carry "TT" after the release letter (`Build/PTT1.190222.001`,
`Build/UTT2.250416.001`), which no phone or tablet build in the fixtures does,
and that is what identifies a television whose model is an opaque code, such as
TD Systems' `K55DLG12US`.

Amazon's Fire TV models are recognised as a family, so a stick released after
this version of the library still reports `DeviceTV`. Fire tablets remain
`DeviceTablet`.

Two notes that look like bugs and are not:

* LG spells its own OS **`Web0S`**, with a zero, on everything since webOS 3.
  All three spellings map to `OSWebOS`.
* The Roku and Apple TV **remote control apps** run on a phone, and report as
  one. Only the players and the boxes report `DeviceTV`.
