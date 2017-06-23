# ![Logo](resources/site/qclauncher_logo_med.png) QCLauncher


----------


*QCLauncher* is a quick launcher utility for [Quake Champions](https://www.quake.com). You enter your account information one time and automatically log into the game on future launches by running QCLauncher, which will exit after launch and will not sit in the background using up system resources. With this utility, you do not need to launch the Bethesda Launcher or have it open to play Quake Champions.

 **However, it is incredibly important to note that QCLauncher is not intended to entirely replace the Bethesda Launcher. Most importantly, you will still need the Bethesda Launcher to download any game updates** or to verify your game files if you need to do that. By default, QCLauncher will not allow you to play if you do not have the latest version of Quake Champions from Bethesda.

----------
![Main window](resources/site/screenshot_1.png)

![Settings](resources/site/screenshot_2.png)

Download
-------------

Both the latest version of the binary and source code are available on the [releases page here on Github.](https://github.com/syncore/qclauncher/releases)

Requirements
-------------

 - A free [Bethesda account](https://account.bethesda.net/en/join) with [Quake Champions](https://quake.bethesda.net/en/signup) installed and added to your account.
 - A 64-bit version of Windows (you can't play QC without this anyway)

How to Use
-------------

:video_camera: If you'd like, [here you can view a a short demo video](https://www.youtube.com/watch?v=z1hN6UCA_zo) on YouTube.

 1. Download the [latest release](https://github.com/syncore/qclauncher/releases) and extract the `qclauncher.exe` file from the zip file.
 2. Double click qclauncher.exe
 3. Enter the user name and password that you use for the Bethesda launcher (or the Bethesda forums, they are the same)
 4. Click the "Select QC EXE" button to locate your Quake Champions exe file; by default, this will be located at: `C:\Program Files (x86)\Bethesda.net Launcher\games\client\bin\pc\QuakeChampions.exe`
 5. If you want to add Quake Champions as a non-Steam game (via QCLauncher), click the 'Add as non-Steam Game' check box.
 5. Click the 'Save' button. Your login information will be verified.
 6. If successful, the launcher will exit. In the future, simply click the `qclauncher.exe` file to launch Quake Champions and play. If you need to re-enter your login info, delete the `data.qcl` file and re-run qclauncher.exe.


Build from Source Code
-------------

 1. Download and install the latest stable release of the Go Programming Language, which is [available here.](https://golang.org/dl/)
 2. Run the included `build.bat` file if you're on Windows or `build.sh` if you're on Linux (though this application only runs on Windows), if everything went well, you should have the `qclauncher.exe` file in the `bin` directory.

Is This Considered a Cheat?
-------------
Not to my knowledge. QCLauncher **does *NOT* touch or modify any game files or game code at all**. It is simply a very lightweight utility that launches the game. Use it if you'd like to, or not. I wrote it as a learning exercise in the [tradition](https://qlprism.syncore.org/) of [contributing](https://ql.syncore.org) to the Quake [community](https://qlprism.syncore.org/qlm/). It's open-source; inspect the code and you will see that there is no funny business going on.

Issues, Contact Me, etc.
-------------

I can be contacted under the name **syncore** on [Discord](https://discordapp.com/). Any other issues can be opened on the [issue tracker here on Github,](https://github.com/syncore/qclauncher/issues) and I will try to address them if I have time.


:thumbsup: :video_game: *Happy fragging, and please support this latest Quake title! Download it for free at https://www.quake.com*