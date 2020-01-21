# killall -9 persona
# nohup ./persona -env debug -ports 9008 &
# # debug alpha release

# # pm2 test
pm2 delete persona_debug
pm2 start persona_debug.sh