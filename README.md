# IDAGIO-Downloader
IDAGIO downloader written in Go.
![](https://i.imgur.com/Lk0Ms3J.png)
[Windows binaries](https://github.com/Sorrow446/IDAGIO-Downloader/releases)

# Setup
Input credentials into config file.
Configure any other options if needed.
|Option|Info|
| --- | --- |
|email|Email address.
|password|Password.
|format|Download quality. 1 = 192 Kbps MP3, 2 = 320 Kbps MP3, 3 = 16-bit / 44.1 kHz FLAC.
|outPath|Where to download to. Path will be made if it doesn't already exist.
|trackTemplate|Track filename template. Vars: album, albumArtist, artist, copyright, genre, title, track, trackPad, trackTotal, upc, year.
|downloadBooklets|Download digital booklets when available.
|maxCoverSize|Fetch covers in their max sizes. true = max, false = 600x600.
|keepCover|Keep covers in album folders.

# Usage
Args take priority over the config file.

Download two albums:   
`idagio_dl_x64.exe https://app.idagio.com/albums/barber-adagio-for-strings-4937b163-2a55-426d-860a-ea47b418a738 https://app.idagio.com/albums/bruckner-symphony-no-3-in-d-minor-wab-103-wagner-1873-version-ed-l-nowak`

Download a single album and from two text files:   
`idagio_dl_x64.exe https://app.idagio.com/albums/barber-adagio-for-strings-4937b163-2a55-426d-860a-ea47b418a738 G:\1.txt G:\2.txt`

```
 _____ ____  _____ _____ _____ _____    ____                _           _
|     |    \|  _  |   __|     |     |  |    \ ___ _ _ _ ___| |___ ___ _| |___ ___
|-   -|  |  |     |  |  |-   -|  |  |  |  |  | . | | | |   | | . | .'| . | -_|  _|
|_____|____/|__|__|_____|_____|_____|  |____/|___|_____|_|_|_|___|__,|___|___|_|

Usage: main.exe [--format FORMAT] [--outpath OUTPATH] URLS [URLS ...]

Positional arguments:
  URLS

Options:
  --format FORMAT, -f FORMAT [default: -1]
  --outpath OUTPATH, -o OUTPATH
  --help, -h             display this help and ex
  ```
  
# Disclaimer
- I will not be responsible for how you use IDAGIO Downloader.    
- IDAGIO brand and name is the registered trademark of its respective owner.    
- IDAGIO Downloader has no partnership, sponsorship or endorsement with IDAGIO.
