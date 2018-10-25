# ![Logo](resources/site/qclauncher_logo_med.png) QCLauncher


----------

What Is QCLauncher?
-------------

QCLauncher is a small tool for [Quake Champions](https://www.quake.com). With it, you can launch Quake Champions *without running the Bethesda Launcher or having the Bethesda Launcher open*. QCLauncher has just 2 files, uses very few resources, and can be configured to immediately exit after launching the game. You can download it [here.](https://github.com/syncore/qclauncher/releases)

 :zap: **It is very important to recognize that QCLauncher does not entirely replace the Bethesda Launcher. Most importantly, if you are not using the [Early Access Steam version of QC](http://store.steampowered.com/app/611500/Quake_Champions/) you will still need the Bethesda Launcher to download any QC game updates** or to verify your game files if necessary. QCLauncher will tell you when Quake Champions game updates are available and you will be unable to play if you do not have the latest version of Quake Champions from Bethesda.

Why?
----
 - The Bethesda Launcher:
	 -  Is basically a [special version](https://bitbucket.org/chromiumembedded/cef) of the Chrome web browser that sits in the background while you play. Depending on your system, this may be heavy on resources for you.
	 - May display [messages at inconvenient times](https://www.reddit.com/r/QuakeChampions/comments/6kffch/dear_bethesda_do_not_do_this/).
 - Quickly and easily enable undocumented game options.

 - Easily add Quake Champions to Steam and use the Steam overlay if you have not purchased the Quake Champions Early Access version through Steam.

----------
![Main window](resources/site/screenshot.png)


Download
-------------

:floppy_disk: You can download both the latest version of the binary and source code [from the releases page.](https://github.com/syncore/qclauncher/releases)

Requirements
-------------

 - A free [Bethesda account](https://account.bethesda.net/en/join) with [Quake Champions](https://quake.bethesda.net/en/signup) installed and added to your account.
 - A 64-bit version of Windows (you can't play QC without this anyway)

How to Use (Setup)
-------------

 1. Download the [latest release](https://github.com/syncore/qclauncher/releases) and extract the `qclauncher.exe` file from the zip file.
 2. Double click `qclauncher.exe` to run QCLauncher.
 3. Click the 'Configure' button and enter the requested information to configure your settings. For the QC user name and password, this will be the same info used for the Bethesda launcher (or the Bethesda forums).
 4. When selecting the QC exe, the default location is: `C:\Program Files (x86)\bethesda.net Launcher\games\quakechampions\client\bin\pc`
 5. *Steam (Optional)*: If you want to add Quake Champions as a non-Steam game, this can be done under the 'Launcher Settings' tab. Click the check box labeled 'Add as a non-Steam Game (for Steam overlay)'. After you save your settings, Steam will open. Find and select `qclauncher.exe` in Steam to add it as a non-Steam game. You can rename it to Quake Champions if you want, so that it will be displayed that way in your friends list.
 6. Click the 'Save All' button. If successful, you should be able to play by clicking the 'Play' button.

New game options have been found since the last QCLauncher release, how can I try these new options?
-------------
Since version 1.01, it has been possible to pass custom Quake Champions start-up options to QCLauncher with the `--customargs` flag. For example, create a shortcut to  QCLauncher or start QCLauncher in this manner:

`qclauncher.exe --customargs="--set /Config/CONFIG/WeaponZScale -10 --set /Config/CONFIG/isLowResParticles 1"`

QCLauncher will then pass these options to Quake Champions on launch.

Developers: Build from Source Code (you can skip this if you don't plan on working on the code)
-------------

 1. Download and install the latest stable release of the Go Programming Language, which is [available here.](https://golang.org/dl/)
 2. *Windows* - To get the QCLauncher source: `go get -ldflags="-H windowsgui -s -w" github.com/syncore/qclauncher`
 3. *Windows* - If you do *not* have Microsoft Visual Studio 2017 *Enterprise Edition*, download the [Build Tools For Visual Studio 2017](https://visualstudio.microsoft.com/downloads/#build-tools-for-visual-studio-2017) from Microsoft. Afterwards, adjust the `msBuildDir` variable in [this file](resources/bin_src/build_blff_src.bat). If you have VS2017, but not the Enterprise edition (i.e. Professional), adjust the `msBuildDir` variable to point to its `MSBuild.exe` executable location.
 4. *Linux* - To get the QCLauncher source: `GOOS=windows GOARCH=amd64 go get -d -ldflags="-H windowsgui -s -w" github.com/syncore/qclauncher`
 5. Find your GOPATH. This can be found by entering:  `go env GOPATH` on the command line.
 6. Change directory to `GOPATH\src\github.com\syncore\qclauncher`
 7. *Windows* - To build, run `build.bat`
 8. *Linux* - To build, run `build.sh` (the application only runs on Windows, but can be built on Linux/OSX).
 9. If everything went well, you should have the `qclauncher.exe` file in the `bin` directory.

Is QCLauncher Considered a Cheat?
-------------
No. QCLauncher **does *NOT* touch or modify any game files or game code at all**. Any additional functionality that QCLauncher provides is derived from the game itself and the game's built-in commands. The tool is simply a very lightweight utility that launches the game. Use it if you'd like to, or not. I wrote it as a learning exercise in the [tradition](https://qlprism.syncore.org/) of [contributing](https://ql.syncore.org) to the Quake [community](https://qlprism.syncore.org/qlm/). It's open-source. Inspect the code and you will see that there is no funny business going on.

Issues, Contact Me, etc.
-------------

I can be contacted under the name **syncore** on [Discord](https://discordapp.com/). Any other issues can be opened on the [issue tracker here on Github,](https://github.com/syncore/qclauncher/issues) and I will try to address them, time permitting. Additionally, there is an ESR thread available [here](http://www.esreality.com/post/2877585/quake-champions-quick-launcher/) and a PlusForward.net thread available [here.](https://www.plusforward.net/quake/post/28904/QCLauncher-Run-Quake-Champions-without-Bethesda-Launcher/)


:thumbsup: :video_game: *Happy fragging, and please support this latest Quake title! Download it for free at https://www.quake.com*
