# killall -9 persona
# nohup ./persona -env release -ports 9008,9009,9010,9011 &
# # debug alpha release
pm2 delete persona_release
pm2 start persona_release.sh

