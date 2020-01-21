# killall -9 persona
# nohup ./persona -env alpha -ports 9008,9009,9010,9011 &
# # debug alpha release

# # pm2 test
pm2 stop persona_alpha
pm2 start persona_alpha.sh
