function wifi_connect()
    wifi.setmode(wifi.STATION)
    wifi.sta.config(cfg.wifi.ssid,cfg.wifi.key) 
    print("MAC")
    print(wifi.sta.getmac())
    
    print("IP")
    print(wifi.sta.getip()) 
end

