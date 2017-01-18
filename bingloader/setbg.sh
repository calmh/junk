#!/bin/bash
set -euo pipefail

image=$(bingloader -dir ~/Pictures/Bing)

/usr/bin/osascript<<END
tell application "Finder"
set desktop picture to POSIX file "$image"
end tell
END
