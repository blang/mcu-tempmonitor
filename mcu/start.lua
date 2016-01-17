print("Process config.lua")
dofile("config.lua")

print("Process wifi.lua")
dofile("wifi.lua")

print("Process temp.lua")
dofile("temp.lua")


print("Process temp_poller.lua")
dofile("temp_poller.lua")

print("Startup")
wifi_connect()
start_temp_poller()
print("Startup finished")