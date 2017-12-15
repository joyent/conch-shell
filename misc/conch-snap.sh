#!/bin/sh

TMPCONCH=`mktemp /tmp/conch.XXXXXX` || exit 1
curl http://us-east.manta.joyent.com/sungo/public/conch-shell/conch-snap 1> $TMPCONCH 2>/dev/null
chmod +x $TMPCONCH
$TMPCONCH $@
rm $TMPCONCH
